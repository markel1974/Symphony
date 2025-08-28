package compiler

import (
	"go/ast"
	"go/token"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// ControlFlow is a structure used to manage and compile AST nodes with the associated file sets and scope information.
// It leverages an IGateKeeper for managing object interactions during the compilation process.
type ControlFlow struct {
	gk      objects.IGateKeeper
	fileSet *token.FileSet
	scopes  *Scopes
	compile func(node ast.Node) error
}

// NewControlFlow creates and returns a new instance of ControlFlow with the specified gatekeeper and scope parameters.
func NewControlFlow(gk objects.IGateKeeper, scopes *Scopes) *ControlFlow {
	return &ControlFlow{
		gk:     gk,
		scopes: scopes,
	}
}

// Setup initializes the Declarations object with a file set and a compile function, returning an error if any occur.
func (c *ControlFlow) Setup(fileSet *token.FileSet, compile func(node ast.Node) error) error {
	c.fileSet = fileSet
	c.compile = compile
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
				return NewCompilerError(c.fileSet, node, err.Error())
			}
		} else if scope.CurrentLoop() != nil { // Otherwise, check if we're in a loop
			breakJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
			if err != nil {
				return err
			}
			if err = scope.AddBreak(breakJumpPos); err != nil {
				return NewCompilerError(c.fileSet, node, err.Error())
			}
		} else {
			return NewCompilerError(c.fileSet, node, "break statement not within a loop or switch")
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
			return NewCompilerError(c.fileSet, node, err.Error())
		}
		return nil
	}

	return NewCompilerError(c.fileSet, node, "unsupported branch statement: %s", node.Tok.String())
}

// SwitchStmt processes a given *ast.SwitchStmt node, handles scopes, and generates bytecode for a switch statement.
// It manages compilation of the switch tag, case clauses, default clause, and handles necessary jump instructions.
func (c *ControlFlow) SwitchStmt(node *ast.SwitchStmt) error {
	scope, err := c.scopes.Current()
	if err != nil {
		return err
	}
	scope.EnterSwitch()

	// 1. Compila l'espressione (tag) e salvala in una variabile temporanea
	//    per evitare di ricalcolarla. Questo è cruciale.
	var tagSymbol *Symbol
	if node.Tag != nil {
		if err := c.compile(node.Tag); err != nil {
			return err
		}
		tagSymbol, err = c.scopes.SymbolDefineUnique("__switch_tag")
		if err != nil {
			return err
		}
		if err := c.scopes.EmitSymbolDefine(tagSymbol); err != nil {
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpPop); err != nil {
			return err
		}
	}

	//if err := c.scopes.Enter("switch_body"); err != nil {
	//	return err
	//}

	var defaultClause *ast.CaseClause
	jumpsToNextCase := []int{}

	// 2. Itera su tutti i 'case' (tranne il 'default', che gestiamo alla fine)
	for _, clauseStmt := range node.Body.List {
		clause := clauseStmt.(*ast.CaseClause)

		// Mettiamo da parte il 'default' per dopo
		if clause.List == nil {
			defaultClause = clause
			continue
		}

		// Back-patch: collega i salti del case precedente a questo
		afterPreviousCasePos := scope.InstructionsLen()
		for _, pos := range jumpsToNextCase {
			if err := c.scopes.ChangeOperand(pos, afterPreviousCasePos); err != nil {
				return err
			}
		}
		jumpsToNextCase = []int{} // Resetta la lista per il case corrente

		// 3. Compila la condizione del case (es. tag == val1)
		if tagSymbol != nil {
			if err := c.scopes.EmitSymbolGet(tagSymbol); err != nil {
				return err
			}
		}
		if err := c.compile(clause.List[0]); err != nil { // Semplificato: assume un valore per case
			return err
		}
		if _, err := c.scopes.Emit(bytecode.OpEqual); err != nil {
			return err
		}

		// 4. Salta al prossimo case se la condizione è falsa
		jumpPos, err := c.scopes.Emit(bytecode.OpJumpFalsy, 9999)
		if err != nil {
			return err
		}
		jumpsToNextCase = append(jumpsToNextCase, jumpPos)

		// 5. Compila il corpo del case
		for _, stmt := range clause.Body {
			if err := c.compile(stmt); err != nil {
				return err
			}
		}

		// 6. Aggiungi un salto alla fine dello switch (Go non ha fall-through)
		endJumpPos, err := c.scopes.Emit(bytecode.OpJump, 9999)
		if err != nil {
			return err
		}
		if err := scope.AddEndJump(endJumpPos); err != nil {
			return err
		}
	}

	// 7. Compila il 'default' se esiste
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

	//if _, err := c.scopes.Leave(); err != nil {
	//	return err
	//}

	// 8. Back-patch finale: aggiorna tutti i salti alla fine
	afterSwitchPos := scope.InstructionsLen()
	for _, pos := range scope.CurrentSwitch().EndJumps {
		if err := c.scopes.ChangeOperand(pos, afterSwitchPos); err != nil {
			return err
		}
	}

	scope.LeaveSwitch()
	return nil
}
