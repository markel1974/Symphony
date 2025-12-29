package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// Loops represents a compilation unit that processes AST nodes and maintains scope and function information.
// It stores references to the file set, scopes, and symbol tables necessary for code compilation and processing.
// The gk field is used to manage object gatekeeping, including allocation, conversion, and adaptation.
// The compile field defines a custom function to compile individual AST nodes within the context of this instance.
type Loops struct {
	gk              objects.IGateKeeper
	fileSet         *token.FileSet
	scopes          *tables.Scopes
	definitionTable *tables.DefinitionTable
	functionTable   *tables.FunctionTable
	compile         func(node ast.Node) error
}

// NewLoops creates and returns a new instance of Loops, initialized with the provided gatekeeper, scopes, and tables.
func NewLoops(gk objects.IGateKeeper, scopes *tables.Scopes, definitionTable *tables.DefinitionTable, functionTable *tables.FunctionTable) *Loops {
	return &Loops{
		gk:              gk,
		scopes:          scopes,
		definitionTable: definitionTable,
		functionTable:   functionTable,
	}
}

// Setup initializes the Loops instance with the provided fileSet and compile function.
func (c *Loops) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes internal structures and performs prerequisites required before the `Compile` method is invoked.
func (c *Loops) Prepare() error {
	return nil
}

// Compile processes the given AST node to perform compilation tasks and returns an error if the process fails.
func (c *Loops) Compile() error {
	return nil
}

// Finalize finalizes the Loops structure by performing necessary cleanup or concluding operations. Returns an error if it fails.
func (c *Loops) Finalize() error {
	return nil
}

// ForStmt compiles an AST for-loop statement by handling initialization, condition, body, post-iteration step, and loop structure.
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
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpTrueId); err != nil {
			return err
		}
	}
	// SymbolEmit a conditional jump to exit the loop if the condition is false
	jumpNotTruthyPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpFalsyId, 9999)
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

	// SymbolEmit an unconditional jump back to the start of the condition
	if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, loopStartPos); err != nil {
		return err
	}

	scope, err = c.scopes.Current()
	if err != nil {
		return err
	}

	// Back-patching: update the conditional jump (OpJumpFalsy)
	afterLoopPos := scope.InstructionsLen()
	if err = c.scopes.InstructionsChangeOperand(node.Pos(), jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}

	// Update all 'break' instructions
	for _, pos := range scope.CurrentLoop().BreakPositions {
		if err = c.scopes.InstructionsChangeOperand(node.Pos(), pos, afterLoopPos); err != nil {
			return err
		}
	}

	// Update all 'continue' instructions
	for _, pos := range scope.CurrentLoop().ContinuePositions {
		if err = c.scopes.InstructionsChangeOperand(node.Pos(), pos, scope.CurrentLoop().ContinueTargetPosition); err != nil {
			return err
		}
	}

	scope.LeaveLoop()

	return nil
}

// RangeStmt processes an AST RangeStmt, emits relevant bytecode, and handles iterator initialization and loop logic.
// It supports range expressions such as identifiers, function calls, and selector expressions.
// Errors are returned for unsupported or invalid range expressions and during bytecode generation.
// It manages loop scope, including `break` and `continue` handling, with back-patching for jump instructions.
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
	if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpIteratorInitId, iteratorSymbol.Index()); err != nil {
		return err
	}
	var returnTypeName string
	switch expr := node.X.(type) {
	case *ast.Ident:
		// ReturnTypeFromSymbol resolves the return type of a symbol by its name and returns it along with a success flag.
		if symbol, ok := c.scopes.SymbolResolve(expr.Name); ok {
			returnTypeName, _ = symbol.ReturnTypeFirst()
		}
	case *ast.CallExpr:
		if ident, ok := expr.Fun.(*ast.Ident); ok {
			if symbol, ok := c.scopes.SymbolResolve(ident.Name); ok {
				returnTypeName, _ = symbol.ReturnTypeFirst()
			}
		}
	case *ast.SelectorExpr:
		// Case: for _, v := range myVar.Items
		if receiverIdent, ok := expr.X.(*ast.Ident); ok {
			//1. Resolve the receiver symbol (myVar)
			returnTypeName, _ = c.definitionTable.StructTypeNameFromSymbolField(receiverIdent.Name, expr.Sel.Name)
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

	if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpIteratorNextId, iteratorSymbol.Index()); err != nil {
		return err
	}
	jumpNotTruthyPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpFalsyId, 9999)
	if err != nil {
		return err
	}
	if valueSymbol != nil {
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpIteratorValueId, iteratorSymbol.Index()); err != nil {
			return err
		}
		if err = c.scopes.SymbolEmitSetAndPop(node.Pos(), valueSymbol); err != nil {
			return err
		}
	}
	if keySymbol != nil {
		if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpIteratorKeyId, iteratorSymbol.Index()); err != nil {
			return err
		}
		if err = c.scopes.SymbolEmitSetAndPop(node.Pos(), keySymbol); err != nil {
			return err
		}
	}
	if err = c.compile(node.Body); err != nil {
		return err
	}
	if _, err = c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, loopStartPos); err != nil {
		return err
	}

	afterLoopPos := scope.InstructionsLen()

	// Back-Patching
	if err = c.scopes.InstructionsChangeOperand(node.Pos(), jumpNotTruthyPos, afterLoopPos); err != nil {
		return err
	}
	for _, pos := range scope.CurrentLoop().BreakPositions {
		if err = c.scopes.InstructionsChangeOperand(node.Pos(), pos, afterLoopPos); err != nil {
			return err
		}
	}
	for _, pos := range scope.CurrentLoop().ContinuePositions {
		if err = c.scopes.InstructionsChangeOperand(node.Pos(), pos, scope.CurrentLoop().ContinueTargetPosition); err != nil {
			return err
		}
	}
	scope.LeaveLoop()
	return nil
}
