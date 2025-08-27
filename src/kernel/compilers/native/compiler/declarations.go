package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Declarations is a structure responsible for managing compiler declarations and scope-related components.
// It holds references to constants, scopes, structs, and a gatekeeper for managing object lifecycle and interactions.
// The fileSet tracks source file information, and the compile function is used for compiling AST nodes.
type Declarations struct {
	gk           objects.IGateKeeper
	references   *Constants
	constants    *Constants
	scopes       *Scopes
	fileSet      *token.FileSet
	structsTable *StructTable
	compile      func(node ast.Node) error
}

// NewDeclarations creates and initializes a new Declarations instance with gatekeeper, constants, scopes, and structs table.
func NewDeclarations(gk objects.IGateKeeper, references *Constants, constants *Constants, scopes *Scopes, structsTable *StructTable) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile:      nil,
		structsTable: structsTable,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *Declarations) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// DeclStmt processes an AST declaration statement node and compiles its declaration, returning any encountered error.
func (c *Declarations) DeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// GenDecl processes a general declaration node, compiling each specification within the declaration. Returns an error on failure.
func (c *Declarations) GenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// TypeSpec processes an AST TypeSpec node, registering structs and their fields in the scope and structs table.
// It validates type uniqueness and collects field details like names, base types, and full types.
func (c *Declarations) TypeSpec(node *ast.TypeSpec) error {
	structType, isStruct := node.Type.(*ast.StructType)
	if !isStruct {
		return nil
	}
	structName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(structName); ok {
		return fmt.Errorf("type '%s' already defined", structName)
	}
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var typeNameBuf bytes.Buffer
			var base = ExtractBaseName(field.Type)
			if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
				return fmt.Errorf("failed to resolve type for field in struct '%s'", structName)
			}
			fieldType := typeNameBuf.String()
			for _, name := range field.Names {
				c.structsTable.Add(structName, name.Name, base, fieldType, field)
			}
		}
	}
	return nil
}

// ValueSpec processes variable declarations and ensures type inference, symbol definition, and bytecode emission.
func (c *Declarations) ValueSpec(node *ast.ValueSpec) error {
	// handles 'var x = 10'
	for i, name := range node.Names {
		if i > len(node.Values)-1 {
			return fmt.Errorf("too few values for %s", name.Name)
		}
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}

		isStruct2 := false
		structName := ""

		// 3. Inferenza del tipo, ora coerente con la nuova logica
		var assignedTypeNames []string
		if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
			if ident, ok := compLit.Type.(*ast.Ident); ok {
				if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
					structName = ident.Name
					assignedTypeNames = []string{typeSymbol.Name()}
					isStruct2 = true
				}
			}
		} else if callExpr, ok := node.Values[i].(*ast.CallExpr); ok {
			var funcName string
			if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
				funcName = ident.Name
			}
			if funcName != "" {
				if funcSymbol, ok := c.scopes.SymbolResolve(funcName); ok {
					returnTypes := funcSymbol.Types()
					if len(returnTypes) != 1 {
						return fmt.Errorf("assignment mismatch: 'var' declaration expects 1 value, but function %s returns %d", funcName, len(returnTypes))
					}
					assignedTypeNames = []string{returnTypes[0]}
				}
			}
		}

		symbol, err := c.scopes.SymbolDefine(name.Name)
		if err != nil {
			return err
		}
		symbol.SetIsStruct(structName, isStruct2)
		if len(assignedTypeNames) > 0 {
			symbol.SetTypes(assignedTypeNames)
			symbol.SetObject(c.gk.NewString(objects.FrameStatic, symbol.Name()+":"+strings.Join(assignedTypeNames, " ")))
		} else {
			symbol.SetObject(c.gk.NewString(objects.FrameStatic, symbol.Name()))
		}

		// 4. Emette bytecode per assegnare il valore dalla cima dello stack alla variabile.
		if err = c.scopes.EmitSymbolDefine(symbol); err != nil {
			return err
		}
		// 5. Pulisce lo stack dal valore ora che è stato assegnato.
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	return nil
}

