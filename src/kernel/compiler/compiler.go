package compiler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strconv"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/compiler/sdk"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// maxScope defines the maximum allowable depth for compilation scopes to prevent excessive recursion or memory use.
const (
	maxScope = 1024
)

// Compiler represents a structure to manage the compilation process, including scopes and associated token file sets.
type Compiler struct {
	factory    *objects.GateKeeper
	loader     *sdk.Loader
	scopes     *Scopes
	constants  *Constants
	references *Constants
	fileSet    *token.FileSet
}

// New creates and returns a new instance of Compiler with initialized scopes using a standard library loader.
func New(factory *objects.GateKeeper) *Compiler {
	loader := sdk.NewLoader(factory)
	op := bytecode.NewOpcodes(factory)
	c := &Compiler{
		factory:    factory,
		loader:     loader,
		constants:  NewConstants(),
		references: NewConstants(),
		scopes:     NewScopes(factory, op),
	}
	return c
}

// Compile parses the provided source file and compiles it into bytecode. Returns compiled bytecode or an error.
func (c *Compiler) Compile(filename string, source any) error {
	for idx := 0; idx < c.loader.BuiltinLen(); idx++ {
		bi := c.loader.Builtin(idx)
		if bi == nil {
			return fmt.Errorf("builtin %d not found", idx)
		}
		c.constants.Add(bi.Name(), bi)
		//c.scopes.SymbolDefine(bi.Name(), GlobalScope, false)
	}
	c.fileSet = token.NewFileSet()
	astFile, err := parser.ParseFile(c.fileSet, filename, source, 0)
	if err != nil {
		return err
	}
	if err = c.compile(astFile); err != nil {
		return err
	}
	return nil
}

// Constants retrieves a slice of IObject containing all constants stored in the current compiler scopes.
func (c *Compiler) Constants() []objects.IObject {
	return c.constants.Retrieve()
}

// References retrieves a list of IObject references from the current compiler scope.
func (c *Compiler) References() []objects.IObject {
	return c.references.Retrieve()
}

// Global retrieves and returns all global objects from the root scope and any objects tracked by references.
func (c *Compiler) Global() []objects.IObject {
	var ret []objects.IObject
	for _, obj := range c.scopes.root.symbols {
		ret = append(ret, c.factory.NewString(objects.FrameStatic, obj.name))
	}
	return ret
}

// Print writes the content of the internal scopes to the provided writer, typically for debugging or inspection.
func (c *Compiler) Print(writer io.Writer) {
	_, _ = fmt.Fprintf(writer, "----- Constants -----")
	c.constants.Print(writer)
	_, _ = fmt.Fprintf(writer, "----- References -----")
	c.references.Print(writer)
	c.scopes.Print(writer)
}

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (c *Compiler) compile(in ast.Node) error {
	var err error = nil

	switch node := in.(type) {
	case *ast.File:
		err = c.doFile(node)
	case *ast.DeclStmt:
		err = c.doDeclStmt(node)
	case *ast.TypeSpec:
		err = c.doTypeSpec(node)
	case *ast.GenDecl: // for `var` and `const` which are handled by AssignStmt
		err = c.doGenDecl(node)
	case *ast.ValueSpec: // handles 'var x = 10'
		err = c.doValueSpec(node)
	case *ast.BlockStmt:
		err = c.doBlockStmt(node)
	case *ast.ExprStmt:
		err = c.doExprStmt(node)
	case *ast.AssignStmt:
		err = c.doAssignStmt(node)
	case *ast.IfStmt:
		err = c.doIfStmt(node)
	case *ast.RangeStmt:
		err = c.doRangeStmt(node)
	case *ast.ForStmt:
		err = c.doForStmt(node)
	case *ast.IncDecStmt:
		err = c.doIncDecStmt(node)
	case *ast.BinaryExpr:
		err = c.doBinaryExpr(node)
	case *ast.UnaryExpr:
		err = c.doUnaryExpr(node)
	case *ast.BasicLit:
		err = c.doBasicLit(node)
	case *ast.Ident:
		err = c.doIdent(node)
	case *ast.CompositeLit:
		err = c.doCompositeLit(node)
	case *ast.FuncDecl:
		err = c.doFuncDecl(node)
	case *ast.CallExpr:
		err = c.doCallExpr(node)
	case *ast.ReturnStmt:
		err = c.doReturnStmt(node)
	case *ast.SelectorExpr:
		err = c.doSelectorExpr(node)
	case *ast.ImportSpec:
		err = c.doImportSpec(node)
	default:
		err = fmt.Errorf("unsupported expression type: %T", node)
	}
	return err
}

