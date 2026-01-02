package compiler

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Declarations is a structure that manages compilation-related components such as scopes, constants, and imports.
type Declarations struct {
	gk              objects.IGateKeeper
	references      *tables.Constants
	constants       *tables.Constants
	scopes          *tables.Scopes
	fileSet         *token.FileSet
	imports         *Imports
	definitionTable *tables.DefinitionTable
	functionTables  *tables.FunctionTable
	initRef         map[string]int
	compile         func(node ast.Node) error
}

// NewDeclarations creates and returns a new instance of Declarations with the provided gatekeeper, tables, and imports.
func NewDeclarations(gk objects.IGateKeeper, references *tables.Constants, constants *tables.Constants, scopes *tables.Scopes, imports *Imports, definitionTable *tables.DefinitionTable, functionTables *tables.FunctionTable) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile:         nil,
		definitionTable: definitionTable,
		functionTables:  functionTables,
		imports:         imports,
	}
}

// Setup initializes the Declarations object with a FileSet and a compile function, and sets up default constants for basic types.
func (c *Declarations) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	boolConst := c.constants.AddOrGet("$_default_bool", c.gk.NewBool(objects.FrameStatic, false))
	intConst := c.constants.AddOrGet("$_default_int", c.gk.NewInt(objects.FrameStatic, 0))
	floatConst := c.constants.AddOrGet("$_default_float", c.gk.NewFloat(objects.FrameStatic, 0.0))
	charConst := c.constants.AddOrGet("$_default_char", c.gk.NewChar(objects.FrameStatic, 0))
	stringConst := c.constants.AddOrGet("$_default_string", c.gk.NewString(objects.FrameStatic, ""))
	c.initRef = map[string]int{
		"bool": boolConst,
		"int":  intConst, "int32": intConst, "int64": intConst,
		"uint": intConst, "uint32": intConst, "uint64": intConst,
		"float": floatConst, "float32": floatConst, "float64": floatConst,
		"char": charConst, "rune": charConst,
		"string": stringConst,
	}
	return nil
}

// Prepare initializes necessary resources or state within the Declarations instance for further operations.
func (c *Declarations) Prepare() error {
	return nil
}

// Compile finalizes the definition table, ensuring all definitions are resolved and ready for use.
func (c *Declarations) Compile() error {
	if err := c.definitionTable.Finalize(); err != nil {
		return err
	}
	return nil
}

// Finalize completes the initialization of Declarations and performs any final setup or cleanup required.
func (c *Declarations) Finalize() error {
	//if err := c.definitionTable.Finalize(); err != nil {
	//	return err
	//}
	return nil
}

// DeclStmt processes an `ast.DeclStmt` node by compiling its declaration and returns an error if the compilation fails.
func (c *Declarations) DeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// GenDecl processes a generic declaration node from the AST, compiling each specification within the declaration.
func (c *Declarations) GenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// TypeSpec processes a type specification node, registers types, and adds relevant fields for structs and interfaces.
func (c *Declarations) TypeSpec(node *ast.TypeSpec) error {
	typeName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(typeName); ok {
		return tables.NewCompilerError(c.fileSet, node, "type '%s' already defined", typeName)
	}

	switch t := node.Type.(type) {
	case *ast.StructType:
		c.definitionTable.StructAdd(typeName)
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				// NAME RESOLUTION (Deep Search)
				var baseStructName string
				if ident := tables.GetIdent(field.Type); ident != nil {
					baseStructName = ident.Name
				} else {
					return fmt.Errorf("unsupported type: %v", field.Type)
				}

				// REGISTRATION
				if len(field.Names) == 0 {
					// Embedding (uses the base struct name as field name)
					c.definitionTable.StructAddField(typeName, baseStructName, baseStructName, field, field.Type)
				} else {
					// Normal fields
					for _, name := range field.Names {
						c.definitionTable.StructAddField(typeName, name.Name, baseStructName, field, field.Type)
					}
				}

				//var base string
				//if ident := tables.GetIdent(field.Type); ident != nil {
				//	base = ident.Name
				//}
				//for _, name := range field.Names {
				//	c.definitionTable.StructAddField(typeName, name.Name, base, field)
				//}
			}
		}
		symbol, err := c.scopes.SymbolDefineType(typeName)
		if err != nil {
			return err
		}
		symbol.SetStruct(typeName, nil)
	case *ast.InterfaceType:
		if err := c.definitionTable.InterfaceAdd(typeName, t); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		symbol, err := c.scopes.SymbolDefineType(typeName)
		if err != nil {
			return err
		}
		symbol.SetInterface(typeName)
		symbol.SetObject(c.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
	}
	return nil
}