// BasicLit compiles a basic literal node (e.g., int, float, string) into a bytecode representation or returns an error.
func (c *Declarations) BasicLit(node *ast.BasicLit) error {
	var obj objects.IObject
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		obj = c.gk.NewInt(objects.FrameStatic, val)
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		obj = c.gk.NewFloat(objects.FrameStatic, val)
	case token.STRING:
		val, _ := strconv.Unquote(node.Value)
		obj = c.gk.NewString(objects.FrameStatic, val)
	default:
		return fmt.Errorf("unhandled literal: %s", node.Kind)
	}
	id := c.constants.Add("", obj)
	if _, err := c.scopes.Emit(bytecode.OpConstant, id); err != nil {
		return err
	}
	return nil
}

// Ident processes an identifier node and emits bytecode if the identifier corresponds to a symbol or keyword in the scope.
func (c *Declarations) Ident(node *ast.Ident) error {
	switch node.Name {
	case "true":
		if _, err := c.scopes.Emit(bytecode.OpTrue); err != nil {
			return err
		}
		return nil
	case "false":
		if _, err := c.scopes.Emit(bytecode.OpFalse); err != nil {
			return err
		}
		return nil
	}
	symbol, ok := c.scopes.SymbolResolve(node.Name)
	if !ok {
		// Se non riusciamo a risolvere il simbolo, NON generiamo un errore.
		// Assumiamo che sia un nome di campo (es. 'Title' in 's.Title')
		// e che il suo nodo genitore (*ast.SelectorExpr) se ne sia già occupato.
		// Semplicemente, non emettiamo alcun bytecode per questo nodo.
		return nil
		//return fmt.Errorf("[Ident] undefined variable: %s", node.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	return nil
}

// AssignStmt processes assignment statements, including multiple and single assignments, and emits corresponding opcodes.
// It supports variable declarations, type inference, field assignments, and checks for type or scope mismatches.
func (c *Declarations) AssignStmt(node *ast.AssignStmt) error {
	// Gestisce l'assegnazione multipla da una chiamata di funzione (es. x, y := f())
	if callExpr, ok := node.Rhs[0].(*ast.CallExpr); ok && len(node.Lhs) > 1 {
		if err := c.compile(callExpr); err != nil {
			return err
		}
		var funcReturnTypes []string
		if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
			if ident.Name != "" {
				if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok {
					funcReturnTypes = funcSymbol.Types()
				}
			}
		}
		if len(node.Lhs) != len(funcReturnTypes) {
			return fmt.Errorf("assignment mismatch: %d variables but %d return values", len(node.Lhs), len(funcReturnTypes))
		}
		for i := len(node.Lhs) - 1; i >= 0; i-- {
			lhs := node.Lhs[i]
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				return fmt.Errorf("unsupported multiple assignment to type %T", lhs)
			}
			var symbol *Symbol
			if node.Tok == token.DEFINE {
				var err error
				symbol, err = c.scopes.SymbolDefine(ident.Name) // Anche qui, l'inferenza del tipo potrebbe essere migliorata
				if err != nil {
					return err
				}
			} else {
				var found bool
				symbol, found = c.scopes.SymbolResolve(ident.Name)
				if !found {
					return fmt.Errorf("undefined variable: %s", ident.Name)
				}
			}
			symbol.SetTypes([]string{funcReturnTypes[i]})
			if err := c.scopes.EmitSymbolSet(symbol); err != nil {
				return err
			}
			if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
				return err
			}
		}
		return nil
	}

	// Gestisce l'assegnazione singola (es. x = 1 o x := 1)
	if err := c.compile(node.Rhs[0]); err != nil {
		return err
	}

	switch lhs := node.Lhs[0].(type) {
	case *ast.Ident:
		name := lhs.Name
		var symbol *Symbol
		if node.Tok == token.DEFINE { // Caso specifico per ':='
			var err error
			// Ispezioniamo il lato destro (RHS) per inferire il tipo
			isStruct2 := false
			structName := ""
			var assignedTypeName []string

			switch rhs := node.Rhs[0].(type) {
			case *ast.BasicLit:
				//nothing to do
			case *ast.CompositeLit: // es. MyStruct{...}
				baseName := ExtractBaseName(rhs.Type)
				if len(baseName) > 0 {
					assignedTypeName = []string{baseName}
					structName = baseName
					isStruct2 = true
				}
			case *ast.CallExpr: // es. NewStruct()
				if ident, ok := rhs.Fun.(*ast.Ident); ok {
					if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
						// Assumiamo il primo tipo di ritorno
						typeName := funcSymbol.Types()[0]
						assignedTypeName = []string{typeName}
						// Verifichiamo se il tipo restituito è uno struct
						if typeSymbol, ok := c.scopes.SymbolResolve(typeName); ok && typeSymbol.IsStruct() {
							isStruct2 = true
							structName = typeName
						}
					}
				}
			case *ast.UnaryExpr: // es. &MyStruct{}
				if rhs.Op == token.AND {
					if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
						if ident, ok := compLit.Type.(*ast.Ident); ok {
							if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
								isStruct2 = true
								structName = typeSymbol.Name()
								assignedTypeName = []string{typeSymbol.Name()}
							}
						}
					}
				}
			default:
				return fmt.Errorf("unsupported right-hand side for assignment: %T", rhs)
			}
			// Definiamo il simbolo usando il valore 'isStruct' appena calcolato
			symbol, err = c.scopes.SymbolDefine(name)
			if err != nil {
				return err
			}
			symbol.SetIsStruct(structName, isStruct2)
			// Se abbiamo un tipo, lo associamo al simbolo
			if len(assignedTypeName) > 0 {
				symbol.SetTypes(assignedTypeName)
			}
		} else { // Caso per l'assegnazione normale '='
			var ok bool
			symbol, ok = c.scopes.SymbolResolve(name)
			if !ok {
				return fmt.Errorf("[AssignStmt] undefined variable: %s", name)
			}
		}
		if err := c.scopes.EmitSymbolSet(symbol); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		return nil

	case *ast.SelectorExpr: // es. myStruct.Field = ...
		if node.Tok == token.DEFINE {
			return fmt.Errorf("cannot define a field with :=")
		}
		receiverIdent, ok := lhs.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported receiver for field assignment")
		}
		symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
		}
		fieldName := lhs.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		const numSelectors = 1
		if symbol.Scope() == GlobalScope {
			if _, err := c.scopes.Emit(bytecode.OpGlobalSelSet, symbol.Index(), numSelectors); err != nil {
				return err
			}
		} else {
			if _, err := c.scopes.Emit(bytecode.OpLocalSelSet, symbol.Index(), numSelectors); err != nil {
				return err
			}
		}
		/*
			if symbol.Scope() == GlobalScope {
				if _, err := c.scopes.Emit(bytecode.OpGlobalSelSet, 1, symbol.Index()); err != nil {
					return err
				}
			} else {
				if _, err := c.scopes.Emit(bytecode.OpLocalSelSet, 1, symbol.Index()); err != nil {
					return err
				}
			}

		*/
		return nil

	default:
		return fmt.Errorf("unsupported left-hand side in assignment: %T", node.Lhs[0])
	}
}

