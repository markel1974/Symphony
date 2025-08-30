package compiler

import (
	"go/ast"
	"go/token"

	tables2 "github.com/markel1974/c64emu/src/kernel/compilers/native/tables"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Functions is a collection that manages a list of function descriptions.
type Functions struct {
	gk             objects.IGateKeeper
	constants      *Constants
	scopes         *tables2.Scopes
	imports        *Imports
	declarations   *Declarations
	functionTable  *tables2.FunctionTable
	structTable    *tables2.StructTable
	interfaceTable *tables2.InterfaceTable
	fileSet        *token.FileSet
	compile        func(node ast.Node) error
}

// NewFunctions initializes and returns a new Functions instance.
func NewFunctions(gk objects.IGateKeeper, constants *Constants, scopes *tables2.Scopes, imports *Imports, declarations *Declarations, structTable *tables2.StructTable, functionTable *tables2.FunctionTable, interfaceTable *tables2.InterfaceTable) *Functions {
	return &Functions{
		gk:             gk,
		constants:      constants,
		scopes:         scopes,
		imports:        imports,
		declarations:   declarations,
		structTable:    structTable,
		functionTable:  functionTable,
		interfaceTable: interfaceTable,
		compile:        nil,
	}
}

// Setup initializes the `Functions` instance with a compile function used for processing AST nodes.
func (c *Functions) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Declare adds a function description derived from the provided function declaration to the Functions container.
func (c *Functions) Declare(funcDecl *ast.FuncDecl) {
	c.functionTable.Add(funcDecl)
}

// Prepare iterates through all function descriptions in the container and prepares their bodies for compilation.
func (c *Functions) Prepare() error {
	for idx := 0; idx < c.functionTable.Len(); idx++ {
		fd, err := c.functionTable.Get(idx)
		if err != nil {
			return err
		}
		if err = c.funcBodyPrepare(fd); err != nil {
			return err
		}
	}
	return nil
}

// Compile iterates over all function descriptions in the container and compiles their bodies. Returns an error if any fails.
func (c *Functions) Compile() error {
	for idx := 0; idx < c.functionTable.Len(); idx++ {
		fd, err := c.functionTable.Get(idx)
		if err != nil {
			return err
		}
		if err = c.funcBodyCompile(fd); err != nil {
			return err
		}
	}
	return nil
}

// funcBodyPrepare prepares the function body by processing its declaration and receiver, and defining its symbol in the scope.
// It returns an error if the receiver type is unsupported or if symbol definition fails.
func (c *Functions) funcBodyPrepare(fd *tables2.FunctionDescription) error {
	node := fd.FuncDecl
	var err error
	if fd.ReturnTypes, err = tables2.GetReceivers(node.Type.Results); err != nil {
		return err
	}
	for _, p := range node.Type.Params.List {
		for _, name := range p.Names {
			fd.InputParams = append(fd.InputParams, name.Name)
		}
	}
	if node.Recv != nil && len(node.Recv.List) > 0 {
		for _, p := range node.Recv.List {
			for _, name := range p.Names {
				fd.StructReceivers = append(fd.StructReceivers, name.Name)
			}
		}
		recvTypeIdent := tables2.GetIdent(node.Recv.List[0])
		if recvTypeIdent == nil {
			return tables2.NewCompilerError(c.fileSet, node, "unsupported method receiver type")
		}
		if !c.structTable.Has(recvTypeIdent.Name) {
			return tables2.NewCompilerError(c.fileSet, node, "undefined type '%s' for method receiver", recvTypeIdent.Name)
		}
		fd.Name = tables2.GetMangledName(recvTypeIdent.Name, node.Name.Name)
		fd.StructName = recvTypeIdent.Name
	} else {
		fd.Name = node.Name.Name
		fd.StructName = ""
	}
	//function symbol placeholder (this is not the real function, it's just a placeholder to be able to compile the body)
	placeHolder, err := c.scopes.SymbolDefine(fd.Name)
	if err != nil {
		return err
	}
	if len(fd.StructName) > 0 {
		if err = c.structTable.AssignSymbol(placeHolder, fd.StructName, nil); err != nil {
			return err
		}
	}
	//placeHolder.SetStruct(fd.StructName)
	return nil
}

// funcBodyCompile compiles the body of a function declaration and generates the necessary bytecode instructions.
func (c *Functions) funcBodyCompile(fd *tables2.FunctionDescription) error {
	node := fd.FuncDecl
	if err := c.scopes.Enter(fd.Name); err != nil {
		return err
	}

	// --- INIZIO MODIFICA ---
	// Sostituiamo la vecchia logica con chiamate dirette a SymbolsFromParameters,
	// che processa correttamente i tipi di struct e interfacce.

	// Definisci i simboli per i ricevitori del metodo (se presenti)
	if node.Recv != nil && len(node.Recv.List) > 0 {
		if _, err := c.functionTable.SymbolsFromParameters(node.Recv); err != nil {
			return err
		}
	}

	// Definisci i simboli per i parametri di input
	if _, err := c.functionTable.SymbolsFromParameters(node.Type.Params); err != nil {
		return err
	}
	// --- FINE MODIFICA ---

	// Il vecchio codice che iterava su fd.StructReceivers e fd.InputParams viene rimosso
	// perché ora è gestito in modo più completo dalla logica qui sopra.

	if err := c.compile(node.Body); err != nil {
		return err
	}

	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	if scope.LastInstruction() == nil || scope.LastInstruction().Opcode() != bytecode.OpReturn {
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

	nParams := c.functionTable.CountParams(node.Type.Params)

	if node.Recv != nil && len(node.Recv.List) > 0 {
		nParams++
	}

	fnSymbol, ok := c.scopes.SymbolRebuildScope(fd.Name, tables2.UnknownScope)
	if !ok {
		return tables2.NewCompilerError(c.fileSet, node, "undefined function: %s", fd.Name)
	}
	compiledFn := c.gk.NewFuncCompiled(objects.FrameStatic, fd.Name, code, nLocals, nParams, false, nil, freeSymbols)
	fnSymbol.SetObject(compiledFn)
	fnSymbol.SetReturnTypes(fd.ReturnTypes)

	if node.Recv == nil && !c.scopes.IsRootScope() {
		if _, err = c.scopes.Emit(bytecode.OpClosure, fnSymbol.Index(), numFree); err != nil {
			return err
		}
		symbol, _ := c.scopes.SymbolResolve(node.Name.Name)
		if err = c.scopes.EmitSymbolSet(symbol); err != nil {
			return err
		}
	}
	return nil
}

// CallExpr processes a single CallExpr node and generates bytecode for function or method calls.
// It resolves the function or method being called, handles arguments, and emits corresponding bytecode.
// It also manages nested function calls by pre-evaluating them and storing results in temporary variables.
// Returns an error if the call expression contains invalid or unresolved references.
func (c *Functions) CallExpr(node *ast.CallExpr) error {
	// Step 1: Pre-evaluate nested function calls
	tempSymbolMap := make(map[ast.Expr]*tables2.Symbol)
	for _, arg := range node.Args {
		if call, isCall := arg.(*ast.CallExpr); isCall {
			if err := c.compile(call); err != nil {
				return err
			}
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_call")
			if err != nil {
				return err
			}
			tempSymbol.SetScope(tables2.LocalScope)
			tempSymbolMap[arg] = tempSymbol
			if err = c.scopes.EmitSymbolSetAndPop(tempSymbol); err != nil {
				return err
			}
		}
	}

	var finalArgs []ast.Expr

	switch fun := node.Fun.(type) {
	case *ast.Ident:
		// Handle a simple function call (unchanged)
		symbol, ok := c.scopes.SymbolResolve(fun.Name)
		if !ok {
			if c.imports.Emit(fun.Name, "") {
				finalArgs = node.Args
				break
			}
			return tables2.NewCompilerError(c.fileSet, node, "undefined function: %s", fun.Name)
		}
		if err := c.scopes.EmitSymbolGet(symbol); err != nil {
			return err
		}
		finalArgs = node.Args

	case *ast.SelectorExpr:
		receiverIdent, ok := fun.X.(*ast.Ident)
		if !ok {
			return tables2.NewCompilerError(c.fileSet, node, "unsupported receiver for selector expression: %T", fun.X)
		}

		// Path 1: Package function (e.g., fmt.Println)
		if c.imports.Emit(receiverIdent.Name, fun.Sel.Name) {
			finalArgs = node.Args
			break
		}

		receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return tables2.NewCompilerError(c.fileSet, node, "undefined variable: %s", receiverIdent.Name)
		}

		// Path 2: Method call on an interface
		if receiverSymbol.IsInterface() {
			// 2a. Load the interface variable (the receiver) onto the stack.
			// The VM will use this object to find the iTable.
			if err := c.compile(fun.X); err != nil {
				return err
			}
			// 2b. Load all call arguments onto the stack.
			for _, arg := range node.Args {
				if tempSymbol, ok := tempSymbolMap[arg]; ok {
					if err := c.scopes.EmitSymbolGet(tempSymbol); err != nil {
						return err
					}
				} else {
					if err := c.compile(arg); err != nil {
						return err
					}
				}
			}
			// 2c. Emit OpCallMethod.
			// The opcode needs the method name index and number of arguments.
			methodName := fun.Sel.Name
			methodNameConstIndex := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, methodName))
			numArgs := len(node.Args)
			if _, err := c.scopes.Emit(bytecode.OpCallMethod, methodNameConstIndex, numArgs); err != nil {
				return err
			}
			// We're done for this case, no need to do anything else.
			return nil
		}

		// Path 3: Method call on a struct
		if receiverSymbol.IsStruct() {
			mangledName := tables2.GetMangledName(receiverSymbol.StructName(), fun.Sel.Name)
			methodSymbol, ok := c.scopes.SymbolResolve(mangledName)
			if !ok {
				return tables2.NewCompilerError(c.fileSet, node, "undefined method '%s' for type '%s'", fun.Sel.Name, receiverSymbol.StructName())
			}
			if err := c.scopes.EmitSymbolGet(methodSymbol); err != nil {
				return err
			}
			finalArgs = append([]ast.Expr{fun.X}, node.Args...)
		} else {
			return tables2.NewCompilerError(c.fileSet, node, "cannot call method on non-struct/non-interface type for '%s'", receiverSymbol.Name())
		}

	default:
		return tables2.NewCompilerError(c.fileSet, node, "unsupported function call type: %T", node.Fun)
	}

	// Step 3 & 4 for simple functions and struct methods
	for _, arg := range finalArgs {
		if tempSymbol, ok := tempSymbolMap[arg]; ok {
			if err := c.scopes.EmitSymbolGet(tempSymbol); err != nil {
				return err
			}
		} else {
			if err := c.compile(arg); err != nil {
				return err
			}
		}
	}
	if _, err := c.scopes.Emit(bytecode.OpCall, len(finalArgs), 0); err != nil {
		return err
	}
	return nil
}

