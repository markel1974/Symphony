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

type Declarations struct {
	gk         objects.IGateKeeper
	references *Constants
	constants  *Constants
	scopes     *Scopes
	fileSet    *token.FileSet
	structs    *Structs
	compile    func(node ast.Node) error
}

func NewDeclarations(gk objects.IGateKeeper, references *Constants, constants *Constants, scopes *Scopes, structs *Structs) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile: nil,
		structs: structs,
	}
}

// Setup initializes the `Others` instance with a compile function used for processing AST nodes.
func (c *Declarations) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// DeclStmt processes a declaration statement node by compiling its declaration content. Returns an error if compilation fails.
func (c *Declarations) DeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// GenDecl processes a general declaration node by compiling each specification within the node. It returns an error if any occur.
func (c *Declarations) GenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// TypeSpec processes a type specification node, validating and defining struct types in the current scope.
func (c *Declarations) TypeSpec(node *ast.TypeSpec) error {
	structType, isStruct := node.Type.(*ast.StructType)
	if !isStruct {
		return nil
	}
	structName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(structName); ok {
		return fmt.Errorf("type '%s' already defined", structName)
	}
	var def []*StructField
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var typeNameBuf bytes.Buffer
			var base = ExtractBaseName(field.Type)
			if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
				return fmt.Errorf("failed to resolve type for field in struct '%s'", structName)
			}
			fieldType := typeNameBuf.String()
			for _, name := range field.Names {
				// here we could add a check for duplicate fields.
				def = append(def, NewStructProperty(name.Name, base, fieldType, field))
			}
		}
	}
	c.structs.Add(structName, def)
	return nil
}

// ValueSpec processes a ValueSpec node to handle variable declarations and assignments within a given scope.
func (c *Declarations) ValueSpec(node *ast.ValueSpec) error {
	// handles 'var x = 10'
	for i, name := range node.Names {
		if i > len(node.Values)-1 {
			return fmt.Errorf("too few values for %s", name.Name)
		}
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}

		isStruct := false

		// 3. Inferenza del tipo, ora coerente con la nuova logica
		var assignedTypeNames []string
		if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
			if ident, ok := compLit.Type.(*ast.Ident); ok {
				if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
					assignedTypeNames = []string{typeSymbol.Name()}
					isStruct = true
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

		symbol, err := c.scopes.SymbolDefine(name.Name, UnknownScope, isStruct)
		if err != nil {
			return err
		}

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

// BasicLit processes an AST BasicLit node and emits the corresponding literal into the current scope.
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

// Ident processes an identifier node, resolving its symbol in the current scope and emitting a symbol get operation.
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

// AssignStmt processes an assignment statement by compiling the right-hand side and resolving variable symbols.
// It also updates the type information for symbols or emits appropriate bytecode for assignments.
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
				symbol, err = c.scopes.SymbolDefine(ident.Name, UnknownScope, false) // Anche qui, l'inferenza del tipo potrebbe essere migliorata
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
			isStruct := false
			var assignedTypeName []string

			switch rhs := node.Rhs[0].(type) {
			case *ast.BasicLit:
				//nothing to do
			case *ast.CompositeLit: // es. MyStruct{...}
				baseName := ExtractBaseName(rhs.Type)
				if len(baseName) > 0 {
					assignedTypeName = []string{baseName}
					isStruct = true
				}
			case *ast.CallExpr: // es. NewStruct()
				if ident, ok := rhs.Fun.(*ast.Ident); ok {
					if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
						// Assumiamo il primo tipo di ritorno
						typeName := funcSymbol.Types()[0]
						assignedTypeName = []string{typeName}
						// Verifichiamo se il tipo restituito è uno struct
						if typeSymbol, ok := c.scopes.SymbolResolve(typeName); ok && typeSymbol.IsStruct() {
							isStruct = true
						}
					}
				}
			case *ast.UnaryExpr: // es. &MyStruct{}
				if rhs.Op == token.AND {
					if compLit, ok := rhs.X.(*ast.CompositeLit); ok {
						if ident, ok := compLit.Type.(*ast.Ident); ok {
							if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
								isStruct = true
								assignedTypeName = []string{typeSymbol.Name()}
							}
						}
					}
				}
			default:
				return fmt.Errorf("unsupported right-hand side for assignment: %T", rhs)
			}
			// Definiamo il simbolo usando il valore 'isStruct' appena calcolato
			symbol, err = c.scopes.SymbolDefine(name, UnknownScope, isStruct)
			if err != nil {
				return err
			}
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
			if _, err := c.scopes.Emit(bytecode.OpSetSelGlobal, symbol.Index(), numSelectors); err != nil {
				return err
			}
		} else {
			if _, err := c.scopes.Emit(bytecode.OpSetSelLocal, symbol.Index(), numSelectors); err != nil {
				return err
			}
		}
		/*
			if symbol.Scope() == GlobalScope {
				if _, err := c.scopes.Emit(bytecode.OpSetSelGlobal, 1, symbol.Index()); err != nil {
					return err
				}
			} else {
				if _, err := c.scopes.Emit(bytecode.OpSetSelLocal, 1, symbol.Index()); err != nil {
					return err
				}
			}

		*/
		return nil

	default:
		return fmt.Errorf("unsupported left-hand side in assignment: %T", node.Lhs[0])
	}
}

// CompositeLit processes the given composite literal node and compiles it into bytecode representation.
// Handles struct, array, and map literals by resolving types, validating fields, and emitting appropriate instructions.
// Returns an error if the composite literal type is unsupported or if any validation or compilation step fails.
// CompositeLit processes the given composite literal node and compiles it into bytecode representation.
// Handles struct, array, and map literals by resolving types, validating fields, and emitting appropriate instructions.
// Returns an error if the composite literal type is unsupported or if any validation or compilation step fails.
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
		structFields := c.structs.Get(t.Name)
		if structFields == nil {
			return fmt.Errorf("unknown composite literal type: %s", t.Name)
		}
		if len(node.Elts) > len(structFields) {
			return fmt.Errorf("too many values in positional struct literal for type '%s'", t.Name)
		}
		symbol, ok := c.scopes.SymbolResolve(t.Name)
		if !ok {
			var err error
			if symbol, err = c.scopes.SymbolDefine(t.Name, UnknownScope, true); err != nil {
				return err
			}
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

// KeyValueExpr processes a KeyValueExpr node and emits the corresponding literal into the current scope.
func (c *Declarations) KeyValueExpr(node *ast.KeyValueExpr) error {
	if err := c.compile(node.Key); err != nil {
		return err
	}
	err := c.compile(node.Value)
	return err
}

func (c *Declarations) StarExpr(node *ast.StarExpr) error {
	// Compiles the expression to the right of the asterisk (e.g. pointer 'p').
	// This will push the ObjectPointer onto the stack.
	if err := c.compile(node.X); err != nil {
		return err
	}
	_, err := c.scopes.Emit(bytecode.OpDeref)
	return err
}

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