// ValueSpec processes an ast.ValueSpec node representing variable declarations, initializing declared variables as needed.
// It supports both initialized declarations (with values) and uninitialized declarations (with types or inferred types).
func (c *Declarations) ValueSpec(node *ast.ValueSpec) error {
	// CASE 1: Declaration has an initialization value (e.g. var x = 10)
	if len(node.Values) > 0 {
		for i, name := range node.Names {
			if i > len(node.Values)-1 {
				return tables.NewCompilerError(c.fileSet, node, "too few values for %s", name.Name)
			}
			if err := c.compile(node.Values[i]); err != nil {
				return err
			}
			symbol, err := c.scopes.SymbolDefine(name.Name)
			if err != nil {
				return err
			}

			var definedSymbol *tables.Symbol
			if node.Type != nil {
				if typeIdent := tables.GetIdent(node.Type); typeIdent != nil {
					typeSymbol, ok := c.scopes.SymbolResolve(typeIdent.Name)
					if ok {
						definedSymbol = typeSymbol
						if definedSymbol.IsInterface() {
							symbol.SetInterface(typeIdent.Name) // Pass the type name (e.g. "Printer")
							symbol.SetObject(c.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
						}
					}
				}
			}

			if definedSymbol == nil {
				if rhsIdent := tables.GetIdent(node.Values[i]); rhsIdent != nil {
					if ass, ok := c.scopes.SymbolResolve(rhsIdent.Name); ok {
						definedSymbol = ass
					}
				} else if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
					if ident := tables.GetIdent(compLit.Type); ident != nil {
						if ass, ok := c.scopes.SymbolResolve(ident.Name); ok {
							definedSymbol = ass
						}
					}
				}
			}

			if definedSymbol != nil && definedSymbol.IsStruct() {
				c.definitionTable.TypeAssign(symbol, definedSymbol.StructName())
			} else if definedSymbol != nil && definedSymbol.IsInterface() {
				inferredTypeName, _ := c.definitionTable.TypeInference(node.Values[i])
				if len(inferredTypeName) == 0 {
					return tables.NewCompilerError(c.fileSet, node, fmt.Sprintf("can't infer struct: %s", node.Values[i]))
				}
				structSymbol, ok := c.scopes.SymbolResolve(inferredTypeName)
				if !ok {
					return tables.NewCompilerError(c.fileSet, node, fmt.Sprintf("cat't resolve struct: %s", inferredTypeName))
				}
				if err = c.handleInterfaceAssign(node.Pos(), symbol, structSymbol); err != nil {
					return tables.NewCompilerError(c.fileSet, node, err.Error())
				}
			} else if inferredTypeName, _ := c.definitionTable.TypeInference(node.Values[i]); len(inferredTypeName) > 0 {
				c.definitionTable.TypeAssign(symbol, inferredTypeName)
			} else {
				symbol.SetObject(c.gk.NewString(objects.FrameStatic, symbol.Name()))
			}
			if err = c.scopes.SymbolEmitDefineAndPop(node.Pos(), symbol); err != nil {
				return err
			}
		}
		return nil
	}

	// CASE 2: Declaration has NO value (e.g. var p1 Printer)
	for _, name := range node.Names {
		symbol, err := c.scopes.SymbolDefine(name.Name)
		if err != nil {
			return err
		}
		if node.Type == nil {
			if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitDefineAndPop(node.Pos(), symbol); err != nil {
				return err
			}
			continue
		}

		// Check if the node type is a selector expression and extract its data
		if selName, selId, ok := tables.GetSelectorData(node.Type); ok {
			if ok = c.imports.EmitPackage(node.Pos(), selName, selId); ok {
				if err = c.scopes.SymbolEmitDefineAndPop(node.Pos(), symbol); err != nil {
					return err
				}
				continue
			}
		}

		if typeIdent := tables.GetIdent(node.Type); typeIdent != nil {
			if typeSymbol, ok := c.scopes.SymbolResolve(typeIdent.Name); ok {
				if typeSymbol.IsInterface() {
					c.definitionTable.InterfaceAssign(symbol, typeIdent.Name)
				} else if typeSymbol.IsStruct() {
					c.definitionTable.TypeAssign(symbol, typeIdent.Name)
				}
			}
		}

		// SymbolEmit zero value for the type. For interfaces, pointers, slices and maps, the zero value is 'nil'
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
		if err = c.scopes.SymbolEmitDefineAndPop(node.Pos(), symbol); err != nil {
			return err
		}
	}
	return nil
}

