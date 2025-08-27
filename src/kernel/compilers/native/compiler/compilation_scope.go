package compiler

import (
	"fmt"
)

// SwitchScope represents the scope of a switch statement during compilation, storing positions for end jump instructions.
type SwitchScope struct {
	EndJumps []int
}

// LoopScope represents the scope of a loop during compilation, tracking break and continue positions within the loop.
type LoopScope struct {
	BreakPositions         []int
	ContinuePositions      []int
	ContinueTargetPosition int
}

// CompilationScope represents the current state of compilation including instructions, emitted instructions,
// loop and switch constructs. It facilitates tracking compilation progress and managing break/continue semantics.
type CompilationScope struct {
	instructions        []byte
	lastInstruction     *EmittedInstruction
	previousInstruction *EmittedInstruction
	loopScopes          []*LoopScope
	switchScopes        []*SwitchScope
}

// NewCompilationScope initializes and returns a new instance of CompilationScope with default values.
func NewCompilationScope() *CompilationScope {
	return &CompilationScope{
		instructions:        []byte{},
		lastInstruction:     nil,
		previousInstruction: nil,
		loopScopes:          []*LoopScope{},
		switchScopes:        []*SwitchScope{},
	}
}

// Instructions returns the bytecode instructions currently stored in the compilation scope.
func (c *CompilationScope) Instructions() []byte {
	return c.instructions
}

// InstructionsLen returns the length of the instructions slice in the current CompilationScope.
func (c *CompilationScope) InstructionsLen() int {
	return len(c.instructions)
}

// LastInstruction returns the last emitted instruction in the current compilation scope.
func (c *CompilationScope) LastInstruction() *EmittedInstruction {
	return c.lastInstruction
}

// PreviousInstruction returns the previously emitted instruction in the current compilation scope.
func (c *CompilationScope) PreviousInstruction() *EmittedInstruction {
	return c.previousInstruction
}

// SetLastInstruction sets the last emitted instruction for this compilation scope.
func (c *CompilationScope) SetLastInstruction(instruction *EmittedInstruction) {
	c.lastInstruction = instruction
}

// SetPreviousInstruction sets the previous emitted instruction for the current compilation scope.
func (c *CompilationScope) SetPreviousInstruction(instruction *EmittedInstruction) {
	c.previousInstruction = instruction
}

// InstructionsSet updates the instruction at the specified position in the current scope with the given byte value.
// Returns an error if the position is invalid or out of bounds.
func (c *CompilationScope) InstructionsSet(pos int, instruction byte) error {
	if pos < 0 || pos >= len(c.instructions) {
		return fmt.Errorf("invalid instruction position: %d", pos)
	}
	c.instructions[pos] = instruction
	return nil
}

// InstructionsGet retrieves the instruction at the specified position in the instructions slice of the compilation scope.
// Returns the instruction as a byte or an error if the position is invalid (e.g., out of bounds).
func (c *CompilationScope) InstructionsGet(pos int) (byte, error) {
	if pos < 0 || pos >= len(c.instructions) {
		return 0, fmt.Errorf("invalid instruction position: %d", pos)
	}
	return c.instructions[pos], nil
}

// InstructionsAppend appends a byte slice to the instructions of the current compilation scope and returns any errors encountered.
func (c *CompilationScope) InstructionsAppend(instruction []byte) error {
	c.instructions = append(c.instructions, instruction...)
	return nil
}

// InstructionsReplace replaces a slice of instructions starting at the specified position with the provided new instructions.
func (c *CompilationScope) InstructionsReplace(pos int, newInstruction []byte) error {
	for i := 0; i < len(newInstruction); i++ {
		if err := c.InstructionsSet(pos+i, newInstruction[i]); err != nil {
			return err
		}
	}
	return nil
}

// EnterLoop appends a new loop context onto the stack, initializing break positions for the loop.
func (c *CompilationScope) EnterLoop() {
	// Aggiunge un nuovo contesto di ciclo allo stack.
	c.loopScopes = append(c.loopScopes, &LoopScope{BreakPositions: []int{}})
}

// LeaveLoop removes the innermost loop context from the stack of loop scopes.
func (c *CompilationScope) LeaveLoop() {
	// Rimuove il contesto del ciclo più interno dallo stack.
	c.loopScopes = c.loopScopes[:len(c.loopScopes)-1]
}

// CurrentLoop retrieves the most recently entered loop scope or returns nil if no loop scope exists.
func (c *CompilationScope) CurrentLoop() *LoopScope {
	if len(c.loopScopes) == 0 {
		return nil
	}
	return c.loopScopes[len(c.loopScopes)-1]
}

// AddBreak adds the position of a break statement to the current loop's break positions. Returns an error if no loop exists.
func (c *CompilationScope) AddBreak(pos int) error {
	loop := c.CurrentLoop()
	if loop == nil {
		return fmt.Errorf("break statement not within a loop")
	}
	loop.BreakPositions = append(loop.BreakPositions, pos)
	return nil
}

// AddContinue adds the given position to the current loop's continue positions if inside a valid loop, otherwise returns an error.
func (c *CompilationScope) AddContinue(pos int) error {
	loop := c.CurrentLoop()
	if loop == nil {
		return fmt.Errorf("continue statement not within a loop")
	}
	loop.ContinuePositions = append(loop.ContinuePositions, pos)
	return nil
}

// EnterSwitch adds a new switch context to the stack of switch scopes.
func (c *CompilationScope) EnterSwitch() {
	c.switchScopes = append(c.switchScopes, &SwitchScope{EndJumps: []int{}})
}

// LeaveSwitch removes the most recent switch scope from the stack of switch scopes.
func (c *CompilationScope) LeaveSwitch() {
	c.switchScopes = c.switchScopes[:len(c.switchScopes)-1]
}

// CurrentSwitch returns the most recently entered SwitchScope, or nil if no switch scopes exist in the CompilationScope.
func (c *CompilationScope) CurrentSwitch() *SwitchScope {
	if len(c.switchScopes) == 0 {
		return nil
	}
	return c.switchScopes[len(c.switchScopes)-1]
}

// AddEndJump adds the position of an end jump to the current switch scope. Returns an error if not within a switch context.
func (c *CompilationScope) AddEndJump(pos int) error {
	s := c.CurrentSwitch()
	if s == nil {
		return fmt.Errorf("break or case statement not within a switch")
	}
	s.EndJumps = append(s.EndJumps, pos)
	return nil
}