// doFile processes an AST file, organizing and compiling its declarations into separate categories.
// It first categorizes declarations into imports, types, functions, and others for targeted compilation.
// Import declarations are processed to resolve module dependencies.
// Type declarations such as structs are compiled to define types used in the program.
// Functions and methods are pre-defined to map their symbols and return types before full compilation.
// Non-function code is compiled after imports and type declarations are processed.
// Finally, the bodies of all functions and methods are compiled for execution.
func (c *Compiler) doFile(node *ast.File) error {
	var importDecls []ast.Decl
	var typeDecls []ast.Decl
	var otherDecls []ast.Decl
	var funcDecls []*ast.FuncDecl

	// step 1: Separate declarations by category
	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			funcDecls = append(funcDecls, d)
		case *ast.GenDecl:
			if d.Tok == token.IMPORT {
				importDecls = append(importDecls, d)
			} else if d.Tok == token.TYPE {
				typeDecls = append(typeDecls, d)
			} else {
				otherDecls = append(otherDecls, d)
			}
		default:
			otherDecls = append(otherDecls, d)
		}
	}

	// step 2: compile all import definitions
	for _, decl := range importDecls {
		if err := c.compile(decl); err != nil {
			return err
		}
	}

	// step 3: compile all type definitions (structs)
	for _, decl := range typeDecls {
		if err := c.compile(decl); err != nil {
			return err
		}
	}

	// step 4: pre-define all functions AND methods, including their return types.
	funcIndexes := make(map[string]int)

	for _, fn := range funcDecls {
		var fnSymbol *Symbol
		fnName := ""
		if fn.Recv != nil && len(fn.Recv.List) > 0 { // Method
			recvTypeIdent := GetIdent(fn.Recv.List[0])
			if recvTypeIdent == nil {
				return fmt.Errorf("unsupported method receiver type")
			}
			fnName = GetMangledName(recvTypeIdent.Name, fn.Name.Name)
			symbol, ok := c.scopes.SymbolResolve(recvTypeIdent.Name)
			if !ok || !symbol.IsStruct() {
				return fmt.Errorf("unknown type '%s' for method receiver", recvTypeIdent.Name)
			}
			fnSymbol = symbol
		} else { // Function
			fnName = fn.Name.Name
			fnSymbol = c.scopes.SymbolDefine(fnName, UnknownScope, false)
		}
		if fnSymbol == nil {
			return fmt.Errorf("unknown function '%s'", fn.Name.Name)
		}

		// function pre-definition
		placeholder := c.factory.NewFuncCompiled(objects.FrameStatic, fnName, nil, 0, 0, false, nil, nil)
		funcIndexes[fnName] = c.constants.Add(fnName, placeholder)
		receiverNames, err := GetReceivers(fn.Type.Results)
		if err != nil {
			return err
		}
		if len(receiverNames) > 0 {
			fnSymbol.SetTypes(receiverNames)
		}
	}

	// step 5: compile all other non-function code
	for _, decl := range otherDecls {
		if err := c.compile(decl); err != nil {
			return err
		}
	}

	// step 6: compile the actual bodies of functions and methods
	for _, fn := range funcDecls {
		var objName string
		var mangledName string
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			recvTypeIdent := GetIdent(fn.Recv.List[0])
			objName = recvTypeIdent.Name
			mangledName = GetMangledName(objName, fn.Name.Name)
		} else {
			mangledName = fn.Name.Name
		}
		idx, ok := funcIndexes[mangledName]
		if !ok {
			return fmt.Errorf("unknown function '%s'", mangledName)
		}
		if err := c.compileFuncBody(fn, objName, mangledName, idx); err != nil {
			return err
		}
	}
	return nil
}