// BasicLit processes an *ast.BasicLit node and converts it into a corresponding constant object based on its type.
func (c *Declarations) BasicLit(node *ast.BasicLit) error {
	var obj objects.IObject
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		obj = c.gk.NewInt(objects.FrameStatic, val)
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		obj = c.gk.NewFloat(objects.FrameStatic, val)
	case token.CHAR:
		unquoted, err := strconv.Unquote(node.Value)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, "malformed character literal %s: %w", node.Value, err)
		}
		runes := []rune(unquoted)
		if len(runes) != 1 {
			return tables.NewCompilerError(c.fileSet, node, "character literal %s must contain exactly one character", node.Value)
		}
		obj = c.gk.NewChar(objects.FrameStatic, runes[0])
	case token.STRING:
		val, err := strconv.Unquote(node.Value)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, "malformed string literal %s: %w", node.Value, err)
		}
		obj = c.gk.NewString(objects.FrameStatic, val)
	default:
		return tables.NewCompilerError(c.fileSet, node, "unhandled literal: %s", node.Kind)
	}
	id := c.constants.Add("", obj)
	if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, id); err != nil {
		return err
	}
	return nil
}

// Ident processes an identifier node, handling special cases like "true", "false", "nil", and performing symbol resolution.
func (c *Declarations) Ident(node *ast.Ident) error {
	switch node.Name {
	case "true":
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpTrueId); err != nil {
			return err
		}
		return nil
	case "false":
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpFalseId); err != nil {
			return err
		}
		return nil
	case "nil":
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
	}
	symbol, ok := c.scopes.SymbolResolve(node.Name)
	if !ok {
		// If we can't resolve the symbol, we DON'T generate an error.
		// We assume it's a field name (e.g. 'Title' in 's.Title')
		// and that its parent node (*ast.SelectorExpr) has already handled it.
		// Simply, we don't emit any bytecode for this node.
		return nil
		//return NewCompilerError(c.FileSet, node,"[Ident] undefined variable: %s", node.Name)
	}
	if err := c.scopes.SymbolEmitGet(node.Pos(), symbol); err != nil {
		return err
	}
	return nil
}

// AssignStmt processes an assignment statement, validating and compiling both the left-hand and right-hand sides.
func (c *Declarations) AssignStmt(node *ast.AssignStmt) error {
	type rhs struct {
		node       ast.Expr
		returnType string
	}
	var rhsContainer []*rhs
	if len(node.Rhs) == 0 {
		return tables.NewCompilerError(c.fileSet, node, "invalid number of values to assign")
	}
	switch rhsType := node.Rhs[0].(type) {
	case *ast.CallExpr:
		// handle x = f(1, 2)
		if len(node.Rhs) > 1 {
			return tables.NewCompilerError(c.fileSet, node, "invalid number of values to assign")
		}
		// Compile the function call (e.g. 'f(1, 2)')
		if err := c.compile(node.Rhs[0]); err != nil {
			return err
		}
		targetIdent := tables.GetIdent(rhsType.Fun)
		if targetIdent == nil {
			return tables.NewCompilerError(c.fileSet, node, "function type not found")
		}
		var returnTypes []string
		if symbol, ok := c.scopes.SymbolResolve(targetIdent.Name); ok {
			returnTypes = symbol.ReturnTypes()
		} else {
			returnTypes = []string{tables.InterfaceDefinition}
		}
		rhsContainer = make([]*rhs, len(returnTypes))
		for idx := range rhsContainer {
			rhsContainer[idx] = &rhs{node: node.Rhs[0], returnType: returnTypes[idx]}
		}
	case *ast.TypeAssertExpr:
		// handle type assertion like val, ok := i.(ConcreteType)
		// the lhs values can be 1 or 2 (e.g. 'val' or 'val, ok')
		argCount := len(node.Lhs)
		if argCount < 1 || argCount > 2 {
			return tables.NewCompilerError(c.fileSet, node, "invalid number of values to assign")
		}
		// Compile the interface object (e.g. 'i')
		if err := c.compile(rhsType.X); err != nil {
			return err
		}
		// Extract the target type name (e.g. "User")
		targetTypeName := rhsType.Type.(*ast.Ident).Name
		constIndex := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, targetTypeName))
		hasOk := 0
		if argCount > 1 {
			hasOk = 1
		}
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpTypeAssertId, hasOk, constIndex); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		rhsContainer = make([]*rhs, len(node.Lhs))
		for idx := range node.Lhs {
			rhsContainer[idx] = &rhs{node: node.Rhs[0], returnType: ""}
		}
	default:
		if len(node.Lhs) != len(node.Rhs) {
			return tables.NewCompilerError(c.fileSet, node, "invalid number of values to assign")
		}
		rhsContainer = make([]*rhs, len(node.Rhs))
		for i := len(node.Rhs) - 1; i >= 0; i-- {
			if err := c.compile(node.Rhs[i]); err != nil {
				return err
			}
			rhsContainer[i] = &rhs{node: node.Rhs[i], returnType: ""}
		}
	}

	if len(node.Lhs) != len(rhsContainer) {
		return tables.NewCompilerError(c.fileSet, node, "assignment mismatch: %d variables but %d return values", len(node.Lhs), len(rhsContainer))
	}

	// Handle multiple assignments (e.g. x, y := 1, 2)
	for i := len(node.Lhs) - 1; i >= 0; i-- {
		pos := node.Pos()
		lhs := node.Lhs[i]
		rhsNode := rhsContainer[i].node
		rhsReturnType := rhsContainer[i].returnType
		err := c.handleVariableAssign(pos, node.Tok, rhsNode, lhs, rhsReturnType)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	return nil
}