// CompositeLit processes and compiles an abstract syntax tree (AST) CompositeLit node into bytecode instructions.
// It handles various composite literal types such as slices, structs, arrays, and maps, emitting corresponding bytecode.
// It also validates keyed and positional struct literals, generates constant keys, and handles null values when needed.
// Returns an error if compilation fails or the composite literal type is unsupported.
func (c *Declarations) CompositeLit(node *ast.CompositeLit) error {
	// Handle slice literals like []int{1, 2, 3} where the parser sets Type to nil.
	if node.Type == nil {
		// Assume it's an array/slice literal.
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpArray, len(node.Elts)); err != nil {
			return err
		}
		return nil
	}

	switch t := node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		structFields, ok := c.structsTable.GetFields(t.Name)
		if !ok {
			return fmt.Errorf("unknown composite literal type: %s", t.Name)
		}
		if len(node.Elts) > len(structFields) {
			return fmt.Errorf("too many values in positional struct literal for type '%s'", t.Name)
		}
		symbol, ok := c.scopes.SymbolResolve(t.Name)
		if !ok {
			var err error
			if symbol, err = c.scopes.SymbolDefine(t.Name); err != nil {
				return err
			}
			symbol.SetIsStruct(t.Name, true)
		}
		symbol.StructPropertyAssign(structFields)
		symbol.SetTypes([]string{t.Name})
		isKeyed := false
		if len(node.Elts) > 0 {
			if _, ok := node.Elts[0].(*ast.KeyValueExpr); ok {
				isKeyed = true
			}
		}
		if isKeyed {
			// key literal (es. Home{Name: "Alfa", Address: "Shanghai"})
			providedFields := make(map[string]ast.Expr)
			for _, elt := range node.Elts {
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					return fmt.Errorf("cannot mix keyed and unkeyed values in struct literal")
				}
				keyIdent, ok := kvExpr.Key.(*ast.Ident)
				if !ok {
					return fmt.Errorf("invalid field name in struct literal")
				}
				providedFields[keyIdent.Name] = kvExpr.Value
			}
			for idx := range symbol.Fields {
				if valueExpr, ok := providedFields[symbol.Fields[idx].name]; ok {
					symbol.Fields[idx].node = valueExpr
				}
			}
		} else {
			// positional literal (es. Home{"Alfa", 20, "Shanghai"}) ---
			for i, elt := range node.Elts {
				symbol.Fields[i].node = elt
			}
		}
		for idx := range symbol.Fields {
			fieldName := symbol.Fields[idx].name
			fieldNode := symbol.Fields[idx].node
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
			if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
				return err
			}
			if fieldNode != nil {
				if err := c.compile(fieldNode); err != nil {
					return err
				}
			} else {
				if _, err := c.scopes.Emit(bytecode.OpNull); err != nil {
					return err
				}
			}
		}
		structLen := len(symbol.Fields) * 2
		if _, err := c.scopes.Emit(bytecode.OpStruct, structLen); err != nil {
			return err
		}
		return nil
	case *ast.ArrayType:
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpArray, len(node.Elts)); err != nil {
			return err
		}
		return nil
	case *ast.MapType:
		for _, elt := range node.Elts {
			kve := elt.(*ast.KeyValueExpr)
			if err := c.compile(kve.Key); err != nil {
				return err
			}
			if err := c.compile(kve.Value); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpMap, len(node.Elts)*2); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported composite literal type: %T", node.Type)
	}
}

