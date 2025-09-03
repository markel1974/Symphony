package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Declarations is a structure responsible for managing compiler declarations and scope-related components.
// It holds references to constants, scopes, structs, and a gatekeeper for managing object lifecycle and interactions.
// The fileSet tracks source file information, and the compile function is used for compiling AST nodes.
type Declarations struct {
	gk             objects.IGateKeeper
	references     *Constants
	constants      *Constants
	scopes         *tables.Scopes
	fileSet        *token.FileSet
	imports        *Imports
	structTable    *tables.StructTable
	functionTables *tables.FunctionTable
	interfaceTable *tables.InterfaceTable
	compile        func(node ast.Node) error
}

// NewDeclarations creates and initializes a new Declarations instance with gatekeeper, constants, scopes, and structs table.
func NewDeclarations(gk objects.IGateKeeper, references *Constants, constants *Constants, scopes *tables.Scopes, imports *Imports, structsTable *tables.StructTable, functionTables *tables.FunctionTable, interfaceTable *tables.InterfaceTable) *Declarations {
	return &Declarations{
		gk: gk, references: references, constants: constants, scopes: scopes,
		compile:        nil,
		structTable:    structsTable,
		functionTables: functionTables,
		interfaceTable: interfaceTable,
		imports:        imports,
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
	typeName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(typeName); ok {
		return tables.NewCompilerError(c.fileSet, node, "type '%s' already defined", typeName)
	}

	switch t := node.Type.(type) {
	case *ast.StructType:
		if t.Fields != nil {
			for _, field := range t.Fields.List {
				var typeNameBuf bytes.Buffer
				var base = c.structTable.ExtractBaseName(field.Type)
				if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
					return tables.NewCompilerError(c.fileSet, node, "failed to resolve type for field in struct '%s'", typeName)
				}
				fieldType := typeNameBuf.String()
				for _, name := range field.Names {
					c.structTable.Add(typeName, name.Name, base, fieldType, field)
				}
			}
		}
		//TODO VERIFCARE SE NECESSARIO (VIENE INSERITO NELLA TABELLA DI STRUTTURE)
		symbol, err := c.scopes.SymbolDefine(typeName)
		if err != nil {
			return err
		}
		symbol.SetStruct(typeName, nil)
	case *ast.InterfaceType:
		// Aggiungi la definizione dell'interfaccia alla nostra nuova tabella
		if err := c.interfaceTable.Add(typeName, t); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		// Aggiungi un simbolo per il tipo interfaccia
		symbol, err := c.scopes.SymbolDefine(typeName)
		if err != nil {
			return err
		}
		// Passa il nome del tipo
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
			var assignedStructSymbol *tables.Symbol
			isInterfaceAssignment := false
			if node.Type != nil {
				if typeIdent, ok := node.Type.(*ast.Ident); ok {
					if typeSymbol, ok := c.scopes.SymbolResolve(typeIdent.Name); ok && typeSymbol.IsInterface() {
						isInterfaceAssignment = true // Only in the first block
						// Pass the type name (e.g. "Printer")
						symbol.SetInterface(typeIdent.Name)
						symbol.SetObject(c.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
					}
				}
			}
			if rhsIdent, ok := node.Values[i].(*ast.Ident); ok {
				assignedStructSymbol, _ = c.scopes.SymbolResolve(rhsIdent.Name)
			} else if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
				if ident, ok := compLit.Type.(*ast.Ident); ok {
					assignedStructSymbol, _ = c.scopes.SymbolResolve(ident.Name)
				}
			}
			if isInterfaceAssignment && assignedStructSymbol != nil && assignedStructSymbol.IsStruct() {
				if err := c.handleInterfaceAssignment(symbol, assignedStructSymbol); err != nil {
					return tables.NewCompilerError(c.fileSet, node, err.Error())
				}
			} else {
				if inferredTypeName, _ := c.structTable.TypeInference(node.Values[i]); len(inferredTypeName) > 0 {
					symbol.SetReturnTypes([]string{inferredTypeName})
					symbol.SetObject(c.gk.NewString(objects.FrameStatic, inferredTypeName+":"+symbol.Name()))
					c.structTable.BindSymbol(symbol, inferredTypeName)
				} else {
					symbol.SetObject(c.gk.NewString(objects.FrameStatic, "unknown:"+symbol.Name()))
				}
			}
			if err = c.scopes.EmitSymbolDefineAndPop(symbol); err != nil {
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
		if node.Type != nil {
			if typeIdent, ok := node.Type.(*ast.Ident); ok {
				if typeSymbol, ok := c.scopes.SymbolResolve(typeIdent.Name); ok {
					if typeSymbol.IsInterface() {
						symbol.SetInterface(typeIdent.Name)
						symbol.SetObject(c.gk.NewString(objects.FrameStatic, "interface:"+symbol.Name()))
					} else if typeSymbol.IsStruct() {
						symbol.SetReturnTypes([]string{typeSymbol.StructName()})
						symbol.SetObject(c.gk.NewString(objects.FrameStatic, typeSymbol.StructName()+":"+symbol.Name()))
						c.structTable.BindSymbol(symbol, typeSymbol.StructName())
					}
				}
			}
		}
		// Emit zero value for the type. For interfaces, pointers, slices and maps, the zero value is 'nil'
		// OpNull opcode does exactly this
		if _, err = c.scopes.Emit(bytecode.OpNull); err != nil {
			return err
		}
		// Define the variable and initialize it with 'nil'
		if err = c.scopes.EmitSymbolDefineAndPop(symbol); err != nil {
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
		return tables.NewCompilerError(c.fileSet, node, "unhandled literal: %s", node.Kind)
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
	case "nil":
		if _, err := c.scopes.Emit(bytecode.OpNull); err != nil {
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
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
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
		var targetName string
		switch funType := rhsType.Fun.(type) {
		case *ast.Ident:
			if len(funType.Name) > 0 {
				targetName = funType.Name
			}
		case *ast.SelectorExpr:
			if receiverIdent, ok := funType.X.(*ast.Ident); ok {
				targetName = receiverIdent.Name
			}
		default:
			return tables.NewCompilerError(c.fileSet, node, "function type %T not supported", funType)
		}
		if len(targetName) == 0 {
			return tables.NewCompilerError(c.fileSet, node, "function symbol name not found")
		}
		symbol, ok := c.scopes.SymbolResolve(targetName)
		if !ok {
			return tables.NewCompilerError(c.fileSet, node, "function symbol not found")
		}
		returnTypes := symbol.ReturnTypes()
		rhsContainer = make([]*rhs, len(returnTypes))
		for idx := range rhsContainer {
			rhsContainer[idx] = &rhs{node: node.Rhs[0], returnType: returnTypes[idx]}
		}
	case *ast.TypeAssertExpr:
		// handle type assertion like val, ok := i.(ConcreteType)
		// the lhs values can be 1 or 2 (e.g. 'val' or 'val, ok')
		if len(node.Lhs) < 1 || len(node.Lhs) > 2 {
			return tables.NewCompilerError(c.fileSet, node, "invalid number of values to assign")
		}
		// Compile the interface object (e.g. 'i')
		if err := c.compile(rhsType.X); err != nil {
			return err
		}
		// Extract the target type name (e.g. "User")
		targetTypeName := rhsType.Type.(*ast.Ident).Name
		constIndex := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, targetTypeName))
		if _, err := c.scopes.Emit(bytecode.OpTypeAssert, constIndex); err != nil {
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
		if err := c.handleVariableDeclaration(node.Tok, rhsContainer[i].node, node.Lhs[i], rhsContainer[i].returnType); err != nil {
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
		if _, err := c.scopes.Emit(bytecode.OpArray, len(node.Elts)); err != nil {
			return err
		}
		return nil
	}

	switch node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		t, ok := node.Type.(*ast.Ident)
		if !ok {
			return tables.NewCompilerError(c.fileSet, node, "unsupported composite literal type: %T", node)
		}
		structName := t.Name
		structFields, err := c.structTable.FieldsFromLiteral(structName, node.Elts)
		if err != nil {
			return err
		}
		for _, field := range structFields {
			keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, field.Name()))
			if _, err = c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
				return err
			}
			if field.Node() != nil {
				if err = c.compile(field.Node()); err != nil {
					return err
				}
			} else {
				if _, err = c.scopes.Emit(bytecode.OpNull); err != nil {
					return err
				}
			}
		}
		structNameIdx := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, structName))
		if _, err = c.scopes.Emit(bytecode.OpConstant, structNameIdx); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
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
	_, err := c.scopes.Emit(bytecode.OpDerefGet)
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
	_, err := c.scopes.Emit(bytecode.OpIndexGet)
	return err
}