// CompositeLit processes an AST composite literal, handling slices, arrays, structs, and maps for code compilation.
func (c *Declarations) CompositeLit(node *ast.CompositeLit) error {
	// Handle slice literals like []int{1, 2, 3} where the parser sets Type to nil.
	if node.Type == nil {
		// Assume it's an array/slice literal.
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpCreateArrayId, len(node.Elts)); err != nil {
			return err
		}
		return nil
	}
	switch val := node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		structName := val.Name
		symbol, err := c.scopes.SymbolResolveOrDefine(structName)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, fmt.Sprintf("can't resolve struct name %s: %s", structName, err.Error()))
		}
		c.definitionTable.TypeAssign(symbol, structName)
		structFields, err := c.definitionTable.StructFieldsFromLiteral(structName, node.Elts)
		if err != nil {
			return err
		}
		for _, field := range structFields {
			keyIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, field.FieldName()))
			if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, keyIdx); err != nil {
				return err
			}

			if field.FieldNode() != nil {
				if err = c.compile(field.FieldNode()); err != nil {
					return err
				}
			} else {
				ref := strings.TrimSpace(strings.ToLower(field.FieldBase())) //field.FieldKind()))
				if valIdx, ok := c.initRef[ref]; ok {
					if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, valIdx); err != nil {
						return err
					}
				} else {
					if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpNullId); err != nil {
						return err
					}
				}
			}
		}
		structNameIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, structName))
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, structNameIdx); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		structLen := len(structFields) * 2
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpCreateStructId, structLen); err != nil {
			return err
		}
		return nil
	case *ast.ArrayType:
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpCreateArrayId, len(node.Elts)); err != nil {
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
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpCreateMapId, len(node.Elts)*2); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
		if receiverIdent, ok := val.X.(*ast.Ident); ok {
			c.imports.EmitPackage(node.Pos(), receiverIdent.Name, val.Sel.Name)
		}
		return nil
	default:
		return tables.NewCompilerError(c.fileSet, node, "unsupported composite literal type: %T", node.Type)
	}
}

// Field processes the given AST field node and emits the appropriate bytecode based on the field's type.
func (c *Declarations) Field(node *ast.Field) error {
	// Analizziamo il tipo del campo per decidere cosa emettere

	switch typeNode := node.Type.(type) {
	case *ast.StarExpr:
		_, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId)
		return err
	case *ast.MapType:
		// make(map[K]V)
		_, err := c.scopes.SymbolEmit(node.Pos(), native.OpCreateMapId, 0)
		return err
	case *ast.ArrayType:
		// make([]T, 0)
		_, err := c.scopes.SymbolEmit(node.Pos(), native.OpCreateArrayId, 0)
		return err
	// 4. IDENTIFICATORE (Struct o Primitivo)
	case *ast.Ident:
		typeName := typeNode.Name

		// Verifichiamo se è una Struct (User Defined Type)
		// Nota: Qui dovresti usare la tua TypeTable/SymbolTable per sapere cos'è "typeName"
		symbol, ok := c.scopes.SymbolResolve(typeName)

		// Se è una STRUCT (e non è un puntatore, siamo nel case Ident diretto)
		// Dobbiamo generarla ricorsivamente!
		if ok && symbol.IsStruct() {
			return c.compileZeroStruct(node.Pos(), typeName)
		}

		// Altrimenti è un Primitivo (int, string, bool) o un alias
		// Cerchiamo il valore zero nei riferimenti initRef
		ref := strings.TrimSpace(strings.ToLower(typeName))
		if valIdx, ok := c.initRef[ref]; ok {
			_, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, valIdx)
			return err
		}

		// Fallback se non riconosciamo il tipo (es. Interface)
		_, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId)
		return err

	// 5. SELECTOR (es. pkg.MyType)
	case *ast.SelectorExpr:
		// Gestione simile a Ident, ma risolvendo il package
		// Per semplicità qui emettiamo Null o gestiamo la struct remota se hai accesso
		_, err := c.scopes.SymbolEmit(node.Pos(), native.OpNullId)
		return err

	default:
		return fmt.Errorf("unsupported generated type: %T", node.Type)
	}
}

