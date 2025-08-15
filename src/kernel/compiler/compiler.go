package compiler

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// maxScope defines the maximum depth allowed for compilation scopes to prevent excessive nesting during processing.
const (
	maxScope = 1024
)

// Compiler manages the organization and tracking of scopes and the main compiled function during a compilation process.
type Compiler struct {
	scopes *Scopes
}

// New initializes and returns a new instance of Compiler.
func New() *Compiler {

	c := &Compiler{
		scopes: NewScopes(),
	}
	return c
}

// Compile compiles the provided source code and returns the compiled bytecode along with any errors encountered during compilation.
// It accepts a filename and source code as arguments and returns a tuple containing the compiled bytecode and any errors encountered.
// The source code can be a string or a []byte.
func (c *Compiler) Compile(filename string, source any) (*bytecode.Bytecode, error) {
	fSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fSet, filename, source, 0)
	if err != nil {
		return nil, err
	}
	if err = c.compile(astFile); err != nil {
		return nil, err
	}
	return c.bytecode()
}

// Bytecode generates and returns the compiled bytecode along with any errors encountered during compilation.
func (c *Compiler) bytecode() (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecode()
	bc.SetConstants(c.scopes.ConstantsRetrieve())
	bc.SetReferences(c.scopes.ReferencesRetrieve())
	return bc, nil
}

// compile processes an AST node and compiles it into bytecode.
// It recursively traverses the AST and calls the appropriate compile method for each node type.
// Returns an error if compilation fails.
func (c *Compiler) compile(in ast.Node) error {
	var err error = nil
	switch node := in.(type) {
	case *ast.File:
		err = c.doFile(node)
	case *ast.DeclStmt:
		err = c.doDeclStmt(node)
	case *ast.GenDecl: // For `var` and `const` which are handled by AssignStmt
		err = c.doGenDecl(node)
	case *ast.ValueSpec: // Handles 'var x = 10'
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

// doFile processes an AST file, separates declarations, and compiles both functions and non-function code.
// It predefines functions for stable indexing and resolves their bodies after global declarations are compiled.
// Returns an error if compilation fails at any stage.
func (c *Compiler) doFile(node *ast.File) error {
	var funcDecls []*ast.FuncDecl
	var otherDecls []ast.Decl
	// Step 1: Separate function declarations from all other declarations.
	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcDecls = append(funcDecls, fn)
		} else {
			otherDecls = append(otherDecls, decl)
		}
	}
	// Step 2: Pre-define functions.
	// This reserves a stable index in the constant table and symbol table
	// for each function before any other bytecode is generated.
	funcIndexes := make(map[string]int)
	for _, fn := range funcDecls {
		placeholder := objects.NewFunctionCompiled(fn.Name.Name, nil, 0, 0, false, nil, nil)
		fnIndex := c.scopes.ConstantsAdd(placeholder)
		c.scopes.SymbolDefine(fn.Name.Name)
		funcIndexes[fn.Name.Name] = fnIndex
	}
	// Step 3: Compile all non-function code (imports, global vars, etc.).
	// This can now add other constants without disturbing function indices.
	for _, decl := range otherDecls {
		if err := c.compile(decl); err != nil {
			return err
		}
	}
	// Step 4: Compile function bodies.
	// Now that indices are stable, we can compile the internal code.
	for _, fn := range funcDecls {
		if err := c.compileFuncBody(fn, funcIndexes[fn.Name.Name]); err != nil {
			return err
		}
	}
	return nil
}

