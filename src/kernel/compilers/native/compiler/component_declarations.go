package compiler

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Declarations is a structure responsible for managing compiler declarations and scope-related components.
// It holds references to constants, scopes, structs, and a gatekeeper for managing object lifecycle and interactions.
// The fileSet tracks source file information, and the compile function is used for compiling AST nodes.
type Declarations struct {
	gk          objects.IGateKeeper
	references  *Constants
	constants   *Constants
	scopes      *Scopes
	fileSet     *token.FileSet
	imports     *Imports
	structTable *StructTable
	compile     func(node ast.Node) error
}

// NewDeclarations creates and initializes a new Declarations instance with gatekeeper, constants, scopes, and structs table.
func NewDeclarations(gk objects.IGateKeeper, references *Constants, constants *Constants, scopes *Scopes, imports *Imports, structsTable *StructTable) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile:     nil,
		structTable: structsTable,
		imports:     imports,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *Declarations) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes the ControlFlow structure, ensuring it is ready for subsequent compilation tasks and operations.
func (c *Declarations) Prepare() error {
	return nil
}

// Compile compiles the AST nodes using the configured compile function and returns an error if the process fails.
func (c *Declarations) Compile() error {
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
		return NewCompilerError(c.fileSet, node, "type '%s' already defined", structName)
	}
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var typeNameBuf bytes.Buffer
			var base = c.structTable.ExtractBaseName(field.Type)
			if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
				return NewCompilerError(c.fileSet, node, "failed to resolve type for field in struct '%s'", structName)
			}
			fieldType := typeNameBuf.String()
			for _, name := range field.Names {
				c.structTable.Add(structName, name.Name, base, fieldType, field)
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
			return NewCompilerError(c.fileSet, node, "too few values for %s", name.Name)
		}
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}

		structName := ""

		// 3. Inferenza del tipo, ora coerente con la nuova logica
		var assignedTypeNames []string
		if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
			if ident, ok := compLit.Type.(*ast.Ident); ok {
				if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
					structName = ident.Name
					assignedTypeNames = []string{typeSymbol.Name()}
				}
			}
		} else if callExpr, ok := node.Values[i].(*ast.CallExpr); ok {
			var funcName string
			if ident, isIdent := callExpr.Fun.(*ast.Ident); isIdent {
				funcName = ident.Name
			}
			if funcName != "" {
				if funcSymbol, ok := c.scopes.SymbolResolve(funcName); ok {
					returnTypes := funcSymbol.ReturnTypes()
					if len(returnTypes) != 1 {
						return NewCompilerError(c.fileSet, node, "assignment mismatch: 'var' declaration expects 1 value, but function %s returns %d", funcName, len(returnTypes))
					}
					assignedTypeNames = []string{returnTypes[0]}
				}
			}
		}
		symbol, err := c.scopes.SymbolDefine(name.Name)
		if err != nil {
			return err
		}
		if len(structName) > 0 {
			if err = c.structTable.AssignSymbol(symbol, structName, assignedTypeNames); err != nil {
				return err
			}
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
		return NewCompilerError(c.fileSet, node, "unhandled literal: %s", node.Kind)
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
		//return NewCompilerError(c.FileSet, node,"[Ident] undefined variable: %s", node.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	return nil
}

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
					funcReturnTypes = funcSymbol.ReturnTypes()
				}
			}
		}
		if len(node.Lhs) != len(funcReturnTypes) {
			return NewCompilerError(c.fileSet, node, "assignment mismatch: %d variables but %d return values", len(node.Lhs), len(funcReturnTypes))
		}
		for i := len(node.Lhs) - 1; i >= 0; i-- {
			lhs := node.Lhs[i]
			ident, ok := lhs.(*ast.Ident)
			if !ok {
				return NewCompilerError(c.fileSet, node, "unsupported multiple assignment to type %T", lhs)
			}
			var symbol *Symbol
			if node.Tok == token.DEFINE {
				var err error
				if symbol, err = c.scopes.SymbolDefine(ident.Name); err != nil {
					return err
				}
			} else {
				var found bool
				if symbol, found = c.scopes.SymbolResolve(ident.Name); !found {
					return NewCompilerError(c.fileSet, node, "undefined variable: %s", ident.Name)
				}
			}
			// Inferenza completa del tipo per ogni variabile.
			inferredTypeName := funcReturnTypes[i]
			if err := c.structTable.AssignSymbol(symbol, inferredTypeName, []string{inferredTypeName}); err != nil {
				return err
			}
			// Emettiamo l'opcode corretto in base a ':=' o '='.
			if node.Tok == token.DEFINE {
				if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
					return err
				}
			} else {
				if err := c.scopes.EmitSymbolSet(symbol); err != nil {
					return err
				}
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
		var err error
		if node.Tok == token.DEFINE { // Caso specifico per ':='
			symbol, err = c.scopes.SymbolDefine(name)
			if err != nil {
				return err
			}
			if structName, assignedTypeName, _ := c.structTable.Inference(node.Rhs[0]); len(structName) > 0 {
				if err = c.structTable.AssignSymbol(symbol, structName, assignedTypeName); err != nil {
					return err
				}
			}
			if err = c.scopes.EmitSymbolDefine(symbol); err != nil {
				return err
			}
		} else { // Caso per l'assegnazione normale '='
			var ok bool
			symbol, ok = c.scopes.SymbolResolve(name)
			if !ok {
				return NewCompilerError(c.fileSet, node, "[AssignStmt] undefined variable: %s", name)
			}
			if err = c.scopes.EmitSymbolSet(symbol); err != nil {
				return err
			}
		}
		if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr: // es. myStruct.Field = ...
		if node.Tok == token.DEFINE {
			return NewCompilerError(c.fileSet, node, "cannot define a field with :=")
		}
		receiverIdent, ok := lhs.X.(*ast.Ident)
		if !ok {
			return NewCompilerError(c.fileSet, node, "unsupported receiver for field assignment")
		}
		symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return NewCompilerError(c.fileSet, node, "undefined variable: %s", receiverIdent.Name)
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
		return nil
	case *ast.StarExpr: // Gestisce casi come '*myVar = ...'
		if node.Tok == token.DEFINE {
			return NewCompilerError(c.fileSet, node, "cannot define a variable with dereference")
		}
		if err := c.compile(lhs.X); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpDerefSet); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		return nil
	default:
		return NewCompilerError(c.fileSet, node, "unsupported left-hand side in assignment: %T", node.Lhs[0])
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

	switch node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		structFields, err := c.structTable.FieldsFromLiteral(node)
		if err != nil {
			return err
		}
		for _, field := range structFields {
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, field.name))
			if _, err = c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
				return err
			}
			if field.node != nil {
				if err = c.compile(field.node); err != nil {
					return err
				}
			} else {
				if _, err = c.scopes.Emit(bytecode.OpNull); err != nil {
					return err
				}
			}
		}
		structLen := len(structFields) * 2
		if _, err = c.scopes.Emit(bytecode.OpStruct, structLen); err != nil {
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
		return NewCompilerError(c.fileSet, node, "unsupported composite literal type: %T", node.Type)
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