// BlockStmt compiles each statement in the provided block and returns an error if any compilation step fails.
func (c *Functions) BlockStmt(node *ast.BlockStmt) error {
	for _, s := range node.List {
		if err := c.compile(s); err != nil {
			return err
		}
	}
	return nil
}

// ExprStmt compiles an expression statement and emits a pop operation to discard its result.
func (c *Functions) ExprStmt(node *ast.ExprStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// ReturnStmt compiles a return statement, handling cases for both void and value returns, and emits corresponding bytecode.
func (c *Functions) ReturnStmt(node *ast.ReturnStmt) error {
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

// FuncDecl processes the function declaration node and compiles its structure into the appropriate bytecode.
func (c *Functions) FuncDecl(_ *ast.FuncDecl) error {
	return nil
}

// FuncLit compiles an anonymous function literal.
// It creates a new scope, compiles the function body, and emits an OpClosure
// instruction to create the closure object at runtime.
func (c *Functions) FuncLit(node *ast.FuncLit) error {
	// 1. Enter a new scope for the anonymous function.
	if err := c.scopes.Enter(""); err != nil { // No struct or func name
		return err
	}
	// 2. Define symbols for the function parameters.
	if _, err := c.functionTable.SymbolsFromParameters(node.Type.Params); err != nil {
		return err
	}
	// 3. Compile the function body.
	if err := c.compile(node.Body); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	if scope.LastInstruction() == nil || scope.LastInstruction().Opcode() != bytecode.OpReturn {
		if _, err = c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
	}
	freeSymbols := c.scopes.SymbolFreeConvert()
	numFree := c.scopes.SymbolFreeCount()
	nParams := c.functionTable.CountParams(node.Type.Params)
	nLocals := c.scopes.SymbolCount()
	code, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	compiledFn := c.gk.NewFuncCompiled(objects.FrameStatic, "", code, nLocals, nParams, false, nil, freeSymbols)
	constIndex := c.constants.Add("", compiledFn)
	if _, err = c.scopes.Emit(bytecode.OpClosure, constIndex, numFree); err != nil {
		return err
	}
	return nil
}