// handleInterfaceAssignment validates and assigns a struct to a variable with an interface type, ensuring compatibility.
// It emits appropriate bytecode for the interface table setup and method bindings required for the variable's interface type.
// Returns an error if the struct does not implement the interface or if bytecode generation fails.
func (c *Declarations) handleInterfaceAssignment(variableSymbol *tables.Symbol, assignedStructSymbol *tables.Symbol) error {
	interfaceName := variableSymbol.InterfaceName()
	structName := assignedStructSymbol.StructName()
	// Compatibility check
	if !c.structTable.Implements(structName, interfaceName) {
		return fmt.Errorf("cannot use value of type %s as type %s: %s does not implement %s",
			structName, interfaceName, structName, interfaceName)
	}
	interfaceDesc, ok := c.interfaceTable.Get(interfaceName)
	if !ok {
		return fmt.Errorf("internal compiler error: unknown interface %s", interfaceName)
	}
	// At this point, the concrete value (struct) is already on the stack.
	// Now we need to push the (method_name, method_function) pairs for the iTable.
	for _, requiredMethod := range interfaceDesc.Methods {
		// Push method name as string constant
		methodNameConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, requiredMethod.Name))
		if _, err := c.scopes.Emit(bytecode.OpConstant, methodNameConst); err != nil {
			return err
		}
		// Push method function
		mangledMethodName := tables.GetMangledName(structName, requiredMethod.Name)
		methodSymbol, ok := c.scopes.SymbolResolve(mangledMethodName)
		if !ok {
			return fmt.Errorf("internal compiler error: could not resolve method %s for struct %s", requiredMethod.Name, structName)
		}
		if err := c.scopes.EmitSymbolGet(methodSymbol); err != nil {
			return err
		}
	}
	// Emit OpInterface opcode to create the object
	if _, err := c.scopes.Emit(bytecode.OpInterface, len(interfaceDesc.Methods)); err != nil {
		return err
	}
	return nil
}