// compileFuncBody compiles the body of a function declaration and generates the necessary bytecode instructions.
func (c *Compiler) compileFuncBody(node *ast.FuncDecl, objName string, mangledName string, constIndex int) error {
	if err := c.scopes.Enter(objName); err != nil {
		return err
	}
	// Aggiunge il ricevitore e i parametri come variabili locali.
	if node.Recv != nil && len(node.Recv.List) > 0 {
		for _, p := range node.Recv.List {
			for _, name := range p.Names {
				c.scopes.SymbolDefine(name.Name, UnknownScope, false)
			}
		}
	}
	for _, p := range node.Type.Params.List {
		for _, name := range p.Names {
			c.scopes.SymbolDefine(name.Name, UnknownScope, false)
		}
	}
	if err := c.compile(node.Body); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	if scope.LastInstruction() == nil || scope.LastInstruction().opcode != bytecode.OpReturn {
		if _, err = c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
	}
	freeSymbols := c.scopes.SymbolFreeConvert()
	numFree := c.scopes.SymbolFreeCount()
	nLocals := c.scopes.SymbolCount()
	code, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	nParams := 0
	if paramL := node.Type.Params; paramL != nil && paramL.List != nil {
		for _, field := range paramL.List {
			nParams += len(field.Names)
		}
	}
	if node.Recv != nil && len(node.Recv.List) > 0 {
		nParams++
	}
	compiledFn := c.factory.NewFuncCompiled(objects.FrameStatic, mangledName, code, nLocals, nParams, false, nil, freeSymbols)
	if err = c.constants.SetIndex(constIndex, compiledFn); err != nil {
		return err
	}
	if node.Recv == nil {
		if _, err = c.scopes.Emit(bytecode.OpClosure, constIndex, numFree); err != nil {
			return err
		}
		symbol, _ := c.scopes.SymbolResolve(node.Name.Name)
		if err = c.scopes.EmitSymbolSet(symbol); err != nil {
			return err
		}
	}
	return nil
}

