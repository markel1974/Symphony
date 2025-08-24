package compiler

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// FunctionDescription represents the metadata of a function including its name, associated struct, parameters, and receiver info.
type FunctionDescription struct {
	Name       string
	Struct     string
	Types      []string
	Params     []string
	Recv       []string
	FuncDecl   *ast.FuncDecl
	StructType bool
}

// NewFunctionDescription creates a new instance of FunctionDescription with the provided function declaration.
func NewFunctionDescription(funcDecl *ast.FuncDecl) *FunctionDescription {
	return &FunctionDescription{
		FuncDecl: funcDecl,
	}
}

// Functions is a collection that manages a list of function descriptions.
type Functions struct {
	gk           objects.IGateKeeper
	constants    *Constants
	scopes       *Scopes
	imports      *Imports
	declarations *Declarations
	container    []*FunctionDescription
}

// NewFunctions initializes and returns a new Functions instance.
func NewFunctions(gk objects.IGateKeeper, constants *Constants, scopes *Scopes, imports *Imports, declarations *Declarations) *Functions {
	return &Functions{
		gk:           gk,
		constants:    constants,
		scopes:       scopes,
		imports:      imports,
		declarations: declarations,
	}
}

// Declare adds a function description derived from the provided function declaration to the Functions container.
func (c *Functions) Declare(funcDecl *ast.FuncDecl) {
	fd := NewFunctionDescription(funcDecl)
	c.container = append(c.container, fd)
}

func (c *Functions) Prepare() error {
	for _, fd := range c.container {
		if err := c.funcBodyPrepare(fd); err != nil {
			return err
		}
	}
	return nil
}