// KeyValueExpr compiles a key-value pair in a composite literal by processing both the key and the value expressions.
func (c *Declarations) KeyValueExpr(node *ast.KeyValueExpr) error {
	if err := c.compile(node.Key); err != nil {
		return err
	}
	err := c.compile(node.Value)
	return err
}

// StarExpr compiles a pointer dereference expression and pushes the resulting ObjectPointer onto the stack.
func (c *Declarations) StarExpr(node *ast.StarExpr) error {
	// Compiles the expression to the right of the asterisk (e.g. pointer 'p').
	// This will push the ObjectPointer onto the stack.
	if err := c.compile(node.X); err != nil {
		return err
	}
	_, err := c.scopes.SymbolEmit(node.Pos(), native.OpDerefGetId)
	return err
}

// IndexExpr compiles an indexed expression, handling array, map, or slice access and emitting the OpIndexGet instruction.
func (c *Declarations) IndexExpr(node *ast.IndexExpr) error {
	// Compile the indexed object (e.g. myArray). This puts the array, map or slice on the stack.
	if err := c.compile(node.X); err != nil {
		return err
	}
	// Compile the index expression (e.g. i). This puts the index value on the stack.
	if err := c.compile(node.Index); err != nil {
		return err
	}
	// SymbolEmit OpIndexGet instruction. The VM will take the index and container from the stack and perform the access.
	_, err := c.scopes.SymbolEmit(node.Pos(), native.OpIndexGetId)
	return err
}

// MapType compiles an AST representation of a map type, emitting its prototype as a constant to the bytecode stack.
func (c *Declarations) MapType(node *ast.MapType) error {
	// Create a static, empty map object to serve as a type descriptor.
	// AddOrGet is used to ensure we only have one such constant.
	mapPrototype := c.gk.NewMap(objects.FrameStatic, make(map[string]objects.IObject))
	constIndex := c.constants.AddOrGet("", mapPrototype)

	// SymbolEmit the opcode to push this constant prototype onto the stack.
	// The 'make' function at runtime will inspect its TypeName.
	if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	return nil
}

// ArrayType compiles an array type declaration by creating a constant prototype and emitting the corresponding opcode.
func (c *Declarations) ArrayType(node *ast.ArrayType) error {
	// Create a static, empty map object to serve as a type descriptor.
	// AddOrGet is used to ensure we only have one such constant.
	arrayProtoType := c.gk.NewArray(objects.FrameStatic, []objects.IObject{})
	constIndex := c.constants.AddOrGet("", arrayProtoType)

	// SymbolEmit the opcode to push this constant prototype onto the stack.
	// The 'make' function at runtime will inspect its TypeName.
	if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	return nil
}

// handleInterfaceDefine binds a concrete symbol to an interface symbol, ensuring the concrete type satisfies the interface.
// It emits the concrete symbol into the current scope and performs the assignment operation.
func (c *Declarations) handleInterfaceDefine(pos token.Pos, iSymbol *tables.Symbol, concreteSymbol *tables.Symbol) error {
	if err := c.scopes.SymbolEmitGet(pos, concreteSymbol); err != nil {
		return err
	}
	if err := c.handleInterfaceAssign(pos, iSymbol, concreteSymbol); err != nil {
		return err
	}
	return nil
}

