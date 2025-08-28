package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/compilers/native/compiler/tables"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

type Loops struct {
	gk            objects.IGateKeeper
	fileSet       *token.FileSet
	scopes        *tables.Scopes
	structTable   *tables.StructTable
	functionTable *tables.FunctionTable
	compile       func(node ast.Node) error
}

func NewLoops(gk objects.IGateKeeper, scopes *tables.Scopes, structTable *tables.StructTable, functionTable *tables.FunctionTable) *Loops {
	return &Loops{
		gk:            gk,
		scopes:        scopes,
		structTable:   structTable,
		functionTable: functionTable,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *Loops) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes the ControlFlow structure, ensuring it is ready for subsequent compilation tasks and operations.
func (c *Loops) Prepare() error {
	return nil
}

// Compile compiles the AST nodes using the configured compile function and returns an error if the process fails.
func (c *Loops) Compile() error {
	return nil
}

// ForStmt compiles a for loop statement, including initialization, condition, post-iteration, and body execution.
func (c *Loops) ForStmt(node *ast.ForStmt) error {
	if node.Init != nil {
		if err := c.compile(node.Init); err != nil {
			return err
		}
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}

	scope.EnterLoop()

	loopStartPos := scope.InstructionsLen()
	// Compile the condition (e.g. x < 10)
	if node.Cond != nil {
		if err = c.compile(node.Cond); err != nil {
			return err
		}
	} else {
		// If there's no condition, it's an infinite loop -> emit 'true'
		if _, err = c.scopes.Emit(bytecode.OpTrue); err != nil {
			return err
		}
	}
	// Emit a conditional jump to exit the loop if the condition is false
	jumpNotTruthyPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
	if err != nil {
		return err
	}

	// Compile the loop body
	if err = c.compile(node.Body); err != nil {
		return err
	}

	// Set the target for 'continue' instructions
	continueTargetPos := scope.InstructionsLen()
	scope.CurrentLoop().ContinueTargetPosition = continueTargetPos

	// Compile the post-iteration instruction (e.g. x++) - ONLY ONCE
	if node.Post != nil {
		if err = c.compile(node.Post); err != nil {
			return err
		}
	}

	// Emit an unconditional jump back to the start of the condition
	if _, err = c.scopes.Emit(bytecode.OpJump, loopStartPos); err != nil {
		return err
	}

	scope, err = c.scopes.Current()
	if err != nil {
		return err
	}

	// Back-patching: update the conditional jump (OpJumpFalsy)
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}

	// Update all 'break' instructions
	for _, pos := range scope.CurrentLoop().BreakPositions {
		if err = c.scopes.ChangeOperand(pos, afterLoopPos); err != nil {
			return err
		}
	}

	// Update all 'continue' instructions
	for _, pos := range scope.CurrentLoop().ContinuePositions {
		if err = c.scopes.ChangeOperand(pos, scope.CurrentLoop().ContinueTargetPosition); err != nil {
			return err
		}
	}

	scope.LeaveLoop()

	return nil
}

// RangeStmt compiles a RangeStmt node into bytecode, handling iterator initialization, key/value assignment, and looping logic.
func (c *Loops) RangeStmt(node *ast.RangeStmt) error {
	if err := c.compile(node.X); err != nil {
		return err
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	scope.EnterLoop()

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
		returnTypeName, _ = c.structTable.ReturnTypeFromSymbol(expr.Name)
	case *ast.CallExpr:
		if ident, ok := expr.Fun.(*ast.Ident); ok {
			returnTypeName, _ = c.structTable.ReturnTypeFromSymbol(ident.Name)
		}
	case *ast.SelectorExpr:
		// Case: for _, v := range myVar.Items
		if receiverIdent, ok := expr.X.(*ast.Ident); ok {
			//1. Resolve the receiver symbol (myVar)
			returnTypeName, _ = c.structTable.TypeNameFromSymbolField(receiverIdent.Name, expr.Sel.Name)
		}
	default:
		return tables.NewCompilerError(c.fileSet, node, "unsupported range expression: %T", node.X)
	}
	keySymbol, err := c.functionTable.RangeKey(node)
	if err != nil {
		return err
	}
	valueSymbol, err := c.functionTable.RangeValue(node, returnTypeName)
	if err != nil {
		return err
	}

	loopStartPos := scope.InstructionsLen()

	scope.CurrentLoop().ContinueTargetPosition = loopStartPos // <-- MODIFICATION

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

	// Back-Patching
	if err = c.scopes.ChangeOperand(jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	for _, pos := range scope.CurrentLoop().BreakPositions {
		if err = c.scopes.ChangeOperand(pos, afterLoopPos); err != nil {
			return err
		}
	}
	for _, pos := range scope.CurrentLoop().ContinuePositions {
		if err = c.scopes.ChangeOperand(pos, scope.CurrentLoop().ContinueTargetPosition); err != nil {
			return err
		}
	}
	scope.LeaveLoop()
	return nil
}
