package compiler

import (
	"fmt"
)

// LoopScope represents the context of a loop, storing positions of break instructions to update their jump targets later.
type LoopScope struct {
	// Posizioni di tutte le istruzioni 'break' incontrate.
	// Queste dovranno essere aggiornate per saltare alla fine del ciclo.
	BreakPositions []int
}

// CompilationScope manages the state of a code compilation process, including instructions and loop contexts.
type CompilationScope struct {
	instructions        []byte
	lastInstruction     *EmittedInstruction
	previousInstruction *EmittedInstruction
	loopScopes          []*LoopScope
}

// NewCompilationScope initializes and returns a new instance of CompilationScope with empty instructions and loop scopes.
func NewCompilationScope() *CompilationScope {
	return &CompilationScope{
		instructions:        []byte{},
		lastInstruction:     nil,
		previousInstruction: nil,
		loopScopes:          []*LoopScope{},
	}
}

// Instructions returns the compiled bytecode instructions of the current compilation scope.
func (c *CompilationScope) Instructions() []byte {
	return c.instructions
}

// InstructionsLen returns the number of instructions in the current compilation scope.
func (c *CompilationScope) InstructionsLen() int {
	return len(c.instructions)
}

// LastInstruction returns the most recently emitted instruction in the current compilation scope.
func (c *CompilationScope) LastInstruction() *EmittedInstruction {
	return c.lastInstruction
}

// PreviousInstruction retrieves the last emitted instruction before the most recent one in the current compilation scope.
func (c *CompilationScope) PreviousInstruction() *EmittedInstruction {
	return c.previousInstruction
}

// SetLastInstruction updates the last emitted instruction in the compilation scope with the provided instruction.
func (c *CompilationScope) SetLastInstruction(instruction *EmittedInstruction) {
	c.lastInstruction = instruction
}

// SetPreviousInstruction updates the previous instruction in the compilation scope with the specified instruction.
func (c *CompilationScope) SetPreviousInstruction(instruction *EmittedInstruction) {
	c.previousInstruction = instruction
}

// InstructionsSet updates the instruction at the specified position in the instructions slice with the given byte value.
// Returns an error if the position is invalid or out of bounds.
func (c *CompilationScope) InstructionsSet(pos int, instruction byte) error {
	if pos < 0 || pos >= len(c.instructions) {
		return fmt.Errorf("invalid instruction position: %d", pos)
	}
	c.instructions[pos] = instruction
	return nil
}

// InstructionsGet retrieves the byte at the specified position in the instructions slice. Returns an error if the position is invalid.
func (c *CompilationScope) InstructionsGet(pos int) (byte, error) {
	if pos < 0 || pos >= len(c.instructions) {
		return 0, fmt.Errorf("invalid instruction position: %d", pos)
	}
	return c.instructions[pos], nil
}

// InstructionsAppend appends the given byte slice to the current scope's instructions and returns an error if unsuccessful.
func (c *CompilationScope) InstructionsAppend(instruction []byte) error {
	c.instructions = append(c.instructions, instruction...)
	return nil
}

// InstructionsReplace replaces instructions starting at the given position with the provided new instructions.
// It returns an error if the position is invalid or the replacement fails.
func (c *CompilationScope) InstructionsReplace(pos int, newInstruction []byte) error {
	for i := 0; i < len(newInstruction); i++ {
		if err := c.InstructionsSet(pos+i, newInstruction[i]); err != nil {
			return err
		}
	}
	return nil
}

// EnterLoop adds a new loop context to the stack, initializing it with an empty list of break positions.
func (c *CompilationScope) EnterLoop() {
	// Aggiunge un nuovo contesto di ciclo allo stack.
	c.loopScopes = append(c.loopScopes, &LoopScope{BreakPositions: []int{}})
}

// LeaveLoop removes the innermost loop context from the stack.
func (c *CompilationScope) LeaveLoop() {
	// Rimuove il contesto del ciclo più interno dallo stack.
	c.loopScopes = c.loopScopes[:len(c.loopScopes)-1]
}

// CurrentLoop retrieves the most recently entered loop scope from the stack or returns nil if no loop is active.
func (c *CompilationScope) CurrentLoop() *LoopScope {
	if len(c.loopScopes) == 0 {
		return nil
	}
	return c.loopScopes[len(c.loopScopes)-1]
}

// AddBreak appends a given position to the BreakPositions of the current loop scope or returns an error if not in a loop.
func (c *CompilationScope) AddBreak(pos int) error {
	loop := c.CurrentLoop()
	if loop == nil {
		return fmt.Errorf("break statement not within a loop")
	}
	loop.BreakPositions = append(loop.BreakPositions, pos)
	return nil
}
