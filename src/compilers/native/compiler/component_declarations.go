package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Declarations is a structure responsible for managing compiler declarations and scope-related components.
// It holds references to constants, scopes, structs, and a gatekeeper for managing object lifecycle and interactions.
// The fileSet tracks source file information, and the compile function is used for compiling AST nodes.
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

// NewDeclarations creates and initializes a new Declarations instance with gatekeeper, constants, scopes, and structs table.
func NewDeclarations(gk objects.IGateKeeper, references *tables.Constants, constants *tables.Constants, scopes *tables.Scopes, imports *Imports, definitionTable *tables.DefinitionTable, functionTables *tables.FunctionTable) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile:         nil,
		definitionTable: definitionTable,
		functionTables:  functionTables,
		imports:         imports,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
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
	typeName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(typeName); ok {
		return tables.NewCompilerError(c.fileSet, node, "type '%s' already defined", typeName)
	}

	switch t := node.Type.(type) {
	case *ast.StructType:
		c.definitionTable.StructAdd(typeName)
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				var typeNameBuf bytes.Buffer
				var base string
				if ident := tables.GetIdent(field.Type); ident != nil {
					base = ident.Name
				}
				if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
					return tables.NewCompilerError(c.fileSet, node, "failed to resolve type for field in struct '%s'", typeName)
				}
				fieldType := typeNameBuf.String()
				for _, name := range field.Names {
					c.definitionTable.StructAddField(typeName, name.Name, base, fieldType, field)
				}
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

