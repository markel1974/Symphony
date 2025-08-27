package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

const (
	UndefinedSymbol = "_"
)

// Functions is a collection that manages a list of function descriptions.
type Functions struct {
	gk            objects.IGateKeeper
	constants     *Constants
	scopes        *Scopes
	imports       *Imports
	declarations  *Declarations
	functionTable *FunctionTable
	structTable   *StructTable
	fileSet       *token.FileSet
	compile       func(node ast.Node) error
}

// NewFunctions initializes and returns a new Functions instance.
func NewFunctions(gk objects.IGateKeeper, constants *Constants, scopes *Scopes, imports *Imports, declarations *Declarations, structTable *StructTable) *Functions {
	return &Functions{
		gk:            gk,
		constants:     constants,
		scopes:        scopes,
		imports:       imports,
		declarations:  declarations,
		structTable:   structTable,
		functionTable: NewFunctionTable(gk, scopes, structTable),
		compile:       nil,
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
func (c *Functions) funcBodyPrepare(fd *FunctionDescription) error {
	node := fd.FuncDecl
	var err error
	if fd.ReturnValues, err = GetReceivers(node.Type.Results); err != nil {
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
		recvTypeIdent := GetIdent(node.Recv.List[0])
		if recvTypeIdent == nil {
			return NewCompilerError(c.fileSet, node, "unsupported method receiver type")
		}
		if !c.structTable.Has(recvTypeIdent.Name) {
			return NewCompilerError(c.fileSet, node, "undefined type '%s' for method receiver", recvTypeIdent.Name)
		}
		fd.Name = GetMangledName(recvTypeIdent.Name, node.Name.Name)
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
func (c *Functions) funcBodyCompile(fd *FunctionDescription) error {
	node := fd.FuncDecl
	if err := c.scopes.Enter(fd.Name); err != nil {
		return err
	}
	for _, recv := range fd.StructReceivers {
		receiverSymbol, err := c.scopes.SymbolDefine(recv)
		if err != nil {
			return err
		}
		if len(fd.StructName) > 0 {
			if err = c.structTable.AssignSymbol(receiverSymbol, fd.StructName, fd.ReturnValues); err != nil {
				return err
			}
		}
	}
	for _, p := range fd.InputParams {
		if _, err := c.scopes.SymbolDefine(p); err != nil {
			return err
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

	nParams := c.functionTable.CountParams(node.Type.Params)

	if node.Recv != nil && len(node.Recv.List) > 0 {
		nParams++
	}

	fnSymbol, ok := c.scopes.SymbolRebuildScope(fd.Name, UnknownScope)
	if !ok {
		return NewCompilerError(c.fileSet, node, "undefined function: %s", fd.Name)
	}
	compiledFn := c.gk.NewFuncCompiled(objects.FrameStatic, fd.Name, code, nLocals, nParams, false, nil, freeSymbols)
	fnSymbol.SetObject(compiledFn)
	fnSymbol.SetTypes(fd.ReturnValues)

	if node.Recv == nil && c.scopes.scopeIndex > 0 {
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
	// Step 1: Pre-evaluate nested function calls and store their results in temporary variables.
	// This ensures that arguments are evaluated before the main function is pushed onto the stack.
	tempSymbolMap := make(map[ast.Expr]*Symbol)
	for _, arg := range node.Args {
		if call, isCall := arg.(*ast.CallExpr); isCall {
			if err := c.compile(call); err != nil {
				return err
			}
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_call")
			if err != nil {
				return err
			}
			tempSymbol.SetScope(LocalScope)
			tempSymbolMap[arg] = tempSymbol
			if err = c.scopes.EmitSymbolSet(tempSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
				return err
			}
		}
	}

	// Step 2: Resolve and push the function/method object onto the stack.
	var finalArgs []ast.Expr

	switch fun := node.Fun.(type) {
	case *ast.Ident: // Chiamata a funzione semplice (es. myFunction())
		symbol, ok := c.scopes.SymbolResolve(fun.Name)
		if !ok {
			// Potrebbe essere una funzione da un pacchetto importato
			if c.imports.Emit(fun.Name, "") {
				finalArgs = node.Args
				break // Fatto, l'import ha già emesso il suo bytecode
			}
			return NewCompilerError(c.fileSet, node, "undefined function: %s", fun.Name)
		}
		// Emette l'opcode corretto (Global, Local, o Free) per caricare la funzione
		if err := c.scopes.EmitSymbolGet(symbol); err != nil {
			return err
		}
		finalArgs = node.Args

	case *ast.SelectorExpr: // Chiamata a metodo (myVar.Method()) o funzione di pacchetto (fmt.Println())
		receiverIdent, ok := fun.X.(*ast.Ident)
		if !ok {
			return NewCompilerError(c.fileSet, node, "unsupported receiver for selector expression: %T", fun.X)
		}

		// Prova a risolverlo come funzione di pacchetto (es. fmt.Println)
		if c.imports.Emit(receiverIdent.Name, fun.Sel.Name) {
			finalArgs = node.Args
			break // Fatto
		}

		// Altrimenti, trattalo come una chiamata a metodo di uno struct
		receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
		if !ok {
			return NewCompilerError(c.fileSet, node, "undefined variable: %s", receiverIdent.Name)
		}

		if !receiverSymbol.IsStruct() {
			return NewCompilerError(c.fileSet, node, "cannot call method on non-struct type '%s'", receiverSymbol.Name())
		}

		// Il nome "mangled" del metodo è 'StructName.MethodName'
		mangledName := GetMangledName(receiverSymbol.StructName(), fun.Sel.Name)
		methodSymbol, ok := c.scopes.SymbolResolve(mangledName)
		if !ok {
			return NewCompilerError(c.fileSet, node, "undefined method '%s' for type '%s'", fun.Sel.Name, receiverSymbol.StructName())
		}

		// Emette il codice per caricare il metodo
		if err := c.scopes.EmitSymbolGet(methodSymbol); err != nil {
			return err
		}

		// Il primo argomento di un metodo è sempre il suo ricevitore (l'istanza dello struct)
		finalArgs = append([]ast.Expr{fun.X}, node.Args...)

	default:
		return NewCompilerError(c.fileSet, node, "unsupported function call type: %T", node.Fun)
	}

	// Step 3: Push all final arguments for the main call onto the stack.
	for _, arg := range finalArgs {
		if tempSymbol, ok := tempSymbolMap[arg]; ok {
			// Argomento pre-calcolato: carica il risultato dalla variabile temporanea
			if err := c.scopes.EmitSymbolGet(tempSymbol); err != nil {
				return err
			}
		} else {
			// Argomento normale: compilalo per mettere il suo valore sullo stack
			if err := c.compile(arg); err != nil {
				return err
			}
		}
	}

	// Step 4: Emit the final OpCall instruction.
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

// IfStmt compiles an if statement, handling both 'then' and optional 'else' blocks with associated bytecode instructions.
func (c *Functions) IfStmt(node *ast.IfStmt) error {
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

// IncDecStmt handles increment and decrement statements for identifiers, updating the corresponding variables and cleaning the stack.
func (c *Functions) IncDecStmt(node *ast.IncDecStmt) error {
	ident, ok := node.X.(*ast.Ident)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unsupported IncDec statement for type %T", node.X)
	}
	symbol, ok := c.scopes.SymbolResolve(ident.Name)
	if !ok {
		return NewCompilerError(c.fileSet, node, "undefined variable: %s", ident.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	// adds constant '1' to the stack
	constIndex := c.constants.Add("", c.gk.NewInt(objects.FrameStatic, 1))
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
		return NewCompilerError(c.fileSet, node, "unsupported IncDec token: %s", node.Tok)
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

// ForStmt compiles a for loop statement, including initialization, condition, post-iteration, and body execution.
func (c *Functions) ForStmt(node *ast.ForStmt) error {
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

// RangeStmt compiles a RangeStmt node into bytecode, handling iterator initialization, key/value assignment, and looping logic.
func (c *Functions) RangeStmt(node *ast.RangeStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	iteratorSymbol, err := c.scopes.SymbolDefineUnique("__iterator")
	if err != nil {
		return err
	}
	if _, err = c.scopes.Emit(bytecode.OpIteratorInit, iteratorSymbol.Index()); err != nil {
		return err
	}
	var returnTypeName string
	switch expr := node.X.(type) {
	case *ast.Ident:
		returnTypeName, _ = c.structTable.TypeNameFromSymbol(expr.Name)
	case *ast.CallExpr:
		if ident, ok := expr.Fun.(*ast.Ident); ok {
			returnTypeName, _ = c.structTable.TypeNameFromSymbol(ident.Name)
		}
	case *ast.SelectorExpr:
		// Caso: for _, v := range myVar.Items
		if receiverIdent, ok := expr.X.(*ast.Ident); ok {
			//1. Risolviamo il simbolo del ricevitore (myVar)
			returnTypeName, _ = c.structTable.TypeNameFromSymbolField(receiverIdent.Name, expr.Sel.Name)
		}
	default:
		return NewCompilerError(c.fileSet, node, "unsupported range expression: %T", node.X)
	}
	keySymbol, err := c.functionTable.RangeKey(node)
	if err != nil {
		return err
	}
	valueSymbol, err := c.functionTable.RangeValue(node, returnTypeName)
	if err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	loopStartPos := scope.InstructionsLen()
	if _, err = c.scopes.Emit(bytecode.OpIteratorNext, iteratorSymbol.Index()); err != nil {
		return err
	}
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
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
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
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

// BinaryExpr processes a binary expression node, compiling both operands and emitting the corresponding binary operation.
func (c *Functions) BinaryExpr(node *ast.BinaryExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if err := c.compile(node.Y); err != nil {
		return err
	}
	z, ok := BinaryAdapterFor(node.Op)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unhandled binary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// UnaryExpr compiles a unary expression by evaluating the operand and applying the specified unary operator.
// It handles special cases for the address-of operator '&', ensuring correct pointer behavior based on operand type.
// Emits appropriate bytecode instructions for each unary operation or returns an error on unsupported cases.
func (c *Functions) UnaryExpr(node *ast.UnaryExpr) error {
	if node.Op == token.AND {
		switch operand := node.X.(type) {
		case *ast.Ident:
			// literal (es. '&h').
			symbol, ok := c.scopes.SymbolResolve(operand.Name)
			if !ok {
				return NewCompilerError(c.fileSet, node, "undefined variable: %s", operand.Name)
			}
			switch symbol.Scope() {
			case LocalScope:
				if _, err := c.scopes.Emit(bytecode.OpLocalPtrGet, symbol.Index()); err != nil {
					return err
				}
			case FreeScope:
				if _, err := c.scopes.Emit(bytecode.OpFreePtrGet, symbol.Index()); err != nil {
					return err
				}
			default:
				return NewCompilerError(c.fileSet, node, "cannot take the address of a global variable")
			}
		case *ast.CompositeLit:
			// literal (es. '&Home{...}').
			if err := c.compile(operand); err != nil {
				return err
			}
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_struct")
			if err != nil {
				return err
			}
			tempSymbol.SetScope(LocalScope)
			if err = c.scopes.EmitSymbolDefine(tempSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.Emit(bytecode.OpLocalPtrGet, tempSymbol.Index()); err != nil {
				return err
			}
		default:
			return NewCompilerError(c.fileSet, node, "cannot take the address of %T", node.X)
		}
		return nil
	}
	// existing logic for other unary operators (e.g. '!', '-', '^')
	if err := c.compile(node.X); err != nil {
		return err
	}
	z, ok := UnaryAdapterFor(node.Op)
	if !ok {
		return NewCompilerError(c.fileSet, node, "unhandled unary op: %s", node.Op)
	}
	if _, err := c.scopes.Emit(z.op, z.arguments...); err != nil {
		return err
	}
	return nil
}

// SelectorExpr processes a selector expression, resolving fields, methods, or package attributes.
// It distinguishes between struct field accesses and package-level selectors.
// Emits appropriate bytecode instructions for each case or returns an error if unsupported.
func (c *Functions) SelectorExpr(node *ast.SelectorExpr) error {
	// analyze the left-hand side of the dot to determine if it's a variable or a package.
	receiverIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// currently not handling complex cases like a[0].field
		return NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported receiver for selector expression: %T", node.X)
	}
	if c.imports.Emit(receiverIdent.Name, node.Sel.Name) {
		return nil
	}
	receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
	if !ok {
		return NewCompilerError(c.fileSet, node, "[SelectorExpr] undefined variable: %s", receiverIdent.Name)
	}
	if receiverSymbol.IsStruct() { // struct
		if err := c.compile(node.X); err != nil {
			return err
		}
		fieldName := node.Sel.Name
		keyConst := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, fieldName))
		if _, err := c.scopes.Emit(bytecode.OpConstant, keyConst); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpIndex); err != nil {
			return err
		}
		return nil
	}
	return NewCompilerError(c.fileSet, node, "[SelectorExpr] unsupported selector expression for symbol %s", receiverSymbol.Name())
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
	if scope.LastInstruction() == nil || scope.LastInstruction().opcode != bytecode.OpReturn {
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