// doAssignStmt processes an assignment statement by compiling the right-hand side and resolving variable symbols.
// It also updates the type information for symbols or emits appropriate bytecode for assignments.
func (c *Compiler) doAssignStmt(node *ast.AssignStmt) error {
	//TODO UNIFY function-only and multi assignment

	// multi assignment function-only (es. x, y := f())
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
				symbol = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
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

	//single assignment
	if err := c.compile(node.Rhs[0]); err != nil {
		return err
	}
	//inferenza del tipo
	var assignedTypeName []string
	switch rhs := node.Rhs[0].(type) {
	case *ast.CompositeLit: // check for variable assignment
		if ident, ok := rhs.Type.(*ast.Ident); ok {
			if typeSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && typeSymbol.IsStruct() {
				assignedTypeName = []string{typeSymbol.Name()}
			}
		}
	case *ast.CallExpr: // check for function call assignment
		if ident, ok := rhs.Fun.(*ast.Ident); ok {
			if funcSymbol, ok := c.scopes.SymbolResolve(ident.Name); ok && len(funcSymbol.Types()) > 0 {
				assignedTypeName = []string{funcSymbol.Types()[0]}
			}
		}
	}

	switch lhs := node.Lhs[0].(type) {
	case *ast.Ident:
		name := lhs.Name
		var symbol *Symbol
		if node.Tok == token.DEFINE {
			symbol = c.scopes.SymbolDefine(name, UnknownScope, false)
		} else {
			var ok bool
			symbol, ok = c.scopes.SymbolResolve(name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", name)
			}
		}
		// Updates the symbol type in both cases (:= and =)
		if len(assignedTypeName) > 0 {
			symbol.SetTypes(assignedTypeName)
		}
		if err := c.scopes.EmitSymbolSet(symbol); err != nil {
			return err
		}
		// OpPop is only needed for simple variable assignments because OpSetLocal/Global don't clean up the stack.
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		return nil
	case *ast.SelectorExpr:
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
		keyConst := c.constants.AddOrGet("", c.factory.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if symbol.Scope() == GlobalScope {
			if _, err := c.scopes.Emit(bytecode.OpSetSelGlobal, symbol.Index(), 1); err != nil {
				return err
			}
		} else {
			if _, err := c.scopes.Emit(bytecode.OpSetSelLocal, symbol.Index(), 1); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported left-hand side in assignment: %T", node.Lhs[0])
	}
}

// doCallExpr processes a single CallExpr node and generates bytecode for function or method calls.
// It resolves the function or method being called, handles arguments, and emits corresponding bytecode.
// It also manages nested function calls by pre-evaluating them and storing results in temporary variables.
// Returns an error if the call expression contains invalid or unresolved references.
func (c *Compiler) doCallExpr(node *ast.CallExpr) error {
	// Step 1: Resolve the function and its arguments for analysis. This part doesn't emit bytecode yet.
	fnOpType := bytecode.OpConstant
	fnIndex := -1
	fnName := ""
	var fnArgs []ast.Expr

	if selExpr, isSelector := node.Fun.(*ast.SelectorExpr); isSelector {
		// Method call
		receiverIdent, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported receiver for selector expression: %T", selExpr.X)
		}
		receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
		}
		if receiverSymbol.Scope() == ImportScope {
			fnOpType = bytecode.OpReferences
			fnName = GetMangledName(receiverIdent.Name, selExpr.Sel.Name)
			found := false
			fnIndex, found = c.references.Get(fnName)
			if !found {
				attrArray := c.factory.NewArray(objects.FrameStatic, []objects.IObject{c.factory.NewString(objects.FrameStatic, receiverIdent.Name), c.factory.NewString(objects.FrameStatic, selExpr.Sel.Name)})
				fnIndex = c.references.Add(fnName, attrArray)
			}
			fnArgs = node.Args
		} else if len(receiverSymbol.Types()) > 0 { // struct method
			structTypeName := receiverSymbol.Types()[0]
			typeSymbol, ok := c.scopes.SymbolResolve(structTypeName)
			if !ok {
				return fmt.Errorf("undefined type: %s", structTypeName)
			}
			methodName := selExpr.Sel.Name
			fnName = GetMangledName(typeSymbol.Name(), methodName)
			fnIndex, ok = c.constants.Get(fnName)
			if !ok {
				return fmt.Errorf("undefined method '%s' for type '%s'", methodName, typeSymbol.Name())
			}
			fnArgs = append(fnArgs, selExpr.X) // The receiver is the first argument
			fnArgs = append(fnArgs, node.Args...)
		} else {
			return fmt.Errorf("cannot call method on untyped variable or undefined package '%s'", receiverSymbol.Name())
		}
	} else {
		//Function call
		ident, ok := node.Fun.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported function call: %T", node.Fun)
		}
		fnName = ident.Name
		fnIndex, ok = c.constants.Get(fnName)
		if !ok {
			return fmt.Errorf("undefined function: %s", ident.Name)
		}
		fnArgs = node.Args
	}
	if fnIndex < 0 {
		return fmt.Errorf("could not resolve function index for '%s'", fnName)
	}

	// 2 pass logic
	// Pass 1: Pre-evaluate nested function calls and store their results in temporary variables.
	// We use a map to link an argument expression to the temporary symbol that holds its result.
	tempSymbolMap := make(map[ast.Expr]*Symbol)
	for _, arg := range fnArgs {
		if call, isCall := arg.(*ast.CallExpr); isCall {
			// This argument is a nested function call.
			// 1.1. Compile the nested call. Its result will be on the stack.
			if err := c.compile(call); err != nil {
				return err
			}
			// 1.2. Create a unique temporary local variable.
			tempSymbol := c.scopes.SymbolDefineUnique("__temp_call", LocalScope, false)
			tempSymbolMap[arg] = tempSymbol // Store the symbol for the second pass
			// 1.3. Emit code to store the result into the temp variable.
			// This generates OpSetLocal, correctly storing the result.
			if err := c.scopes.EmitSymbolSet(tempSymbol); err != nil {
				return err
			}
			// 1.4. Pop the result from the stack to keep it clean for the next operations.
			// This is crucial as OpSetLocal only peeks at the stack value.
			if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
				return err
			}
		}
	}

	// Step 2: Push the main function object itself onto the stack.
	// This ensures the stack layout is [function, ...args]
	if _, err := c.scopes.Emit(fnOpType, fnIndex); err != nil {
		return err
	}
	// Pass 3: Push all final arguments for the main call onto the stack.
	for _, arg := range fnArgs {
		if tempSymbol, ok := tempSymbolMap[arg]; ok {
			// This was a nested call; load its pre-computed result from the temporary variable.
			if err := c.scopes.EmitSymbolGet(tempSymbol); err != nil {
				return err
			}
		} else {
			// This is a regular argument; compile it directly to push its value.
			if err := c.compile(arg); err != nil {
				return err
			}
		}
	}
	// Step 4: Emit the final OpCall instruction.
	if _, err := c.scopes.Emit(bytecode.OpCall, len(fnArgs), 0); err != nil {
		return err
	}
	return nil
}