func (c *Functions) Compile() error {
	for _, fd := range c.container {
		if err := c.funcBodyCompile(fd); err != nil {
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
	if fd.Types, err = GetReceivers(node.Type.Results); err != nil {
		return err
	}
	for _, p := range node.Type.Params.List {
		for _, name := range p.Names {
			fd.Params = append(fd.Params, name.Name)
		}
	}
	if node.Recv != nil && len(node.Recv.List) > 0 {
		for _, p := range node.Recv.List {
			for _, name := range p.Names {
				fd.Recv = append(fd.Recv, name.Name)
			}
		}
		recvTypeIdent := GetIdent(node.Recv.List[0])
		if recvTypeIdent == nil {
			return fmt.Errorf("unsupported method receiver type")
		}
		baseSymbol, ok := c.scopes.SymbolResolve(recvTypeIdent.Name)
		if !ok {
			return fmt.Errorf("undefined type '%s' for method receiver", recvTypeIdent.Name)
		}
		if !baseSymbol.IsStruct() {
			return fmt.Errorf("unknown type '%s' for method receiver", recvTypeIdent.Name)
		}
		fd.Name = GetMangledName(recvTypeIdent.Name, node.Name.Name)
		fd.Struct = recvTypeIdent.Name
		fd.StructType = true
	} else {
		fd.Name = node.Name.Name
		fd.Struct = ""
		fd.StructType = false
	}
	//function symbol placeholder (this is not the real function, it's just a placeholder to be able to compile the body)
	if _, err = c.scopes.SymbolDefine(fd.Name, UnknownScope, false); err != nil {
		return err
	}
	return nil
}

// funcBodyCompile compiles the body of a function declaration and generates the necessary bytecode instructions.
func (c *Functions) funcBodyCompile(fd *FunctionDescription) error {
	node := fd.FuncDecl
	if err := c.scopes.Enter(fd.Struct, fd.Name); err != nil {
		return err
	}
	for _, p := range fd.Recv {
		if _, err := c.scopes.SymbolDefine(p, UnknownScope, fd.StructType); err != nil {
			return err
		}
	}
	for _, p := range fd.Params {
		if _, err := c.scopes.SymbolDefine(p, UnknownScope, false); err != nil {
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
	fnSymbol, err := c.scopes.SymbolReset(fd.Name, UnknownScope, false)
	if err != nil {
		return err
	}
	compiledFn := c.gk.NewFuncCompiled(objects.FrameStatic, fd.Name, code, nLocals, nParams, false, nil, freeSymbols)
	fnSymbol.SetObject(compiledFn)
	fnSymbol.SetTypes(fd.Types)

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

// compile traverses the provided AST node and compiles it into bytecode, handling various node types in a switch block.
func (c *Functions) compile(in ast.Node) error {
	var err error = nil

	switch node := in.(type) {
	case *ast.GenDecl:
		err = c.declarations.GenDecl(node)
	case *ast.Ident:
		err = c.declarations.Ident(node)
	case *ast.AssignStmt:
		err = c.declarations.AssignStmt(node)
	case *ast.BlockStmt:
		err = c.doBlockStmt(node)
	case *ast.ExprStmt:
		err = c.doExprStmt(node)
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
	case *ast.FuncDecl:
		err = c.doFuncDecl(node)
	case *ast.CallExpr:
		err = c.doCallExpr(node)
	case *ast.ReturnStmt:
		err = c.doReturnStmt(node)
	case *ast.SelectorExpr:
		err = c.doSelectorExpr(node)
	default:
		err = fmt.Errorf("[functions] unsupported expression type: %T", node)
	}
	return err
}

// doCallExpr processes a single CallExpr node and generates bytecode for function or method calls.
// It resolves the function or method being called, handles arguments, and emits corresponding bytecode.
// It also manages nested function calls by pre-evaluating them and storing results in temporary variables.
// Returns an error if the call expression contains invalid or unresolved references.
func (c *Functions) doCallExpr(node *ast.CallExpr) error {
	// Step 1: Resolve the function and its arguments for analysis. This part doesn't emit bytecode yet.
	fnOpType := bytecode.OpConstant
	fnIndex := -1
	fnName := ""
	var fnArgs []ast.Expr

	if selExpr, isSelector := node.Fun.(*ast.SelectorExpr); isSelector {
		// Import or Method call
		receiverIdent, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return fmt.Errorf("unsupported receiver for selector expression: %T", selExpr.X)
		}
		if c.imports.Contains(receiverIdent.Name) {
			var err error
			fnOpType = bytecode.OpReferences
			fnArgs = node.Args
			fnName, fnIndex, err = c.imports.Create(receiverIdent.Name, selExpr.Sel.Name)
			if err != nil {
				return err
			}
		} else {
			receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
			}
			if len(receiverSymbol.Types()) > 0 { // struct method
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
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_call", LocalScope, false)
			if err != nil {
				return err
			}
			tempSymbolMap[arg] = tempSymbol // Store the symbol for the second pass
			// 1.3. Emit code to store the result into the temp variable.
			// This generates OpSetLocal, correctly storing the result.
			if err = c.scopes.EmitSymbolSet(tempSymbol); err != nil {
				return err
			}
			// 1.4. Pop the result from the stack to keep it clean for the next operations.
			// This is crucial as OpSetLocal only peeks at the stack value.
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

// doBlockStmt compiles each statement in the provided block and returns an error if any compilation step fails.
func (c *Functions) doBlockStmt(node *ast.BlockStmt) error {
	for _, s := range node.List {
		if err := c.compile(s); err != nil {
			return err
		}
	}
	return nil
}

// doExprStmt compiles an expression statement and emits a pop operation to discard its result.
func (c *Functions) doExprStmt(node *ast.ExprStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// doIfStmt compiles an if statement, handling both 'then' and optional 'else' blocks with associated bytecode instructions.
func (c *Functions) doIfStmt(node *ast.IfStmt) error {
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
func (c *Functions) doIncDecStmt(node *ast.IncDecStmt) error {
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

// doForStmt compiles a for loop statement, including initialization, condition, post-iteration, and body execution.
func (c *Functions) doForStmt(node *ast.ForStmt) error {
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
func (c *Functions) doRangeStmt(node *ast.RangeStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	iteratorSymbol, err := c.scopes.SymbolDefineUnique("__iterator", UnknownScope, false)
	if err != nil {
		return err
	}
	if _, err = c.scopes.Emit(bytecode.OpIteratorInit, iteratorSymbol.Index()); err != nil {
		return err
	}
	var keySymbol, valueSymbol *Symbol
	if node.Key != nil {
		if ident, ok := node.Key.(*ast.Ident); ok && ident.Name != "_" {
			keySymbol, err = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
			if err != nil {
				return err
			}
		}
	}
	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != "_" {
			valueSymbol, err = c.scopes.SymbolDefine(ident.Name, UnknownScope, false)
			if err != nil {
				return err
			}
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
func (c *Functions) doReturnStmt(node *ast.ReturnStmt) error {
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
func (c *Functions) doBinaryExpr(node *ast.BinaryExpr) error {
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
func (c *Functions) doUnaryExpr(node *ast.UnaryExpr) error {
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
			tempSymbol, err := c.scopes.SymbolDefineUnique("__temp_struct", LocalScope, false)
			if err != nil {
				return err
			}
			if err = c.scopes.EmitSymbolDefine(tempSymbol); err != nil {
				return err
			}
			if _, err = c.scopes.Emit(bytecode.OpGetLocalPtr, tempSymbol.Index()); err != nil {
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

// doSelectorExpr processes a selector expression, resolving fields, methods, or package attributes.
// It distinguishes between struct field accesses and package-level selectors.
// Emits appropriate bytecode instructions for each case or returns an error if unsupported.
func (c *Functions) doSelectorExpr(node *ast.SelectorExpr) error {
	// analyze the left-hand side of the dot to determine if it's a variable or a package.
	receiverIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// currently not handling complex cases like a[0].field
		return fmt.Errorf("unsupported receiver for selector expression: %T", node.X)
	}
	if c.imports.Contains(receiverIdent.Name) {
		if _, _, err := c.imports.Create(receiverIdent.Name, node.Sel.Name); err != nil {
			return err
		}
		return nil
	}
	receiverSymbol, ok := c.scopes.SymbolResolve(receiverIdent.Name)
	if !ok {
		return fmt.Errorf("undefined variable: %s", receiverIdent.Name)
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
	return fmt.Errorf("unsupported selector expression: %T", node.X)
}

// doFuncDecl processes the function declaration node and compiles its structure into the appropriate bytecode.
func (c *Functions) doFuncDecl(_ *ast.FuncDecl) error {
	return nil
}
