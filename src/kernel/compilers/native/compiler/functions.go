package compiler

import (
	"fmt"
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
		functionTable: NewFunctionTable(),
		compile:       nil,
	}
}

// Setup initializes the `Functions` instance with a compile function used for processing AST nodes.
func (c *Functions) Setup(compile func(node ast.Node) error) error {
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
			return fmt.Errorf("unsupported method receiver type")
		}
		if !c.structTable.Has(recvTypeIdent.Name) {
			return fmt.Errorf("undefined type '%s' for method receiver", recvTypeIdent.Name)
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
	nParams := 0
	if paramL := node.Type.Params; paramL != nil && paramL.List != nil {
		for _, field := range paramL.List {
			nParams += len(field.Names)
		}
	}
	if node.Recv != nil && len(node.Recv.List) > 0 {
		nParams++
	}

	fnSymbol, ok := c.scopes.SymbolRebuildScope(fd.Name, UnknownScope)
	if !ok {
		return fmt.Errorf("undefined function: %s", fd.Name)
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
	// Step 1: Resolve the function and its arguments for analysis. This part doesn't emit bytecode yet.
	identName := ""
	receiverIdentName := ""
	selName := ""
	commonName := ""

	selExpr, isSelector := node.Fun.(*ast.SelectorExpr)
	if isSelector {
		receiverIdent, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported receiver for selector expression: %T", selExpr.X)
		}
		commonName = receiverIdent.Name
		receiverIdentName = receiverIdent.Name
		selName = selExpr.Sel.Name
	} else {
		ident, ok := node.Fun.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported function call: %T", node.Fun)
		}
		commonName = ident.Name
		identName = ident.Name
	}

	fnOpType := bytecode.OpGlobalGet
	var fnArgs []ast.Expr
	fnName, fnIndex, ok := c.imports.Attach(commonName, selName)
	if ok {
		fnOpType = bytecode.OpReferences
		fnArgs = node.Args
	} else {
		if selExpr != nil {
			receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdentName)
			if !ok {
				return fmt.Errorf("undefined variable: %s", receiverIdentName)
			}
			if len(receiverSymbol.Types()) > 0 { // struct method
				structTypeName := receiverSymbol.Types()[0]
				if !c.structTable.Has(structTypeName) {
					return fmt.Errorf("undefined type: %s", structTypeName)
				}
				fnName = GetMangledName(structTypeName, selName)
				fnSymbol, ok := c.scopes.SymbolResolve(fnName)
				if !ok {
					return fmt.Errorf("undefined method '%s' for type '%s' [%s]", selName, structTypeName, fnName)
				}
				fnIndex = fnSymbol.Index()
				fnArgs = append(fnArgs, selExpr.X) // The receiver is the first argument
				fnArgs = append(fnArgs, node.Args...)
			} else {
				return fmt.Errorf("cannot call method on untyped variable or undefined package '%s'", receiverSymbol.Name())
			}
		} else {
			symbol, ok := c.scopes.SymbolResolve(identName)
			if !ok {
				return fmt.Errorf("undefined function: %s", identName)
			}
			fnName = identName
			fnIndex = symbol.Index()
			fnArgs = node.Args
		}
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
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_call")
			if err != nil {
				return err
			}
			tempSymbol.SetScope(LocalScope)
			tempSymbolMap[arg] = tempSymbol // Store the symbol for the second pass
			// 1.3. Emit code to store the result into the temp variable.
			// This generates OpLocalSet, correctly storing the result.
			if err = c.scopes.EmitSymbolSet(tempSymbol); err != nil {
				return err
			}
			// 1.4. Pop the result from the stack to keep it clean for the next operations.
			// This is crucial as OpLocalSet only peeks at the stack value.
			if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
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
	scope, err := c.scopes.Current()
	if err != nil {
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
		if symbol, ok := c.scopes.SymbolResolve(expr.Name); ok {
			if len(symbol.Types()) > 0 {
				returnTypeName = symbol.Types()[0]
			}
		}
	case *ast.CallExpr:
		if ident, ok := expr.Fun.(*ast.Ident); ok {
			// Simbolo della funzione, per inferire il tipo di ritorno
			if symbol, ok := c.scopes.SymbolResolve(ident.Name); ok {
				if len(symbol.Types()) > 0 {
					returnTypeName = symbol.Types()[0]
				}
			}
		}
	case *ast.SelectorExpr:
		// Caso: for _, v := range myVar.Items
		if receiverIdent, ok := expr.X.(*ast.Ident); ok {
			// 1. Risolviamo il simbolo del ricevitore (myVar)
			if receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name); ok {
				if returnTypeName, ok = c.structTable.GetTypeNameFromFields(receiverSymbol.StructName(), expr.Sel.Name); !ok {
					return fmt.Errorf("undefined field: %s.%s", receiverSymbol.StructName(), expr.Sel.Name)
				}
			}
		}
	default:
		return fmt.Errorf("unsupported range expression: %T", node.X)
	}

	var keySymbol, valueSymbol *Symbol
	if node.Key != nil {
		if ident, ok := node.Key.(*ast.Ident); ok && ident.Name != UndefinedSymbol {
			keySymbol, err = c.scopes.SymbolDefine(ident.Name)
			if err != nil {
				return err
			}
		}
	}

	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != UndefinedSymbol {
			valueSymbol, err = c.scopes.SymbolDefine(ident.Name)
			if err != nil {
				return err
			}
			if len(returnTypeName) > 0 {
				if err = c.structTable.AssignSymbol(valueSymbol, returnTypeName, []string{returnTypeName}); err != nil {
					return err
				}
			}
		}
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
		return fmt.Errorf("unhandled binary op: %s", node.Op)
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
				return fmt.Errorf("undefined variable: %s", operand.Name)
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
				return fmt.Errorf("cannot take the address of a global variable")
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