// handleVariableDeclaration processes variable declarations and assignments, handling various cases like ':=' and '='.
// It resolves symbols, infers types, and manages scope and structure assignments based on the given token and expressions.
// Returns an error if the operation fails or if invalid assignments are attempted.
func (c *Declarations) handleVariableDeclaration(tok token.Token, rhsIn ast.Expr, lhsIn ast.Expr, inferredTypeName string) error {
	switch lhs := lhsIn.(type) {
	case *ast.Ident:
		if lhs.Name == tables.UndefinedSymbol {
			// if is '_', we don't create a symbol, simply discard the corresponding value from the top of the stack.
			if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
				return err
			}
			return nil
		}
		if tok == token.DEFINE {
			// Specific case for ':='
			symbol, err := c.scopes.SymbolDefine(lhs.Name)
			if err != nil {
				return err
			}
			if len(inferredTypeName) == 0 {
				inferredTypeName, _ = c.structTable.TypeInference(rhsIn)
			}
			if len(inferredTypeName) > 0 {
				symbol.SetReturnTypes([]string{inferredTypeName})
				symbol.SetObject(c.gk.NewString(objects.FrameStatic, inferredTypeName+":"+symbol.Name()))
				c.structTable.BindSymbol(symbol, inferredTypeName)
			}
			if err = c.scopes.EmitSymbolDefineAndPop(symbol); err != nil {
				return err
			}
		} else {
			// Case for normal assignment '='
			symbol, ok := c.scopes.SymbolResolve(lhs.Name)
			if !ok {
				return fmt.Errorf("[AssignStmt] undefined variable: %s", lhs.Name)
			}
			if symbol.IsInterface() {
				var rhsName string
				switch rhs := rhsIn.(type) {
				case *ast.Ident:
					rhsName = rhs.Name
				case *ast.CompositeLit:
					if ident, ok := rhs.Type.(*ast.Ident); ok {
						rhsName = ident.Name
					}
				}
				if len(rhsName) > 0 {
					if assignedStructSymbol, ok := c.scopes.SymbolResolve(rhsName); ok && assignedStructSymbol.IsStruct() {
						if err := c.handleInterfaceAssignment(symbol, assignedStructSymbol); err != nil {
							return err
						}
					}
				} else {
					return fmt.Errorf("cannot assign interface to interface")
				}
			} else {
				if len(inferredTypeName) == 0 {
					inferredTypeName, _ = c.structTable.TypeInference(rhsIn)
				}
				if len(inferredTypeName) > 0 {
					symbol.SetReturnTypes([]string{inferredTypeName})
					symbol.SetObject(c.gk.NewString(objects.FrameStatic, inferredTypeName+":"+symbol.Name()))
					c.structTable.BindSymbol(symbol, inferredTypeName)
				}
			}
			if err := c.scopes.EmitSymbolSetAndPop(symbol); err != nil {
				return err
			}
		}
		return nil
	case *ast.IndexExpr:
		// Handles cases like 'myMap[key] = value' or 'mySlice[index] = value'
		if tok == token.DEFINE {
			return fmt.Errorf("cannot define variable with index assignment using :=")
		}
		tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_assign_rhs")
		if err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSetAndPop(tempSymbol); err != nil {
			return err
		}
		if err = c.compile(lhs.X); err != nil { // Compiles 'm'
			return err
		}
		if err = c.compile(lhs.Index); err != nil { // Compiles "three"
			return err
		}
		if err = c.scopes.EmitSymbolGet(tempSymbol); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpIndexSet); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
		if tok == token.DEFINE {
			return fmt.Errorf("cannot define a field with :=")
		}
		// Try to use the fast path for simple receivers (e.g. myVar.Field)
		if receiverIdent, ok := lhs.X.(*ast.Ident); ok {
			if symbol, ok := c.scopes.SymbolResolve(receiverIdent.Name); ok {
				// It's a known symbol, use specific Op...SelSet opcodes
				// RHS value is already on stack, leave it there
				// Push the field name as key.
				fieldName := lhs.Sel.Name
				keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
				if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
					return err
				}
				// Stack is now: [..., value, "fieldName"]
				const numSelectors = 1
				scope := bytecode.OpLocalSelSet
				if symbol.Scope() == tables.GlobalScope {
					scope = bytecode.OpGlobalSelSet
				}
				if err := c.scopes.EmitAndPop(scope, symbol.Index(), numSelectors); err != nil {
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
		if err = c.scopes.EmitSymbolSet(tempSymbol); err != nil {
			return err
		}
		if err = c.compile(lhs.X); err != nil {
			return err
		}
		fieldName := lhs.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err = c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolGet(tempSymbol); err != nil {
			return err
		}
		if err = c.scopes.EmitAndPop(bytecode.OpIndexSet); err != nil {
			return err
		}
		return nil
	case *ast.StarExpr:
		// Handles cases like '*myVar = ...'
		if tok == token.DEFINE {
			return fmt.Errorf("cannot define a variable with dereference")
		}
		if err := c.compile(lhs.X); err != nil {
			return err
		}
		if err := c.scopes.EmitAndPop(bytecode.OpDerefSet); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported left-hand side in assignment: %T", lhs)
	}
}