// handleInterfaceAssign validates and assigns a concrete type to an interface, ensuring the implementation is compatible.
func (c *Declarations) handleInterfaceAssign(pos token.Pos, iSymbol *tables.Symbol, concreteSymbol *tables.Symbol) error {
	interfaceName := iSymbol.InterfaceName()
	structName := concreteSymbol.StructName()
	// Compatibility check
	if !c.definitionTable.StructImplements(structName, interfaceName) {
		return fmt.Errorf("cannot use value of type %s as type %s: %s does not implement %s", structName, interfaceName, structName, interfaceName)
	}
	interfaceDesc, ok := c.definitionTable.InterfaceGet(interfaceName)
	if !ok {
		return fmt.Errorf("internal compiler error: unknown interface %s", interfaceName)
	}
	// At this point, the concrete value (struct) is already on the stack.
	// Now we need to push the (method_name, method_function) pairs for the iTable.
	for _, requiredMethod := range interfaceDesc.Methods {
		// Push method name as string constant
		methodNameConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, requiredMethod.Name))
		if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, methodNameConst); err != nil {
			return err
		}
		// Push method function
		mangledMethodName := tables.GetMangledName(structName, requiredMethod.Name)
		methodSymbol, ok := c.scopes.SymbolResolve(mangledMethodName)
		if !ok {
			return fmt.Errorf("internal compiler error: could not resolve method %s for struct %s", requiredMethod.Name, structName)
		}
		if err := c.scopes.SymbolEmitGet(pos, methodSymbol); err != nil {
			return err
		}
	}
	// SymbolEmit OpCreateInterface opcode to create the object
	if _, err := c.scopes.SymbolEmit(pos, native.OpCreateInterfaceId, len(interfaceDesc.Methods)); err != nil {
		return err
	}
	return nil
}

// handleVariableAssignNew processes variable assignment operations, supporting definition, updates, and various LHS types.
// It handles identifiers, indexed expressions, selectors, and dereference expressions, based on token type and context.
// Returns an error if the assignment is invalid or encounters a compile-time/emit-time error.
func (c *Declarations) handleVariableAssignNew(pos token.Pos, tok token.Token, rhsIn ast.Expr, lhsIn ast.Expr, inferredTypeName string) error {
	switch lhs := lhsIn.(type) {
	case *ast.Ident:
		if lhs.Name == tables.UndefinedSymbol {
			if _, err := c.scopes.SymbolEmit(pos, native.OpPopId); err != nil {
				return err
			}
			return nil
		}

		switch tok {
		case token.DEFINE:
			symbol, err := c.scopes.SymbolDefine(lhs.Name)
			if err != nil {
				return err
			}
			c.definitionTable.InferAssign(symbol, inferredTypeName, rhsIn)
			if err = c.scopes.SymbolEmitDefineAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		case token.ASSIGN:
			symbol, ok := c.scopes.SymbolResolve(lhs.Name)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] undefined variable: %s", lhs.Name)
			}
			if symbol.IsInterface() {
				var rhsName string
				if ident := tables.GetIdent(rhsIn); ident != nil {
					rhsName = ident.Name
				}
				if len(rhsName) == 0 {
					return fmt.Errorf("[handleVariableAssign] cannot assign interface to interface")
				}
				if assignedStructSymbol, ok := c.scopes.SymbolResolve(rhsName); ok && assignedStructSymbol.IsStruct() {
					if err := c.handleInterfaceAssign(pos, symbol, assignedStructSymbol); err != nil {
						return err
					}
				}
			} else {
				c.definitionTable.InferAssign(symbol, inferredTypeName, rhsIn)
			}
			if err := c.scopes.SymbolEmitSetAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		default:
			symbol, ok := c.scopes.SymbolResolve(lhs.Name)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] undefined variable: %s", lhs.Name)
			}
			if err := c.scopes.SymbolEmitGet(pos, symbol); err != nil {
				return err
			}
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			if err := c.scopes.SymbolEmitSetAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		}

	case *ast.IndexExpr:
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define variable with index assignment using %s", tok)
		case token.ASSIGN:
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if err := c.compile(lhs.Index); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		default:
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if err := c.compile(lhs.Index); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexGetId); err != nil {
				return err
			}
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if err := c.compile(lhs.Index); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		}

	case *ast.SelectorExpr:
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define a field using %s", tok)
		case token.ASSIGN:
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			fieldName := lhs.Sel.Name
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
			if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		default:
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			fieldName := lhs.Sel.Name
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
			if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexGetId); err != nil {
				return err
			}
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil { // Field (di nuovo)
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		}

	case *ast.StarExpr:
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define a dereferenced variable using %s", tok)
		default:
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if err := c.scopes.SymbolEmitAndPop(pos, native.OpDerefSetId); err != nil {
				return err
			}
			return nil
		}
	default:
		return fmt.Errorf("[handleVariableAssign] unsupported left-hand side in assignment: %T", lhs)
	}
}