// SelectorExpr processes a selector expression, resolving fields, methods, or package attributes.
// It distinguishes between struct field accesses and package-level selectors.
// Emits appropriate bytecode instructions for each case or returns an error if unsupported.
func (c *Functions) SelectorExpr(node *ast.SelectorExpr) error {
	// analyze the left-hand side of the dot to determine if it's a variable or a package.
	receiverIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// currently not handling complex cases like a[0].field
		return fmt.Errorf("[SelectorExpr] unsupported receiver for selector expression: %T", node.X)
	}
	if c.imports.Emit(receiverIdent.Name, node.Sel.Name) {
		return nil
	}
	receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
	if !ok {
		return fmt.Errorf("[SelectorExpr] undefined variable: %s", receiverIdent.Name)
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
	return fmt.Errorf("[SelectorExpr] unsupported selector expression for symbol %s", receiverSymbol.Name())
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
	var paramNames []string
	if node.Type.Params != nil {
		for _, p := range node.Type.Params.List {
			var typeName string
			switch t := p.Type.(type) {
			case *ast.Ident:
				typeName = t.Name
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					typeName = ident.Name
				}
			}
			structName := ""
			isStruct := c.structTable.Has(typeName)
			if isStruct {
				structName = typeName
			}
			for _, name := range p.Names {
				paramNames = append(paramNames, name.Name)
				zSymbol, err := c.scopes.SymbolDefine(name.Name)
				if err != nil {
					return err
				}
				zSymbol.SetScope(LocalScope)
				if len(structName) > 0 {
					if err = c.structTable.AssignSymbol(zSymbol, structName, nil); err != nil {
						return err
					}
				}
			}
		}
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
	nLocals := c.scopes.SymbolCount()
	code, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	compiledFn := c.gk.NewFuncCompiled(objects.FrameStatic, "", code, nLocals, len(paramNames), false, nil, freeSymbols)
	constIndex := c.constants.Add("", compiledFn)
	if _, err = c.scopes.Emit(bytecode.OpClosure, constIndex, numFree); err != nil {
		return err
	}
	return nil
}