// KeyValueExpr processes an ast.KeyValueExpr node by compiling its Key and Value components in sequence. Returns an error if any occurs.
func (c *Declarations) KeyValueExpr(node *ast.KeyValueExpr) error {
	if err := c.compile(node.Key); err != nil {
		return err
	}
	err := c.compile(node.Value)
	return err
}

// StarExpr compiles the expression to the right of the asterisk, pushing the ObjectPointer onto the stack.
func (c *Declarations) StarExpr(node *ast.StarExpr) error {
	// Compiles the expression to the right of the asterisk (e.g. pointer 'p').
	// This will push the ObjectPointer onto the stack.
	if err := c.compile(node.X); err != nil {
		return err
	}
	_, err := c.scopes.Emit(bytecode.OpDeref)
	return err
}

// IndexExpr compiles an indexed object and its index expression, emitting VM instructions to perform the index operation.
func (c *Declarations) IndexExpr(node *ast.IndexExpr) error {
	// Compile the indexed object (e.g. myArray). This puts the array, map or slice on the stack.
	if err := c.compile(node.X); err != nil {
		return err
	}
	// Compile the index expression (e.g. i). This puts the index value on the stack.
	if err := c.compile(node.Index); err != nil {
		return err
	}
	// Emit OpIndex instruction. The VM will take the index and container from the stack and perform the access.
	_, err := c.scopes.Emit(bytecode.OpIndex)
	return err

}
