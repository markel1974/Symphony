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

const (
	mainFnName = "main"
)

// maxScope defines the maximum allowable depth for scopes in the compiler to prevent excessive nesting or stack overflow.
const (
	maxScope = 1024
)

// Compiler manages the compilation process, including constant storage, scopes, and symbol resolution during program compilation.
type Compiler struct {
	mainFn *objects.FunctionCompiled
	scopes *Scopes
}

// New creates and initializes a new Compiler instance with a fresh symbol table and main compilation scope.
func New() *Compiler {
	return &Compiler{
		mainFn: nil,
		scopes: NewScopes(),
	}
}

// Compile traverses the provided AST node and compiles it into bytecode for execution by the virtual machine.
func (c *Compiler) Compile(in ast.Node) error {
	switch node := in.(type) {
	case *ast.File:
		for _, s := range node.Decls {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.DeclStmt:
		if err := c.Compile(node.Decl); err != nil {
			return err
		}
	case *ast.GenDecl: // For `var` and `const` which are handled by AssignStmt
		for _, spec := range node.Specs {
			if err := c.Compile(spec); err != nil {
				return err
			}
		}
	case *ast.ValueSpec: // Handles 'var x = 10'
		for i, name := range node.Names {
			if err := c.Compile(node.Values[i]); err != nil {
				return err
			}
			symbol := c.scopes.SymbolDefine(name.Name)
			if err := c.scopes.EmitSymbolDefine(symbol); err != nil {
				return err
			}
		}

	// --- Statements ---
	case *ast.BlockStmt:
		for _, s := range node.List {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		// Remove value from stack if unused
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	case *ast.AssignStmt:
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
	case *ast.IfStmt:
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
	case *ast.ForStmt:
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
	case *ast.BinaryExpr:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		if err := c.Compile(node.Y); err != nil {
			return err
		}
		if err := c.scopes.EmitBinaryOp(node.Op); err != nil {
			return err
		}
	case *ast.UnaryExpr:
		if err := c.Compile(node.X); err != nil {
			return err
		}
		if err := c.scopes.EmitUnaryOp(node.Op); err != nil {
			return err
		}
	case *ast.BasicLit:
		if err := c.scopes.EmitLiteral(node); err != nil {
			return err
		}
	case *ast.Ident:
		symbol, ok := c.scopes.SymbolResolve(node.Name)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Name)
		}
		if err := c.scopes.EmitSymbolGet(symbol); err != nil {
			return err
		}
	case *ast.CompositeLit:
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
	case *ast.FuncDecl:
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
		fnIndex := c.scopes.ConstantsAdd("", compiledFn)
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
	case *ast.CallExpr:
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
	case *ast.ReturnStmt:
		if len(node.Results) == 0 {
			// Return 'undefined'
			if _, err := c.scopes.Emit(bytecode.OpReturn, 0); err != nil {
				return err
			}
		} else {
			if err := c.Compile(node.Results[0]); err != nil {
				return err
			}
			// Return a value
			if _, err := c.scopes.Emit(bytecode.OpReturn, 1); err != nil {
				return err
			}
		}
	case *ast.SelectorExpr:
		// 1. Compila l'espressione a sinistra (es. 'fmt').
		// Questo emetterà un 'OpGetGlobal' che a runtime metterà
		// l'oggetto modulo MapImmutable sullo stack.
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

		nameIndex, found := c.scopes.ConstantsGet(cacheKey)
		if !found {
			attrArray := objects.NewArray([]objects.IObject{objects.NewStringNoSize(mName), objects.NewStringNoSize(sName)})
			nameIndex = c.scopes.ConstantsAdd(cacheKey, attrArray)
		}
		if _, err := c.scopes.Emit(bytecode.OpGetAttr, nameIndex); err != nil {
			return err
		}
	case *ast.ImportSpec:
		moduleName := node.Path.Value
		c.scopes.SymbolDefine(strings.Trim(moduleName, "\"'"))
	default:
		return fmt.Errorf("unsupported expression type: %T", node)
	}
	return nil
}

// Bytecode generates and returns the compiled bytecode representation from the compiler's current state.
// It encapsulates the main function and constants into the bytecode structure to produce an executable output.
// Returns an error if there's an issue retrieving the current compilation scope.
func (c *Compiler) Bytecode() (*bytecode.Bytecode, error) {
	bc := bytecode.NewBytecode()
	if c.mainFn == nil {
		return nil, errors.New("main function not found")
	}
	bc.SetMainFunction(c.mainFn)
	bc.SetConstants(c.scopes.ConstantsRetrieve())
	return bc, nil
}