// doDeclStmt processes a declaration statement node by compiling its declaration content. Returns an error if compilation fails.
func (c *Compiler) doDeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// doBlockStmt compiles each statement in the provided block and returns an error if any compilation step fails.
func (c *Compiler) doBlockStmt(node *ast.BlockStmt) error {
	for _, s := range node.List {
		if err := c.compile(s); err != nil {
			return err
		}
	}
	return nil
}

// doExprStmt compiles an expression statement and emits a pop operation to discard its result.
func (c *Compiler) doExprStmt(node *ast.ExprStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// doTypeSpec processes a type specification node, validating and defining struct types in the current scope.
func (c *Compiler) doTypeSpec(node *ast.TypeSpec) error {
	structType, isStruct := node.Type.(*ast.StructType)
	if !isStruct {
		return nil
	}
	structName := node.Name.Name
	if _, ok := c.scopes.SymbolResolve(structName); ok {
		return fmt.Errorf("type '%s' already defined", structName)
	}
	var fields []*FieldDef
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var typeNameBuf bytes.Buffer
			if err := printer.Fprint(&typeNameBuf, c.fileSet, field.Type); err != nil {
				return fmt.Errorf("failed to resolve type for field in struct '%s'", structName)
			}
			fieldType := typeNameBuf.String()
			for _, name := range field.Names {
				// here we could add a check for duplicate fields.
				fields = append(fields, NewFieldDef(name.Name, fieldType, nil))
			}
		}
	}
	symbol := c.scopes.SymbolDefine(structName, UnknownScope, true)
	//symbol.Scope() == GlobalScope
	//c.constants.Add("", c.factory.NewStruct(objects.FrameStatic, map[string]objects.IObject{}))
	//symbol.SetScope(TypeScope)
	//symbol := c.scopes.SymbolDefine(structName, TypeScope)
	symbol.Fields = fields
	return nil
}