// compileFuncBody compiles the body of a function, managing scoping, symbol resolution, and bytecode generation.
func (c *Compiler) compileFuncBody(node *ast.FuncDecl, constIndex int) error {
	symbol, _ := c.scopes.SymbolResolve(node.Name.Name)
	if err := c.scopes.Enter(); err != nil {
		return err
	}
	for _, p := range node.Type.Params.List {
		for _, name := range p.Names {
			c.scopes.SymbolDefine(name.Name)
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
	nLocals := c.scopes.SymbolCount()
	freeSymbols := c.scopes.SymbolFreeConvert()
	numFree := c.scopes.SymbolFreeCount()
	code, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	nParams := 0
	varArgs := false
	if paramL := node.Type.Params; paramL != nil && paramL.List != nil {
		if nParams = len(paramL.List); nParams > 0 {
			lastParam := paramL.List[nParams-1]
			if _, ok := lastParam.Type.(*ast.Ellipsis); ok {
				varArgs = true
			}
		}
	}
	compiledFn := objects.NewFunctionCompiled(node.Name.Name, code, nLocals, nParams, varArgs, nil, freeSymbols)
	if err = c.scopes.ConstantsSetIndex(constIndex, compiledFn); err != nil {
		return err
	}
	if _, err = c.scopes.Emit(bytecode.OpClosure, constIndex, numFree); err != nil {
		return err
	}
	if err = c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	return nil
}

// doDeclStmt processes a declaration statement within the AST and compiles it, returning an error if compilation fails.
func (c *Compiler) doDeclStmt(node *ast.DeclStmt) error {
	if err := c.compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// doBlockStmt processes a block statement by iterating through its statements and compiling each one.
// Returns an error if compilation of any statement fails.
func (c *Compiler) doBlockStmt(node *ast.BlockStmt) error {
	for _, s := range node.List {
		if err := c.compile(s); err != nil {
			return err
		}
	}
	return nil
}

// doExprStmt compiles an expression statement and removes its value from the stack if it is unused.
func (c *Compiler) doExprStmt(node *ast.ExprStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// doAssignStmt processes assignment statements, compiling their expressions and emitting the necessary symbol operations.
// Handles both variable definitions (`x := 10`) and regular assignments (`x = 20`).
// Returns an error if compilation or symbol resolution encounters a problem.
func (c *Compiler) doAssignStmt(node *ast.AssignStmt) error {
	for i, expr := range node.Rhs {
		if err := c.compile(expr); err != nil {
			return err
		}
		ident := node.Lhs[i].(*ast.Ident)
		if node.Tok == token.DEFINE { // Handles 'x := 10'
			symbol := c.scopes.SymbolDefine(ident.Name)
			if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
				return err
			}
		} else { // Handles 'x = 20'
			symbol, ok := c.scopes.SymbolResolve(ident.Name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", ident.Name)
			}
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

// doGenDecl processes a generic declaration node and compiles each specification contained within the declaration.
func (c *Compiler) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// doValueSpec processes a ValueSpec node, compiles its values, and defines symbols in the current scope.
func (c *Compiler) doValueSpec(node *ast.ValueSpec) error {
	// Handles 'var x = 10'
	for i, name := range node.Names {
		if err := c.compile(node.Values[i]); err != nil {
			return err
		}
		symbol := c.scopes.SymbolDefine(name.Name)
		if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
			return err
		}
	}
	return nil
}

// doCallExpr compiles a function call expression, including the function to call and its arguments.
// It emits the OpCall bytecode to execute the function call at runtime.
// Returns an error if compilation of the function or its arguments fails.
func (c *Compiler) doCallExpr(node *ast.CallExpr) error {
	// Compile function to call (e.g. identifier)
	if err := c.compile(node.Fun); err != nil {
		return err
	}
	// Compile arguments
	for _, arg := range node.Args {
		if err := c.compile(arg); err != nil {
			return err
		}
	}
	// 0 for non-spread call
	if _, err := c.scopes.Emit(bytecode.OpCall, len(node.Args), 0); err != nil {
		return err
	}
	return nil
}

// doIfStmt compiles an if statement, translating its condition, then block, and optional else block into bytecode.
// It manages jump instructions to handle conditional execution and ensures proper operand updates for jumps.
func (c *Compiler) doIfStmt(node *ast.IfStmt) error {
	// Compile condition
	if err := c.compile(node.Cond); err != nil {
		return err
	}
	// Emit conditional jump with temporary address
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	// Compile 'then' block
	if err = c.compile(node.Body); err != nil {
		return err
	}
	// If there's an 'else' block, emit jump to skip it
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
	// Update conditional jump address
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, scope.InstructionsLen()); err != nil {
		return err
	}
	// Compile 'else' block if it exists
	if node.Else != nil {
		if err = c.compile(node.Else); err != nil {
			return err
		}
		scope, err = c.scopes.Current()
		if err != nil {
			return err
		}
		// Update jump address to skip else
		if err = c.scopes.ChangeOperand(jumpToEndPos, scope.InstructionsLen()); err != nil {
			return err
		}
	}
	return nil
}

// doIncDecStmt processes increment and decrement statements in the AST and generates corresponding bytecode.
// It ensures the variable is defined, fetches its value, performs the operation, and saves the result.
// Returns an error if the variable is undefined or if an unsupported token is encountered.
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
	// 3. Aggiunge la costante '1' allo stack
	constIndex := c.scopes.ConstantsAdd(objects.NewInt(1))
	if _, err := c.scopes.Emit(bytecode.OpConstant, constIndex); err != nil {
		return err
	}
	if node.Tok == token.INC {
		if _, err := c.scopes.Emit(bytecode.OpBinaryOp, int(objects.OperatorAdd)); err != nil {
			return err
		}
	} else if node.Tok == token.DEC {
		if _, err := c.scopes.Emit(bytecode.OpBinaryOp, int(objects.OperatorSub)); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported IncDec token: %s", node.Tok)
	}
	if err := c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	// --- INIZIO DELLA CORREZIONE ---
	// L'operazione di incremento/decremento lascia il risultato sullo stack.
	// Dato che è un'istruzione, dobbiamo pulire questo valore.
	if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

// doForStmt compiles an AST ForStmt node into bytecode, handling the loop's
// initialization, condition, post-statement, body, and jump logic.
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
	// updates (back-patching) conditional jump address (OpJumpFalsy)
	// to point to loop end
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	return nil
}

// doRangeStmt processes a range statement, handling iteration and variable assignment in the compiled instructions.
func (c *Compiler) doRangeStmt(node *ast.RangeStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	iteratorSymbol := c.scopes.SymbolDefineUnique("__iterator")
	if _, err = c.scopes.Emit(bytecode.OpIteratorInit, iteratorSymbol.Index); err != nil {
		return err
	}
	var keySymbol, valueSymbol *Symbol
	if node.Key != nil {
		if ident, ok := node.Key.(*ast.Ident); ok && ident.Name != "_" {
			keySymbol = c.scopes.SymbolDefine(ident.Name)
		}
	}
	if node.Value != nil {
		if ident, ok := node.Value.(*ast.Ident); ok && ident.Name != "_" {
			valueSymbol = c.scopes.SymbolDefine(ident.Name)
		}
	}

	// 4. Inizio del ciclo
	loopStartPos := scope.InstructionsLen()

	// 5. Controlla se ci sono altri elementi, passando l'indice dell'iteratore.
	if _, err := c.scopes.Emit(bytecode.OpIteratorNext, iteratorSymbol.Index); err != nil {
		return err
	}
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}

	// 6. Assegna i valori e pulisce lo stack degli operandi
	if valueSymbol != nil {
		if _, err = c.scopes.Emit(bytecode.OpIteratorValue, iteratorSymbol.Index); err != nil {
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
		if _, err = c.scopes.Emit(bytecode.OpIteratorKey, iteratorSymbol.Index); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpIteratorValue, iteratorSymbol.Index); err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolSet(keySymbol); err != nil {
			return err
		}
		if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}

	// 7. Compila il corpo del ciclo
	if err = c.compile(node.Body); err != nil {
		return err
	}

	// 8. Salta all'inizio
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}

	// 9. Back-patching del jump di uscita
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) doFuncDecl(node *ast.FuncDecl) error {
	return nil
}

/*
// doFuncDecl compiles a function declaration into bytecode and manages the function's scope, parameters, and body.
func (c *Compiler) doFuncDecl(node *ast.FuncDecl) error {
	// Global function declaration
	if err := c.scopes.Enter(); err != nil {
		return err
	}
	// Receiver (methods) not supported in this version Parameters
	for _, p := range node.Type.Params.List {
		for _, name := range p.Names {
			c.scopes.SymbolDefine(name.Name)
		}
	}
	// Body
	if err := c.compile(node.Body); err != nil {
		return err
	}
	// Implicit return if missing
	if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
		return err
	}
	nLocals := c.scopes.SymbolCount()
	code, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	nParams := 0
	varArgs := false
	if paramL := node.Type.Params.List; paramL != nil {
		if nParams = len(paramL); nParams > 0 {
			if _, ok := paramL[nParams-1].Type.(*ast.Ellipsis); ok {
				varArgs = true
			}
		}
	}

	fName := node.Name.Name
	//TODO sourceMap
	compiledFn := objects.NewFunctionCompiled(fName, code, nLocals, nParams, varArgs, nil, c.scopes.SymbolFreeConvert())
	fnIndex := c.scopes.ConstantsAdd(compiledFn)
	if _, err = c.scopes.Emit(bytecode.OpClosure, fnIndex, c.scopes.SymbolFreeCount()); err != nil { //bytecode.OpClosure
		return err
	}
	// Define function in current scope
	symbol := c.scopes.SymbolDefine(fName)
	if err = c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	return nil
}

*/

// doReturnStmt compiles a return statement, handling both void and value returns, and emits corresponding bytecode.
func (c *Compiler) doReturnStmt(node *ast.ReturnStmt) error {
	if len(node.Results) == 0 {
		// Return 'undefined'
		if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
		return nil
	}
	if err := c.compile(node.Results[0]); err != nil {
		return err
	}
	// Return a value
	if _, err := c.scopes.Emit(bytecode.OpReturn, 1); err != nil {
		return err
	}
	return nil
}

// doBinaryExpr compiles a binary expression by processing its left and right operands and emitting the corresponding operation.
func (c *Compiler) doBinaryExpr(node *ast.BinaryExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if err := c.compile(node.Y); err != nil {
		return err
	}
	if err := c.scopes.EmitBinaryOp(node.Op); err != nil {
		return err
	}
	return nil
}

// doUnaryExpr compiles a unary expression by processing its operand and emitting the associated unary operation.
func (c *Compiler) doUnaryExpr(node *ast.UnaryExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	if err := c.scopes.EmitUnaryOp(node.Op); err != nil {
		return err
	}
	return nil
}

// doBasicLit processes a BasicLit AST node, emitting its corresponding literal and handling errors during compilation.
func (c *Compiler) doBasicLit(node *ast.BasicLit) error {
	if err := c.scopes.EmitLiteral(node); err != nil {
		return err
	}
	return nil
}

// doIdent resolves an identifier, checks its existence, and emits the corresponding symbol retrieval code.
func (c *Compiler) doIdent(node *ast.Ident) error {
	symbol, ok := c.scopes.SymbolResolve(node.Name)
	if !ok {
		return fmt.Errorf("undefined variable: %s", node.Name)
	}
	if err := c.scopes.EmitSymbolGet(symbol); err != nil {
		return err
	}
	return nil
}

// doCompositeLit processes composite literals such as arrays or maps and compiles their elements into bytecode.
// Returns an error if the literal type or elements cannot be compiled.
func (c *Compiler) doCompositeLit(node *ast.CompositeLit) error {
	// Check type to determine if array or map
	switch node.Type.(type) {
	case *ast.ArrayType:
		// Array literal (e.g. []int{1, 2, 3})
		for _, elt := range node.Elts {
			if err := c.compile(elt); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpArray, len(node.Elts)); err != nil {
			return err
		}
	case *ast.MapType:
		// Map literal (e.g. map[string]int{"a": 1})
		for _, elt := range node.Elts {
			kve := elt.(*ast.KeyValueExpr)
			if err := c.compile(kve.Key); err != nil {
				return err
			}
			if err := c.compile(kve.Value); err != nil {
				return err
			}
		}
		if _, err := c.scopes.Emit(bytecode.OpMap, len(node.Elts)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported composite literal type")
	}
	return nil
}

// doImportSpec processes an ImportSpec node, extracts the module name, and defines it in the current symbol scope.
func (c *Compiler) doImportSpec(node *ast.ImportSpec) error {
	moduleName := node.Path.Value
	c.scopes.SymbolDefine(strings.Trim(moduleName, "\"'"))
	return nil
}

// doSelectorExpr compiles a selector expression (e.g., 'fmt.Println') and emits the necessary bytecode instructions.
func (c *Compiler) doSelectorExpr(node *ast.SelectorExpr) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	moduleIdent, ok := node.X.(*ast.Ident)
	if !ok {
		// Per ora, gestiamo solo il caso semplice come 'fmt.print'
		// e non casi complessi come 'a[0].print()'.
		return fmt.Errorf("unsupported selector expression: %T", node.X)
	}
	mName := moduleIdent.Name
	sName := node.Sel.Name
	cacheKey := "selector:" + mName + "." + sName

	nameIndex, found := c.scopes.ReferencesGet(cacheKey)
	if !found {
		attrArray := objects.NewArray([]objects.IObject{objects.NewStringNoSize(mName), objects.NewStringNoSize(sName)})
		nameIndex = c.scopes.ReferencesAdd(cacheKey, attrArray)
	}
	if _, err := c.scopes.Emit(bytecode.OpReferences, nameIndex); err != nil {
		return err
	}
	return nil
}