// ValueSpec processes a value specification, handling variable declarations with and without initialization values.
// It performs type checking, symbol definition, and code generation for variable bindings within the current scope.
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
			if err = c.scopes.EmitSymbolDefineAndPop(node.Pos(), symbol); err != nil {
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
			if _, err = c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
				return err
			}
			if err = c.scopes.EmitSymbolDefineAndPop(node.Pos(), symbol); err != nil {
				return err
			}
			continue
		}

		// Check if the node type is a selector expression and extract its data
		if selName, selId, ok := tables.GetSelectorData(node.Type); ok {
			if ok = c.imports.EmitPackage(node.Pos(), selName, selId); ok {
				if err = c.scopes.EmitSymbolDefineAndPop(node.Pos(), symbol); err != nil {
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

		// Emit zero value for the type. For interfaces, pointers, slices and maps, the zero value is 'nil'
		if _, err = c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolDefineAndPop(node.Pos(), symbol); err != nil {
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
	if _, err := c.scopes.Emit(node.Pos(), native.OpConstantId, id); err != nil {
		return err
	}
	return nil
}

// Ident processes an identifier node and emits bytecode if the identifier corresponds to a symbol or keyword in the scope.
func (c *Declarations) Ident(node *ast.Ident) error {
	switch node.Name {
	case "true":
		if _, err := c.scopes.Emit(node.Pos(), native.OpTrueId); err != nil {
			return err
		}
		return nil
	case "false":
		if _, err := c.scopes.Emit(node.Pos(), native.OpFalseId); err != nil {
			return err
		}
		return nil
	case "nil":
		if _, err := c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
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
	if err := c.scopes.EmitSymbolGet(node.Pos(), symbol); err != nil {
		return err
	}
	return nil
}

// AssignStmt processes assignment statements in the AST, handling variable assignment, declaration, and struct inference.
// It supports single assignments, multiple return from functions, and advanced cases like selector and pointer assignments.
// Returns an error if there are issues during statement compilation, undefined variables, or invalid assignment targets.
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
		if _, err := c.scopes.Emit(node.Pos(), native.OpTypeAssertId, hasOk, constIndex); err != nil {
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
		if _, err := c.scopes.Emit(node.Pos(), native.OpCreateArrayId, len(node.Elts)); err != nil {
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
			keyIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, field.Name()))
			if _, err = c.scopes.Emit(node.Pos(), native.OpConstantId, keyIdx); err != nil {
				return err
			}
			if field.Node() != nil {
				if err = c.compile(field.Node()); err != nil {
					return err
				}
			} else {
				ref := strings.TrimSpace(strings.ToLower(field.Kind()))
				if valIdx, ok := c.initRef[ref]; ok {
					if _, err = c.scopes.Emit(node.Pos(), native.OpConstantId, valIdx); err != nil {
						return err
					}
				} else {
					if _, err = c.scopes.Emit(node.Pos(), native.OpNullId); err != nil {
						return err
					}
				}
			}
		}
		structNameIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, structName))
		if _, err = c.scopes.Emit(node.Pos(), native.OpConstantId, structNameIdx); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		structLen := len(structFields) * 2
		if _, err = c.scopes.Emit(node.Pos(), native.OpCreateStructId, structLen); err != nil {
			return err
		}
		return nil
	case *ast.ArrayType:
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(node.Pos(), native.OpCreateArrayId, len(node.Elts)); err != nil {
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
		if _, err := c.scopes.Emit(node.Pos(), native.OpCreateMapId, len(node.Elts)*2); err != nil {
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
	_, err := c.scopes.Emit(node.Pos(), native.OpDerefGetId)
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
	// Emit OpIndexGet instruction. The VM will take the index and container from the stack and perform the access.
	_, err := c.scopes.Emit(node.Pos(), native.OpIndexGetId)
	return err
}

// MapType handles map type nodes, primarily when they appear as arguments
// to built-in functions like 'make'. It pushes a prototype Map object onto
// the stack, which acts as a type descriptor for the runtime function.
func (c *Declarations) MapType(node *ast.MapType) error {
	// Create a static, empty map object to serve as a type descriptor.
	// AddOrGet is used to ensure we only have one such constant.
	mapPrototype := c.gk.NewMap(objects.FrameStatic, make(map[string]objects.IObject))
	constIndex := c.constants.AddOrGet("", mapPrototype)

	// Emit the opcode to push this constant prototype onto the stack.
	// The 'make' function at runtime will inspect its TypeName.
	if _, err := c.scopes.Emit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	return nil
}

// ArrayType handles map type nodes, primarily when they appear as arguments
// to built-in functions like 'make'. It pushes a prototype Map object onto
// the stack, which acts as a type descriptor for the runtime function.
func (c *Declarations) ArrayType(node *ast.ArrayType) error {
	// Create a static, empty map object to serve as a type descriptor.
	// AddOrGet is used to ensure we only have one such constant.
	arrayProtoType := c.gk.NewArray(objects.FrameStatic, []objects.IObject{})
	constIndex := c.constants.AddOrGet("", arrayProtoType)

	// Emit the opcode to push this constant prototype onto the stack.
	// The 'make' function at runtime will inspect its TypeName.
	if _, err := c.scopes.Emit(node.Pos(), native.OpConstantId, constIndex); err != nil {
		return err
	}
	return nil
}

// handleInterfaceDefine handles the process of defining an interface, ensuring proper symbol emission and assignment.
func (c *Declarations) handleInterfaceDefine(pos token.Pos, iSymbol *tables.Symbol, concreteSymbol *tables.Symbol) error {
	if err := c.scopes.EmitSymbolGet(pos, concreteSymbol); err != nil {
		return err
	}
	if err := c.handleInterfaceAssign(pos, iSymbol, concreteSymbol); err != nil {
		return err
	}
	return nil
}

// handleInterfaceAssign validates and assigns a struct to a variable with an interface type, ensuring compatibility.
// It emits appropriate bytecode for the interface table setup and method bindings required for the variable's interface type.
// Returns an error if the struct does not implement the interface or if bytecode generation fails.
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
		if _, err := c.scopes.Emit(pos, native.OpConstantId, methodNameConst); err != nil {
			return err
		}
		// Push method function
		mangledMethodName := tables.GetMangledName(structName, requiredMethod.Name)
		methodSymbol, ok := c.scopes.SymbolResolve(mangledMethodName)
		if !ok {
			return fmt.Errorf("internal compiler error: could not resolve method %s for struct %s", requiredMethod.Name, structName)
		}
		if err := c.scopes.EmitSymbolGet(pos, methodSymbol); err != nil {
			return err
		}
	}
	// Emit OpCreateInterface opcode to create the object
	if _, err := c.scopes.Emit(pos, native.OpCreateInterfaceId, len(interfaceDesc.Methods)); err != nil {
		return err
	}
	return nil
}