// doCompositeLit processes the given composite literal node and compiles it into bytecode representation.
// Handles struct, array, and map literals by resolving types, validating fields, and emitting appropriate instructions.
// Returns an error if the composite literal type is unsupported or if any validation or compilation step fails.
func (c *Compiler) doCompositeLit(node *ast.CompositeLit) error {
	switch t := node.Type.(type) {
	case *ast.Ident:
		// struct literal (es. MyStruct{...})
		symbol, ok := c.scopes.SymbolResolve(t.Name)
		if !ok || !symbol.IsStruct() {
			return fmt.Errorf("unknown composite literal type: %s", t.Name)
		}
		if len(node.Elts) > len(symbol.Fields) {
			return fmt.Errorf("too many values in positional struct literal for type '%s'", symbol.Name())
		}
		for idx := range symbol.Fields {
			symbol.Fields[idx].SetNode(nil)
		}
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
				if valueExpr, ok := providedFields[symbol.Fields[idx].Name()]; ok {
					symbol.Fields[idx].SetNode(valueExpr)
				}
			}
		} else {
			// positional literal (es. Home{"Alfa", 20, "Shanghai"}) ---
			for i, elt := range node.Elts {
				symbol.Fields[i].SetNode(elt)
			}
		}
		for idx := range symbol.Fields {
			fieldName := symbol.Fields[idx].Name()
			fieldNode := symbol.Fields[idx].Node()
			keyConst := c.constants.AddOrGet("", c.factory.NewString(objects.FrameStatic, fieldName))
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

// doGenDecl processes a general declaration node by compiling each specification within the node. It returns an error if any occur.
func (c *Compiler) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// doValueSpec processes a ValueSpec node to handle variable declarations and assignments within a given scope.
func (c *Compiler) doValueSpec(node *ast.ValueSpec) error {
	// handles 'var x = 10'
	for i, name := range node.Names {
		if i > len(node.Values)-1 {
			return fmt.Errorf("too few values for %s", name.Name)
		}
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}
		symbol := c.scopes.SymbolDefine(name.Name, UnknownScope, false)

		// 3. Inferenza del tipo, ora coerente con la nuova logica
		var assignedTypeNames []string
		if compLit, ok := node.Values[i].(*ast.CompositeLit); ok {
			if ident, ok := compLit.Type.(*ast.Ident); ok {
				if typeSymbol, isType := c.scopes.SymbolResolve(ident.Name); isType && typeSymbol.IsStruct() {
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
					returnTypes := funcSymbol.Types()
					if len(returnTypes) != 1 {
						return fmt.Errorf("assignment mismatch: 'var' declaration expects 1 value, but function %s returns %d", funcName, len(returnTypes))
					}
					assignedTypeNames = []string{returnTypes[0]}
				}
			}
		}

		if len(assignedTypeNames) > 0 {
			symbol.SetTypes(assignedTypeNames)
		}

		// 4. Emette bytecode per assegnare il valore dalla cima dello stack alla variabile.
		if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
			return err
		}
		// 5. Pulisce lo stack dal valore ora che è stato assegnato.
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	return nil
}

// doIfStmt compiles an if statement, handling both 'then' and optional 'else' blocks with associated bytecode instructions.
func (c *Compiler) doIfStmt(node *ast.IfStmt) error {
	if err := c.compile(node.Cond); err != nil {
		return err
	}
	// emit conditional jump with temporary address
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	// compile 'then' block
	if err = c.compile(node.Body); err != nil {
		return err
	}
	// if there's an 'else' block, emit jump to skip it
	jumpToEndPos := 0
	if node.Else != nil {
		jumpToEndPos, err = c.scopes.Emit(bytecode.OpJump, 9999)
		if err != nil {
			return err
		}
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	// update conditional jump address
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, scope.InstructionsLen()); err != nil {
		return err
	}
	// compile 'else' block if it exists
	if node.Else != nil {
		if err = c.compile(node.Else); err != nil {
			return err
		}
		scope, err = c.scopes.Current()
		if err != nil {
			return err
		}
		// update jump address to skip else
		if err = c.scopes.ChangeOperand(jumpToEndPos, scope.InstructionsLen()); err != nil {
			return err
		}
	}
	return nil
}

