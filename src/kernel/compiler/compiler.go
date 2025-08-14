package compiler

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// mainFnName is the constant string representing the name of the main entry function in the program.
const (
	mainFnName = "main"
)

// maxScope defines the maximum depth allowed for compilation scopes to prevent excessive nesting during processing.
const (
	maxScope = 1024
)

// Compiler manages the organization and tracking of scopes and the main compiled function during a compilation process.
type Compiler struct {
	mainFn *objects.FunctionCompiled
	scopes *Scopes
}

// New initializes and returns a new instance of Compiler.
func New() *Compiler {
	return &Compiler{
		mainFn: nil,
		scopes: NewScopes(),
	}
}

// Bytecode generates and returns the compiled bytecode along with any errors encountered during compilation.
func (c *Compiler) Bytecode() (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecode()
	if c.mainFn == nil {
		return nil, errors.New("main function not found")
	}
	bc.SetMainFunction(c.mainFn)
	bc.SetConstants(c.scopes.ConstantsRetrieve())
	bc.SetReferences(c.scopes.ReferencesRetrieve())
	return bc, nil
}

// Compile processes the provided AST node and invokes the appropriate handler based on the node's type.
// It returns an error if the node type is unsupported or any processing issue occurs.
func (c *Compiler) Compile(in ast.Node) error {
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
	case *ast.ForStmt:
		err = c.doForStmt(node)
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

// doFile processes an AST file node by compiling its declarations and returns an error if the compilation fails.
func (c *Compiler) doFile(node *ast.File) error {
	for _, s := range node.Decls {
		if err := c.Compile(s); err != nil {
			return err
		}
	}
	return nil
}

// doDeclStmt processes a declaration statement within the AST and compiles it, returning an error if compilation fails.
func (c *Compiler) doDeclStmt(node *ast.DeclStmt) error {
	if err := c.Compile(node.Decl); err != nil {
		return err
	}
	return nil
}

// doBlockStmt processes a block statement by iterating through its statements and compiling each one.
// Returns an error if compilation of any statement fails.
func (c *Compiler) doBlockStmt(node *ast.BlockStmt) error {
	for _, s := range node.List {
		if err := c.Compile(s); err != nil {
			return err
		}
	}
	return nil
}

// doExprStmt compiles an expression statement and removes its value from the stack if it is unused.
func (c *Compiler) doExprStmt(node *ast.ExprStmt) error {
	if err := c.Compile(node.X); err != nil {
		return err
	}
	// Remove value from stack if unused
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
		if err := c.Compile(expr); err != nil {
			return err
		}
		ident := node.Lhs[i].(*ast.Ident)
		if node.Tok == token.DEFINE { // Handles 'x := 10'
			symbol := c.scopes.SymbolDefine(ident.Name)
			// AND USE THE NEW FUNCTION HERE TOO
			if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
				return err
			}
		} else { // Handles 'x = 20'
			symbol, ok := c.scopes.SymbolResolve(ident.Name)
			if !ok {
				return fmt.Errorf("undefined variable: %s", ident.Name)
			}
			// Assignment continues to use the old function
			if err := c.scopes.EmitSymbolSet(symbol); err != nil {
				return err
			}
		}
	}
	return nil
}

// doGenDecl processes a generic declaration node and compiles each specification contained within the declaration.
func (c *Compiler) doGenDecl(node *ast.GenDecl) error {
	for _, spec := range node.Specs {
		if err := c.Compile(spec); err != nil {
			return err
		}
	}
	return nil
}

// doValueSpec processes a ValueSpec node, compiles its values, and defines symbols in the current scope.
func (c *Compiler) doValueSpec(node *ast.ValueSpec) error {
	// Handles 'var x = 10'
	for i, name := range node.Names {
		if err := c.Compile(node.Values[i]); err != nil {
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
	if err := c.Compile(node.Fun); err != nil {
		return err
	}
	// Compile arguments
	for _, arg := range node.Args {
		if err := c.Compile(arg); err != nil {
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
	if err := c.Compile(node.Cond); err != nil {
		return err
	}
	// Emit conditional jump with temporary address
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	// Compile 'then' block
	if err = c.Compile(node.Body); err != nil {
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
		if err = c.Compile(node.Else); err != nil {
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

// doForStmt compiles an AST ForStmt node into bytecode, handling the loop's condition, body, jump logic, and stack cleanup.
func (c *Compiler) doForStmt(node *ast.ForStmt) error {
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	// Starting position for loop condition
	loopStartPos := scope.InstructionsLen()
	// Compile condition
	if err = c.Compile(node.Cond); err != nil {
		return err
	}
	// Emit conditional jump to exit loop
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}
	// Compile loop body
	if err = c.Compile(node.Body); err != nil {
		return err
	}
	// Emit unconditional jump to return to condition start
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}
	// Update (back-patching) OpJumpFalsy address to point to loop end
	scope, err = c.scopes.Current()
	if err != nil {
		return err
	}
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	// Remove condition value from stack after loop terminates
	if _, err = c.scopes.Emit(bytecode.OpPop); err != nil {
		return err
	}
	return nil
}

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
	if err := c.Compile(node.Body); err != nil {
		return err
	}
	// Implicit return if missing
	if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
		return err
	}
	numLocals := c.scopes.SymbolCount()
	instructions, err := c.scopes.Leave()
	if err != nil {
		return err
	}
	numParams := 0
	varArgs := false
	if paramL := node.Type.Params.List; paramL != nil {
		if numParams = len(paramL); numParams > 0 {
			if _, ok := paramL[numParams-1].Type.(*ast.Ellipsis); ok {
				varArgs = true
			}
		}
	}
	// Create compiled function object
	//TODO sourceMap
	compiledFn := objects.NewFunctionCompiled(node.Name.String(), instructions, numLocals, numParams, varArgs, nil, c.scopes.SymbolFreeConvert())
	fnIndex := c.scopes.ConstantsAdd(compiledFn)
	if _, err = c.scopes.Emit(bytecode.OpClosure, fnIndex, c.scopes.SymbolFreeCount()); err != nil {
		return err
	}
	// Define function in current scope
	symbol := c.scopes.SymbolDefine(node.Name.Name)
	if err = c.scopes.EmitSymbolSet(symbol); err != nil {
		return err
	}
	if node.Name.Name == mainFnName {
		c.mainFn = compiledFn
	}
	return nil
}

// doReturnStmt compiles a return statement, handling both void and value returns, and emits corresponding bytecode.
func (c *Compiler) doReturnStmt(node *ast.ReturnStmt) error {
	if len(node.Results) == 0 {
		// Return 'undefined'
		if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
			return err
		}
		return nil
	}
	if err := c.Compile(node.Results[0]); err != nil {
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
	if err := c.Compile(node.X); err != nil {
		return err
	}
	if err := c.Compile(node.Y); err != nil {
		return err
	}
	if err := c.scopes.EmitBinaryOp(node.Op); err != nil {
		return err
	}
	return nil
}

// doUnaryExpr compiles a unary expression by processing its operand and emitting the associated unary operation.
func (c *Compiler) doUnaryExpr(node *ast.UnaryExpr) error {
	if err := c.Compile(node.X); err != nil {
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
			if err := c.Compile(elt); err != nil {
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
			if err := c.Compile(kve.Key); err != nil {
				return err
			}
			if err := c.Compile(kve.Value); err != nil {
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
	if err := c.Compile(node.X); err != nil {
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