// handleVariableAssign processes variable declarations and assignments, handling various cases like ':=' and '='.
// It resolves symbols, infers types, and manages scope and structure assignments based on the given token and expressions.
// Returns an error if the operation fails or if invalid assignments are attempted.
func (c *Declarations) handleVariableAssign(pos token.Pos, tok token.Token, rhsIn ast.Expr, lhsIn ast.Expr, inferredTypeName string) error {
	switch lhs := lhsIn.(type) {
	case *ast.Ident:
		if lhs.Name == tables.UndefinedSymbol {
			// if is '_', we don't create a symbol, simply discard the corresponding value from the top of the stack.
			if _, err := c.scopes.Emit(pos, native.OpPopId); err != nil {
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
			//must be a concrete value
			c.definitionTable.InferAssign(symbol, inferredTypeName, rhsIn)
			if err = c.scopes.EmitSymbolDefineAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		case token.ASSIGN:
			// Case for normal assignment '='
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
			if err := c.scopes.EmitSymbolSetAndPop(pos, symbol); err != nil {
				return err
			}
			return nil
		default:
			// TODO see BinaryAdapterFor token.AND_ASSIGN, token.OR_ASSIGN, etc.
			return fmt.Errorf("[handleVariableAssign] invalid token %s for variable assignment", tok)
		}
	case *ast.IndexExpr:
		// Handles cases like 'myMap[key] = value' or 'mySlice[index] = value'
		if tok == token.DEFINE {
			return fmt.Errorf("[handleVariableAssign] cannot define variable with index assignment using :=")
		}
		tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_assign_rhs")
		if err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSetAndPop(pos, tempSymbol); err != nil {
			return err
		}
		if err = c.compile(lhs.X); err != nil { // Compiles 'm'
			return err
		}
		if err = c.compile(lhs.Index); err != nil { // Compiles "three"
			return err
		}
		if err = c.scopes.EmitSymbolGet(pos, tempSymbol); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(pos, native.OpIndexSetId); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
		if tok == token.DEFINE {
			return fmt.Errorf("[handleVariableAssign] cannot define a field with :=")
		}
		// Try to use the fast path for simple receivers (e.g. myVar.Field)
		if receiverIdent, ok := lhs.X.(*ast.Ident); ok {
			if symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name); ok {
				// It's a known symbol, use specific Op...SelSet opcodes
				// RHS value is already on stack, leave it there
				// Push the field name as key.
				fieldName := lhs.Sel.Name
				keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
				if _, err := c.scopes.Emit(pos, native.OpConstantId, keyConst); err != nil {
					return err
				}
				// Stack is now: [..., value, "fieldName"]
				const numSelectors = 1
				op := native.OpLocalIndexId
				if symbol.Scope() == tables.GlobalScope {
					op = native.OpGlobalIndexId
				}
				if _, err := c.scopes.Emit(pos, op, numSelectors, symbol.Index()); err != nil {
					return err
				}
				return nil
			}
		}
		// fallback to a general path for complex receivers (e.g. mySlice[0].Field)
		tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_assign_rhs")
		if err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSet(pos, tempSymbol); err != nil {
			return err
		}
		if err = c.compile(lhs.X); err != nil {
			return err
		}
		fieldName := lhs.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err = c.scopes.Emit(pos, native.OpConstantId, keyConst); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolGet(pos, tempSymbol); err != nil {
			return err
		}
		if err = c.scopes.EmitAndPop(pos, native.OpIndexSetId); err != nil {
			return err
		}
		return nil
	case *ast.StarExpr:
		// Handles cases like '*myVar = ...'
		if tok == token.DEFINE {
			return fmt.Errorf("[handleVariableAssign] cannot define a variable with dereference")
		}
		if err := c.compile(lhs.X); err != nil {
			return err
		}
		if err := c.scopes.EmitAndPop(pos, native.OpDerefSetId); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("[handleVariableAssign] unsupported left-hand side in assignment: %T", lhs)
	}
}