// doIncDecStmt handles increment and decrement statements for identifiers, updating the corresponding variables and cleaning the stack.
func (c *Compiler) doIncDecStmt(node *ast.IncDecStmt) error {
	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return fmt.Errorf("unsupported IncDec statement for type %T", node.X)
	}
	symbol, ok := c.scopes.SymbolResolve(ident.Name)
	if !ok {
		return fmt.Errorf("undefined variable: %s", ident.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	// adds constant '1' to the stack
	constIndex := c.constants.Add("", c.factory.NewInt(objects.FrameStatic, 1))
	if _, err := c.scopes.Emit(bytecode.OpConstant, constIndex); err != nil {
		return err
	}
	if node.Tok == token.INC {
		if _, err := c.scopes.Emit(bytecode.OpBinary, int(objects.OperatorAdd)); err != nil {
			return err
		}
	} else if node.Tok == token.DEC {
		if _, err := c.scopes.Emit(bytecode.OpBinary, int(objects.OperatorSub)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported IncDec token: %s", node.Tok)
	}
	if err := c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	// the increment/decrement operation leaves the result on the stack. Since this is a statement, we need to clean this value.
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// doForStmt compiles a for loop statement, including initialization, condition, post-iteration, and body execution.
func (c *Compiler) doForStmt(node *ast.ForStmt) error {
	if node.Init != nil {
		if err := c.compile(node.Init); err != nil {
			return err
		}
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	loopStartPos := scope.InstructionsLen()
	// compiles condition (e.g. x < 10)
	if node.Cond != nil {
		if err = c.compile(node.Cond); err != nil {
			return err
		}
	} else {
		// if no condition is provided, it's an infinite loop - for simplicity emit 'true'
		if _, err = c.scopes.Emit(bytecode.OpTrue); err != nil {
			return err
		}
	}
	// emits a conditional jump to exit the loop if condition is false
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	if err = c.compile(node.Body); err != nil {
		return err
	}
	// compiles post-iteration statement (e.g. x++)
	if node.Post != nil {
		if err = c.compile(node.Post); err != nil {
			return err
		}
	}
	// emits an unconditional jump to return to condition start
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}
	scope, err = c.scopes.Current()
	if err != nil {
		return err
	}
	// updates (back-patching) conditional jump address (OpJumpFalsy) to point to loop end
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	return nil
}

// doRangeStmt compiles a RangeStmt node into bytecode, handling iterator initialization, key/value assignment, and looping logic.
func (c *Compiler) doRangeStmt(node *ast.RangeStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	iteratorSymbol := c.scopes.SymbolDefineUnique("__iterator", UnknownScope, false)
	if _, err = c.scopes.Emit(bytecode.OpIteratorInit, iteratorSymbol.Index()); err != nil {
		return err
	}
	var keySymbol, valueSymbol *Symbol
	if node.Key != nil {
		if ident, ok := node.Key.(*ast.Ident); ok && ident.Name != "_" {
			keySymbol = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
		}
	}
	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != "_" {
			valueSymbol = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
		}
	}
	// Loop start
	loopStartPos := scope.InstructionsLen()
	// Check if there are more elements, passing the iterator index
	if _, err := c.scopes.Emit(bytecode.OpIteratorNext, iteratorSymbol.Index()); err != nil {
		return err
	}
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	// Assign values and clean operand stack
	if valueSymbol != nil {
		if _, err = c.scopes.Emit(bytecode.OpIteratorValue, iteratorSymbol.Index()); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSet(valueSymbol); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	if keySymbol != nil {
		if _, err = c.scopes.Emit(bytecode.OpIteratorKey, iteratorSymbol.Index()); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpIteratorValue, iteratorSymbol.Index()); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSet(keySymbol); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}
	if err = c.compile(node.Body); err != nil {
		return err
	}
	// Jump to start
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}
	// Back-patching of exit jump
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	return nil
}

// doReturnStmt compiles a return statement, handling cases for both void and value returns, and emits corresponding bytecode.
func (c *Compiler) doReturnStmt(node *ast.ReturnStmt) error {
	if len(node.Results) == 0 {
		if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
		return nil
	}
	for _, result := range node.Results {
		if err := c.compile(result); err != nil {
			return err
		}
	}
	if _, err := c.scopes.Emit(bytecode.OpReturn, len(node.Results)); err != nil {
		return err
	}
	return nil
}

// doBinaryExpr processes a binary expression node, compiling both operands and emitting the corresponding binary operation.
func (c *Compiler) doBinaryExpr(node *ast.BinaryExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if err := c.compile(node.Y); err != nil {
		return err
	}
	z, ok := BinaryAdapterFor(node.Op)
	if !ok {
		return fmt.Errorf("unhandled binary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// doUnaryExpr compiles a unary expression by evaluating the operand and applying the specified unary operator.
// It handles special cases for the address-of operator '&', ensuring correct pointer behavior based on operand type.
// Emits appropriate bytecode instructions for each unary operation or returns an error on unsupported cases.
func (c *Compiler) doUnaryExpr(node *ast.UnaryExpr) error {
	if node.Op == token.AND {
		switch operand := node.X.(type) {
		case *ast.Ident:
			// literal (es. '&h').
			symbol, ok := c.scopes.SymbolResolve(operand.Name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", operand.Name)
			}
			switch symbol.Scope() {
			case LocalScope:
				if _, err := c.scopes.Emit(bytecode.OpGetLocalPtr, symbol.Index()); err != nil {
					return err
				}
			case FreeScope:
				if _, err := c.scopes.Emit(bytecode.OpGetFreePtr, symbol.Index()); err != nil {
					return err
				}
			default:
				return fmt.Errorf("cannot take the address of a global variable")
			}
		case *ast.CompositeLit:
			// literal (es. '&Home{...}').
			if err := c.compile(operand); err != nil {
				return err
			}
			tempSymbol := c.scopes.SymbolDefineUnique("__temp_struct", LocalScope, false)
			if err := c.scopes.EmitSymbolDefine(tempSymbol); err != nil {
				return err
			}
			if _, err := c.scopes.Emit(bytecode.OpGetLocalPtr, tempSymbol.Index()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot take the address of %T", node.X)
		}
		return nil
	}
	// existing logic for other unary operators (e.g. '!', '-', '^')
	if err := c.compile(node.X); err != nil {
		return err
	}
	z, ok := UnaryAdapterFor(node.Op)
	if !ok {
		return fmt.Errorf("unhandled unary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// doBasicLit processes an AST BasicLit node and emits the corresponding literal into the current scope.
func (c *Compiler) doBasicLit(node *ast.BasicLit) error {
	var obj objects.IObject
	switch node.Kind {
	case token.INT:
		val, _ := strconv.ParseInt(node.Value, 0, 64)
		obj = c.factory.NewInt(objects.FrameStatic, val)
	case token.FLOAT:
		val, _ := strconv.ParseFloat(node.Value, 64)
		obj = c.factory.NewFloat(objects.FrameStatic, val)
	case token.STRING:
		val, _ := strconv.Unquote(node.Value)
		obj = c.factory.NewString(objects.FrameStatic, val)
	default:
		return fmt.Errorf("unhandled literal: %s", node.Kind)
	}
	id := c.constants.Add("", obj)
	if _, err := c.scopes.Emit(bytecode.OpConstant, id); err != nil {
		return err
	}
	return nil
}

// doIdent processes an identifier node, resolving its symbol in the current scope and emitting a symbol get operation.
func (c *Compiler) doIdent(node *ast.Ident) error {
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
		return fmt.Errorf("undefined variable: %s", node.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	return nil
}

// doImportSpec handles an import specification by defining the imported module name in the current scope.
func (c *Compiler) doImportSpec(node *ast.ImportSpec) error {
	moduleName := node.Path.Value
	c.scopes.SymbolDefine(strings.Trim(moduleName, "\"'"), ImportScope, false)
	return nil
}

// doSelectorExpr processes a selector expression, resolving fields, methods, or package attributes.
// It distinguishes between struct field accesses and package-level selectors.
// Emits appropriate bytecode instructions for each case or returns an error if unsupported.
func (c *Compiler) doSelectorExpr(node *ast.SelectorExpr) error {
	// analyze the left-hand side of the dot to determine if it's a variable or a package.
	receiverIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// currently not handling complex cases like a[0].field
		return fmt.Errorf("unsupported receiver for selector expression: %T", node.X)
	}
	receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
	if !ok {
		return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
	}
	if receiverSymbol.Scope() == ImportScope {
		cacheKey := GetMangledName(receiverIdent.Name, node.Sel.Name)
		nameIndex, found := c.references.Get(cacheKey)
		if !found {
			attrArray := c.factory.NewArray(objects.FrameStatic, []objects.IObject{c.factory.NewString(objects.FrameStatic, receiverIdent.Name), c.factory.NewString(objects.FrameStatic, node.Sel.Name)})
			nameIndex = c.references.Add(cacheKey, attrArray)
		}
		if _, err := c.scopes.Emit(bytecode.OpReferences, nameIndex); err != nil {
			return err
		}
		//} else if len(receiverSymbol.Object()) > 0 { // struct
	} else if receiverSymbol.IsStruct() { // struct
		if err := c.compile(node.X); err != nil {
			return err
		}
		fieldName := node.Sel.Name
		keyConst := c.constants.AddOrGet("", c.factory.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpIndex); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unsupported selector expression: %T", node.X)
}

// doFuncDecl processes the function declaration node and compiles its structure into the appropriate bytecode.
func (c *Compiler) doFuncDecl(_ *ast.FuncDecl) error {
	return nil
}