// handleVariableAssign handles the assignment of a variable based on its left-hand side and operator token.
func (c *Declarations) handleVariableAssign(pos token.Pos, tok token.Token, rhsIn ast.Expr, lhsIn ast.Expr, inferredTypeName string) error {
	switch lhs := lhsIn.(type) {
	case *ast.Ident:
		if lhs.Name == tables.UndefinedSymbol {
			// if is '_', we don't create a symbol, simply discard the corresponding value from the top of the stack.
			if _, err := c.scopes.SymbolEmit(pos, native.OpPopId); err != nil {
				return err
			}
			return nil
		}
		switch tok {
		case token.DEFINE:
			symbol, err := c.scopes.SymbolDefine(lhs.Name)
			if err != nil {
				return err
			}
			c.definitionTable.InferAssign(symbol, inferredTypeName, rhsIn)
			if err = c.scopes.SymbolEmitDefineAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		case token.ASSIGN:
			symbol, ok := c.scopes.SymbolResolve(lhs.Name)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] undefined variable: %s", lhs.Name)
			}
			if symbol.IsInterface() {
				var rhsName string
				if ident := tables.GetIdent(rhsIn); ident != nil {
					rhsName = ident.Name
				}
				if len(rhsName) == 0 {
					return fmt.Errorf("[handleVariableAssign] cannot assign interface to interface")
				}
				if assignedStructSymbol, ok := c.scopes.SymbolResolve(rhsName); ok && assignedStructSymbol.IsStruct() {
					if err := c.handleInterfaceAssign(pos, symbol, assignedStructSymbol); err != nil {
						return err
					}
				}
			} else {
				c.definitionTable.InferAssign(symbol, inferredTypeName, rhsIn)
			}
			if err := c.scopes.SymbolEmitSetAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		default:
			symbol, ok := c.scopes.SymbolResolve(lhs.Name)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] undefined variable: %s", lhs.Name)
			}
			if err := c.scopes.SymbolEmitGet(pos, symbol); err != nil {
				return err
			}
			if err := c.compile(rhsIn); err != nil {
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			if err := c.scopes.SymbolEmitSetAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		}
	case *ast.IndexExpr:
		// Handles cases like 'myMap[key] = value' or 'mySlice[index] = value'
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define variable with index assignment using :=")
		case token.ASSIGN:
			tempSymbol, err := c.scopes.SymbolDefineUnique("__tmp_index_assign_rhs_")
			if err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitSetAndPop(pos, tempSymbol); err != nil {
				return err
			}
			if err = c.compile(lhs.X); err != nil { // Compiles 'm'
				return err
			}
			if err = c.compile(lhs.Index); err != nil { // Compiles "three"
				return err
			}
			if err = c.scopes.SymbolEmitGet(pos, tempSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		default:
			// Handles a[i] += v
			// 1. Get
			if err := c.compile(lhs.X); err != nil { // a
				return err
			}
			if err := c.compile(lhs.Index); err != nil { // i
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexGetId); err != nil {
				return err
			}
			// 2. Operate
			if err := c.compile(rhsIn); err != nil { // v
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			// 3. Set
			if err := c.compile(lhs.X); err != nil { // a
				return err
			}
			if err := c.compile(lhs.Index); err != nil { // i
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		}
	case *ast.SelectorExpr:
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define a field with :=")
		case token.ASSIGN:
			// Try to use the fast path for simple receivers (e.g. myVar.Field)

			if receiverIdent := tables.GetIdent(lhs.X); receiverIdent != nil {
				if symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name); ok {
					// It's a known symbol, use specific Op...SelSet opcodes
					// RHS value is already on stack, leave it there push the field name as key.
					lhsFieldName := lhs.Sel.Name
					keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, lhsFieldName))
					if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil {
						return err
					}
					// Stack is now: [..., value, "lhsFieldName"]
					const numSelectors = 1
					op := native.OpLocalIndexId
					if symbol.Scope() == tables.GlobalScope {
						op = native.OpGlobalIndexId
					}
					if _, err := c.scopes.SymbolEmit(pos, op, numSelectors, symbol.Index()); err != nil {
						return err
					}
					return nil
				}
			}

			// fallback to a general path for complex receivers (e.g. mySlice[0].Field)
			tmpSymbol, err := c.scopes.SymbolDefineUnique("__tmp_selector_assign_rhs_")
			if err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitSet(pos, tmpSymbol); err != nil {
				return err
			}
			if err = c.compile(lhs.X); err != nil {
				return err
			}
			lshFieldName := lhs.Sel.Name
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, lshFieldName))
			if _, err = c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitGet(pos, tmpSymbol); err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitAndPop(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		default:
			// Handles s.Field += v
			if err := c.compile(lhs.X); err != nil { // s
				return err
			}
			lshFieldName := lhs.Sel.Name
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, lshFieldName))
			if _, err := c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil {
				return err
			}
			if _, err := c.scopes.SymbolEmit(pos, native.OpIndexGetId); err != nil {
				return err
			}
			if err := c.compile(rhsIn); err != nil { // v
				return err
			}
			adapter, ok := BinaryAdapterFor(tok)
			if !ok {
				return fmt.Errorf("[handleVariableAssign] unsupported assignment operator: %s", tok)
			}
			if _, err := c.scopes.SymbolEmit(pos, adapter.op, adapter.arguments...); err != nil {
				return err
			}
			tmpResultSymbol, err := c.scopes.SymbolDefineUnique("__tmp_selector_assign_rhs_")
			if err != nil {
				return err
			}
			if err = c.scopes.SymbolEmitSet(pos, tmpResultSymbol); err != nil {
				return err
			}
			if err = c.compile(lhs.X); err != nil { // s
				return err
			}
			if _, err = c.scopes.SymbolEmit(pos, native.OpConstantId, keyConst); err != nil { // Field
				return err
			}
			if err = c.scopes.SymbolEmitGet(pos, tmpResultSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.SymbolEmit(pos, native.OpIndexSetId); err != nil {
				return err
			}
			return nil
		}
	case *ast.StarExpr:
		// Handles cases like '*myVar = ...'
		switch tok {
		case token.DEFINE:
			return fmt.Errorf("[handleVariableAssign] cannot define a variable with dereference")
		default:
			if err := c.compile(lhs.X); err != nil {
				return err
			}
			if err := c.scopes.SymbolEmitAndPop(pos, native.OpDerefSetId); err != nil {
				return err
			}
			return nil
		}
	default:
		return fmt.Errorf("[handleVariableAssign] unsupported left-hand side in assignment: %T", lhs)
	}
}

