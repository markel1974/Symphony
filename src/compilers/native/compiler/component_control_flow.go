package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/compilers/native/tables"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/sequencers/native"
)

// ControlFlow is a structure used to manage and compile AST nodes with the associated file sets and scope information.
// It leverages an IGateKeeper for managing object interactions during the compilation process.
type ControlFlow struct {
	gk              objects.IGateKeeper
	fileSet         *token.FileSet
	scopes          *tables.Scopes
	definitionTable *tables.DefinitionTable
	constants       *tables.Constants
	compile         func(node ast.Node) error
}

// NewControlFlow creates and returns a new instance of ControlFlow with the specified gatekeeper and scope parameters.
func NewControlFlow(gk objects.IGateKeeper, constants *tables.Constants, scopes *tables.Scopes, definitionTable *tables.DefinitionTable) *ControlFlow {
	return &ControlFlow{
		gk:              gk,
		constants:       constants,
		scopes:          scopes,
		definitionTable: definitionTable,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *ControlFlow) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
	return nil
}

// Prepare initializes the ControlFlow structure, ensuring it is ready for subsequent compilation tasks and operations.
func (c *ControlFlow) Prepare() error {
	return nil
}

// Compile compiles the AST nodes using the configured compile function and returns an error if the process fails.
func (c *ControlFlow) Compile() error {
	return nil
}

