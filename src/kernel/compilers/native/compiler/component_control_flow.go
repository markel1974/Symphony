package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/compilers/native/tables"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// ControlFlow is a structure used to manage and compile AST nodes with the associated file sets and scope information.
// It leverages an IGateKeeper for managing object interactions during the compilation process.
type ControlFlow struct {
	gk          objects.IGateKeeper
	fileSet     *token.FileSet
	scopes      *tables.Scopes
	structTable *tables.StructTable
	constants   *Constants
	compile     func(node ast.Node) error
}

// NewControlFlow creates and returns a new instance of ControlFlow with the specified gatekeeper and scope parameters.
func NewControlFlow(gk objects.IGateKeeper, constants *Constants, scopes *tables.Scopes, structTable *tables.StructTable) *ControlFlow {
	return &ControlFlow{
		gk:          gk,
		constants:   constants,
		scopes:      scopes,
		structTable: structTable,
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

// BranchStmt compiles a branch statement, handling 'break' and 'continue' operations within loops or switch statements.
func (c *ControlFlow) BranchStmt(node *ast.BranchStmt) error {
	scope, scopeErr := c.scopes.Current()
	if scopeErr != nil {
		return scopeErr
	}
	if node.Tok == token.BREAK {
		if scope.CurrentSwitch() != nil {
			breakJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
			if err != nil {
				return err
			}
			if err = scope.AddEndJump(breakJumpPos); err != nil {
				return tables.NewCompilerError(c.fileSet, node, err.Error())
			}
		} else if scope.CurrentLoop() != nil { // Otherwise, check if we're in a loop
			breakJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
			if err != nil {
				return err
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
		// Emit an unconditional jump with a temporary address
		continueJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
		if err != nil {
			return err
		}
		// Add the position of this 'continue' to the loop context
		if err := scope.AddContinue(continueJumpPos); err != nil {
			return tables.NewCompilerError(c.fileSet, node, err.Error())
		}
		return nil
	}

	return tables.NewCompilerError(c.fileSet, node, "unsupported branch statement: %s", node.Tok.String())
}

// SwitchStmt processes a given *ast.SwitchStmt node, handles scopes, and generates bytecode for a switch statement.
// It manages compilation of the switch tag, case clauses, default clause, and handles necessary jump instructions.
func (c *ControlFlow) SwitchStmt(node *ast.SwitchStmt) error {
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	scope.EnterSwitch()

	// 1. Compile the expression (tag) and store it in a temporary variable to avoid recalculation.
	var tagSymbol *tables.Symbol
	if node.Tag != nil {
		if err := c.compile(node.Tag); err != nil {
			return err
		}
		tagSymbol, err = c.scopes.SymbolDefineUnique("__switch_tag")
		if err != nil {
			return err
		}
		if err = c.scopes.EmitSymbolDefineAndPop(tagSymbol); err != nil {
			return err
		}
	}

	var defaultClause *ast.CaseClause
	var jumpsToNextCase []int
	// 2. Iterate over all 'case' statements (except 'default', which we handle at the end)
	for _, clauseStmt := range node.Body.List {
		clause := clauseStmt.(*ast.CaseClause)
		// Save the default clause for later
		if clause.List == nil {
			defaultClause = clause
			continue
		}
		// Back-patch: connect previous case jumps to this one
		afterPreviousCasePos := scope.InstructionsLen()
		for _, pos := range jumpsToNextCase {
			if err := c.scopes.ChangeOperand(pos, afterPreviousCasePos); err != nil {
				return err
			}
		}
		jumpsToNextCase = []int{} // Reset list for current-case
		// 3. Compile case condition (e.g. tag == val1)
		if tagSymbol != nil {
			if err := c.scopes.EmitSymbolGet(tagSymbol); err != nil {
				return err
			}
		}
		if err := c.compile(clause.List[0]); err != nil { // Simplified: assumes one value per case
			return err
		}
		eql, ok := BinaryAdapterFor(token.EQL)
		if !ok {
			return tables.NewCompilerError(c.fileSet, node, "unhandled binary op: %s", token.EQL)
		}
		if _, err := c.scopes.Emit(eql.op, eql.arguments...); err != nil {
			return err
		}
		//if _, err := c.scopes.Emit(bytecode.OpLogical, int(objects.OperatorLogicalEq)); err != nil {
		//	return err
		//}
		// 4. Jump to the next-case if the condition is false
		jumpPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
		if err != nil {
			return err
		}
		jumpsToNextCase = append(jumpsToNextCase, jumpPos)
		// 5. Compile case body
		for _, stmt := range clause.Body {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
		// 6. Add jump to end of switch (Go has no fall-through)
		endJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
		if err != nil {
			return err
		}
		if err := scope.AddEndJump(endJumpPos); err != nil {
			return err
		}
	}
	// 7. Compile 'default' if it exists
	afterLastCasePos := scope.InstructionsLen()
	for _, pos := range jumpsToNextCase {
		if err := c.scopes.ChangeOperand(pos, afterLastCasePos); err != nil {
			return err
		}
	}
	if defaultClause != nil {
		for _, stmt := range defaultClause.Body {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
	}
	// 8. Final back-patch: update all jumps to end
	afterSwitchPos := scope.InstructionsLen()
	for _, pos := range scope.CurrentSwitch().EndJumps {
		if err := c.scopes.ChangeOperand(pos, afterSwitchPos); err != nil {
			return err
		}
	}
	scope.LeaveSwitch()
	return nil
}

// TypeSwitchStmt processes a given *ast.TypeSwitchStmt node, handles scopes, and generates bytecode for a type switch statement.
// It manages compilation of the type assertion, case clauses, default clause, and handles necessary jump instructions.
func (c *ControlFlow) TypeSwitchStmt(node *ast.TypeSwitchStmt) error {
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	scope.EnterSwitch()
	// 1. Compile the interface object (e.g. 'i' in 'i.(type)') and save it in a temp variable to access in each case.
	if err = c.compile(node.Assign.(*ast.AssignStmt).Rhs[0].(*ast.TypeAssertExpr).X); err != nil {
		return err
	}
	interfaceSymbol, err := c.scopes.SymbolDefineUnique("__type_switch_interface")
	if err != nil {
		return err
	}
	if err = c.scopes.EmitSymbolDefine(interfaceSymbol); err != nil {
		return err
	}
	var defaultClause *ast.CaseClause
	var jumpsToNextCase []int
	var endJumps []int

	// 2. Iterate over all 'case' clauses
	for _, clauseStmt := range node.Body.List {
		clause := clauseStmt.(*ast.CaseClause)
		// Save default for later
		if clause.List == nil {
			defaultClause = clause
			continue
		}
		// Back-patch jumps from previous case to this one
		afterPreviousCasePos := scope.InstructionsLen()
		for _, pos := range jumpsToNextCase {
			if err := c.scopes.ChangeOperand(pos, afterPreviousCasePos); err != nil {
				return err
			}
		}
		jumpsToNextCase = []int{}
		// 3. Execute type assertion for this case
		// Load interface from temporary variable
		if err := c.scopes.EmitSymbolGet(interfaceSymbol); err != nil {
			return err
		}
		// Get type name to check against (e.g. "int", "string")
		targetTypeName := clause.List[0].(*ast.Ident).Name
		constIndex := c.constants.AddOrGet("", c.gk.NewString(objects.FrameStatic, targetTypeName))
		if _, err := c.scopes.Emit(bytecode.OpTypeAssert, constIndex); err != nil {
			return err
		}
		// Stack now contains [converted_value, success_boolean]
		// Use boolean for conditional jump
		jumpPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
		if err != nil {
			return err
		}
		jumpsToNextCase = append(jumpsToNextCase, jumpPos)
		// 4. If type matches, define new variable (e.g. 'v')
		//	in a new scope for case body
		if err := c.scopes.Enter(tables.UnknownScope, ""); err != nil {
			return err
		}
		assignStmt := node.Assign.(*ast.AssignStmt)
		variableName := assignStmt.Lhs[0].(*ast.Ident).Name
		// Define symbol 'v' with correct type
		caseVarSymbol, err := c.scopes.SymbolDefine(variableName)
		if err != nil {
			return err
		}
		// Type inference for symbol
		if err = c.structTable.AssignSymbol(caseVarSymbol, targetTypeName, []string{targetTypeName}); err != nil {
			return err
		}
		// Assign value (already on stack) to new variable
		if err = c.scopes.EmitSymbolSetAndPop(caseVarSymbol); err != nil {
			return err
		}
		// 5. Compile case body
		for _, stmt := range clause.Body {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
		// Exit case scope
		if _, err = c.scopes.Leave(); err != nil {
			return err
		}
		// Add jump to end of switch
		endJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
		if err != nil {
			return err
		}
		endJumps = append(endJumps, endJumpPos)
	}
	// 6. Compile 'default' if it exists
	afterLastCasePos := scope.InstructionsLen()
	for _, pos := range jumpsToNextCase {
		if err := c.scopes.ChangeOperand(pos, afterLastCasePos); err != nil {
			return err
		}
	}
	if defaultClause != nil {
		// In default, converted value not needed
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
		for _, stmt := range defaultClause.Body {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}
	}
	// 7. Final back-patch
	afterSwitchPos := scope.InstructionsLen()
	for _, pos := range endJumps {
		if err := c.scopes.ChangeOperand(pos, afterSwitchPos); err != nil {
			return err
		}
	}
	scope.LeaveSwitch()
	return nil
}