// compileZeroStruct generates code for initializing a struct with zero-values for all its fields based on its definition.
// It retrieves the struct's fields, emits key-value pairs for each field, and creates the struct instance.
func (c *Declarations) compileZeroStruct(pos token.Pos, structName string) error {
	// 1. Retrieve struct definition from the table
	// Note: StructFieldsFromLiteral is used "empty" to get the fields
	// Or better: use a direct method c.definitionTable.GetFields(structName)
	// We use StructFieldsFromLiteral passing nil as 'elts' if supported, or a dedicated method.
	structFields, err := c.definitionTable.StructFieldsFromLiteral(structName, nil)
	if err != nil {
		return fmt.Errorf("cannot resolve default struct for %s: %v", structName, err)
	}

	// 2. Iterate over ALL fields to generate their zero-values
	for _, field := range structFields {
		// Emit Key
		keyIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, field.FieldName()))
		if _, err = c.scopes.SymbolEmit(pos, native.OpConstantId, keyIdx); err != nil {
			return err
		}

		// Emit Value (Recursion!)
		// Here too we need to decide if it's a nested struct or a primitive
		isPtr, cont, k := field.Options()
		if !isPtr && cont == "" && k == "struct" {
			// RECURSION: A struct inside a struct
			if err = c.compileZeroStruct(pos, field.FieldBase()); err != nil {
				return err
			}
		} else {
			// Primitive / Container / Pointer -> Zero Value
			ref := strings.TrimSpace(strings.ToLower(field.FieldBase()))
			if valIdx, ok := c.initRef[ref]; ok {
				if _, err = c.scopes.SymbolEmit(pos, native.OpConstantId, valIdx); err != nil {
					return err
				}
			} else {
				if _, err = c.scopes.SymbolEmit(pos, native.OpNullId); err != nil {
					return err
				}
			}
		}
	}

	// 3. Emit struct name (Optional, depends on your OpCreateStruct)
	nameIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, structName))
	if _, err = c.scopes.SymbolEmit(pos, native.OpConstantId, nameIdx); err != nil {
		return err
	}

	// 4. Emit struct creation
	// Len * 2 because we emitted (Key, Value) for each field
	if _, err = c.scopes.SymbolEmit(pos, native.OpCreateStructId, len(structFields)*2); err != nil {
		return err
	}

	return nil
}