// IfStmt compiles an if statement, handling both 'then' and optional 'else' blocks with associated bytecode instructions.
func (c *ControlFlow) IfStmt(node *ast.IfStmt) error {
	if err := c.compile(node.Cond); err != nil {
		return tables.NewCompilerError(c.fileSet, node, err.Error())
	}
	// emit conditional jump with temporary address
	jumpNotTruthyPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpFalsyId, 9999)
	if err != nil {
		return tables.NewCompilerError(c.fileSet, node, err.Error())
	}
	// compile 'then' block
	if err = c.compile(node.Body); err != nil {
		return tables.NewCompilerError(c.fileSet, node, err.Error())
	}
	// if there's an 'else' block, emit jump to skip it
	jumpToEndPos := 0
	if node.Else != nil {
		jumpToEndPos, err = c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	scope, err := c.scopes.Current()
	if err != nil {
		return tables.NewCompilerError(c.fileSet, node, err.Error())
	}
	// update conditional jump address
	if err = c.scopes.InstructionsChangeOperand(node.Pos(), jumpNotTruthyPos, scope.InstructionsLen()); err != nil {
		return tables.NewCompilerError(c.fileSet, node, err.Error())
	}
	// compile 'else' block if it exists
	if node.Else != nil {
		if err = c.compile(node.Else); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		scope, err = c.scopes.Current()
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		// update jump address to skip else
		if err = c.scopes.InstructionsChangeOperand(node.Pos(), jumpToEndPos, scope.InstructionsLen()); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	return nil
}

// BranchStmt compiles a branch statement, handling 'break' and 'continue' operations within loops or switch statements.
func (c *ControlFlow) BranchStmt(node *ast.BranchStmt) error {
	scope, scopeErr := c.scopes.Current()
	if scopeErr != nil {
		return tables.NewCompilerError(c.fileSet, node, scopeErr.Error())
	}
	if node.Tok == token.BREAK {
		if scope.CurrentSwitch() != nil {
			breakJumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
			if err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			if err = scope.AddEndJump(breakJumpPos); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		} else if scope.CurrentLoop() != nil { // Otherwise, check if we're in a loop
			breakJumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
			if err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			if err = scope.AddBreak(breakJumpPos); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		} else {
			return tables.NewCompilerError(c.fileSet, node, "break statement not within a loop or switch")
		}
		return nil
	}

	if node.Tok == token.CONTINUE {
		// SymbolEmit an unconditional jump with a temporary address
		continueJumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		// add the position of this 'continue' to the loop context
		if err = scope.AddContinue(continueJumpPos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		return nil
	}

	return tables.NewCompilerError(c.fileSet, node, "unsupported branch statement: %s", node.Tok.String())
}

// SwitchStmt processes a given *ast.SwitchStmt node, handles scopes, and generates bytecode for a switch statement.
// It manages compilation of the switch tag, case clauses, default clause, and handles necessary jump instructions.
func (c *ControlFlow) SwitchStmt(node *ast.SwitchStmt) error {
	scope, scopeErr := c.scopes.Current()
	if scopeErr != nil {
		return tables.NewCompilerError(c.fileSet, node, scopeErr.Error())
	}
	scope.EnterSwitch()

	// 1. Compile the expression (tag) and store it in a temporary variable to avoid recalculation.
	var tagSymbol *tables.Symbol
	if node.Tag != nil {
		err := c.compile(node.Tag)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		tagSymbol, err = c.scopes.SymbolDefineUnique("__switch_tag_")
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		if err = c.scopes.SymbolEmitDefineAndPop(node.Pos(), tagSymbol); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}

	var defaultClause *ast.CaseClause
	var jumpsToNextCase []int
	// 2. Iterate over all 'case' statements (except 'default', which we handle at the end)
	for _, clauseStmt := range node.Body.List {
		clause := clauseStmt.(*ast.CaseClause)
		for _, clauseEntry := range clause.List {
			// Save the default clause for later
			if clause.List == nil {
				defaultClause = clause
				continue
			}
			// Back-patch: connect previous case jumps to this one
			afterPreviousCasePos := scope.InstructionsLen()
			for _, pos := range jumpsToNextCase {
				if err := c.scopes.InstructionsChangeOperand(node.Pos(), pos, afterPreviousCasePos); err != nil {
					return tables.NewCompilerError(c.fileSet, node, err.Error())
				}
			}
			jumpsToNextCase = []int{}
			// 3. Compile case condition (e.g. tag == val1)
			if tagSymbol != nil {
				if err := c.scopes.SymbolEmitGet(node.Pos(), tagSymbol); err != nil {
					return tables.NewCompilerError(c.fileSet, node, err.Error())
				}
			}
			if err := c.compile(clauseEntry); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			eql, ok := BinaryAdapterFor(token.EQL)
			if !ok {
				return tables.NewCompilerError(c.fileSet, node, "unhandled binary op: %s", token.EQL)
			}
			if _, err := c.scopes.SymbolEmit(node.Pos(), eql.op, eql.arguments...); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			// 4. Jump to the next-case if the condition is false
			jumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpFalsyId, 9999)
			if err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			jumpsToNextCase = append(jumpsToNextCase, jumpPos)
			// 5. Compile case body
			for _, stmt := range clause.Body {
				if err = c.compile(stmt); err != nil {
					return tables.NewCompilerError(c.fileSet, node, err.Error())
				}
			}
			// 6. add jump to end of switch (Go has no fall-through)
			endJumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
			if err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			if err = scope.AddEndJump(endJumpPos); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		}
	}
	// 7. Compile 'default' if it exists
	afterLastCasePos := scope.InstructionsLen()
	for _, pos := range jumpsToNextCase {
		if err := c.scopes.InstructionsChangeOperand(node.Pos(), pos, afterLastCasePos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	if defaultClause != nil {
		for _, stmt := range defaultClause.Body {
			if err := c.compile(stmt); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		}
	}
	// 8. Final back-patch: update all jumps to end
	afterSwitchPos := scope.InstructionsLen()
	for _, pos := range scope.CurrentSwitch().EndJumps {
		if err := c.scopes.InstructionsChangeOperand(node.Pos(), pos, afterSwitchPos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	scope.LeaveSwitch()
	return nil
}

// TypeSwitchStmt processes a type switch statement, handling case clauses, default clause, and generating necessary bytecode.
func (c *ControlFlow) TypeSwitchStmt(node *ast.TypeSwitchStmt) error {
	scope, scopeErr := c.scopes.Current()
	if scopeErr != nil {
		return tables.NewCompilerError(c.fileSet, node, scopeErr.Error())
	}
	scope.EnterSwitch()

	assignStmt, assignStmtOk := node.Assign.(*ast.AssignStmt)
	if !assignStmtOk {
		return tables.NewCompilerError(c.fileSet, node, "unexpected assign statement for type switch %w", node.Assign)
	}
	if len(assignStmt.Rhs) == 0 {
		return tables.NewCompilerError(c.fileSet, node, "empty right hand side for type switch")
	}
	typeAssertExpr, typeAssertExprOk := assignStmt.Rhs[0].(*ast.TypeAssertExpr)
	if !typeAssertExprOk {
		return tables.NewCompilerError(c.fileSet, node, "unexpected assign statement for type switch %w", assignStmt.Rhs[0])
	}
	interfaceExpr := typeAssertExpr.X

	var jumpsToEnd []int
	var jumpToNextCasePos = -1

	// iterate over all 'case' clauses (excluding 'default')
	for _, clauseStmt := range node.Body.List {
		clause, ok := clauseStmt.(*ast.CaseClause)
		if !ok || clause.List == nil {
			continue
		}

		if jumpToNextCasePos != -1 {
			if err := c.scopes.InstructionsChangeOperand(node.Pos(), jumpToNextCasePos, scope.InstructionsLen()); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		}

		if err := c.compile(interfaceExpr); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		targetTypeName := clause.List[0].(*ast.Ident).Name
		hasOk := 1 //default emit hasOk
		constIndex := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, targetTypeName))
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpTypeAssertId, hasOk, constIndex); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		var err error
		jumpToNextCasePos, err = c.scopes.SymbolEmit(node.Pos(), native.OpJumpFalsyId, 9999)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}

		// success path
		if err = c.scopes.Enter(tables.LocalScope, ""); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		variableName := assignStmt.Lhs[0].(*ast.Ident).Name
		caseVarSymbol, err := c.scopes.SymbolDefine(variableName)
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		c.definitionTable.TypeAssign(caseVarSymbol, targetTypeName)

		if err = c.scopes.SymbolEmitSetAndPop(node.Pos(), caseVarSymbol); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		for _, stmt := range clause.Body {
			if err = c.compile(stmt); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		}
		// 1. Before leaving the scope, save the last instruction generated within it.
		innerScope, err := c.scopes.Current()
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		lastInstructionInCase := innerScope.LastInstruction()

		// 2. Now leave the scope and get the generated bytecode.
		caseBytecode, source, err := c.scopes.Leave()
		if err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}

		// 3. add case bytecode to outer scope.
		if _, err = c.scopes.InstructionsAppend(caseBytecode, source); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}

		// 4. Now, in the outer scope, check the instruction we saved.
		// If it's not a 'return', add the jump to end.
		if lastInstructionInCase == nil || lastInstructionInCase.Opcode() != native.OpReturnId {
			jumpPos, err := c.scopes.SymbolEmit(node.Pos(), native.OpJumpId, 9999)
			if err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
			jumpsToEnd = append(jumpsToEnd, jumpPos)
		}
	}

	// Final landing for last failed case
	if jumpToNextCasePos != -1 {
		if err := c.scopes.InstructionsChangeOperand(node.Pos(), jumpToNextCasePos, scope.InstructionsLen()); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}

	// Clean the stack from last remaining 'result' (undefined)
	if len(node.Body.List) > 0 {
		if _, err := c.scopes.SymbolEmit(node.Pos(), native.OpPopId); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}

	// Compile 'default' if exists
	if defaultClause := c.findDefault(node.Body); defaultClause != nil {
		for _, stmt := range defaultClause.Body {
			if err := c.compile(stmt); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		}
	}

	// Final back-patching of all jumps to end
	endPos := scope.InstructionsLen()
	for _, pos := range jumpsToEnd {
		if err := c.scopes.InstructionsChangeOperand(node.Pos(), pos, endPos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}
	for _, pos := range scope.CurrentSwitch().EndJumps {
		if err := c.scopes.InstructionsChangeOperand(node.Pos(), pos, endPos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
	}

	scope.LeaveSwitch()
	return nil
}

// findDefault traverses a block statement to find the first default case clause and returns it, or nil if none exists.
func (c *ControlFlow) findDefault(list *ast.BlockStmt) *ast.CaseClause {
	for _, clauseStmt := range list.List {
		if clause, ok := clauseStmt.(*ast.CaseClause); ok {
			if clause.List == nil {
				return clause
			}
		}
	}
	return nil
}
