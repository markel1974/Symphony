package compiler

import "fmt"

// CompilationScope represents a compilation context for bytecode with instructions and metadata about emissions.
type CompilationScope struct {
	instructions        []byte
	lastInstruction     *EmittedInstruction
	previousInstruction *EmittedInstruction
}

// NewCompilationScope creates and returns a new instance of CompilationScope with initialized fields.
func NewCompilationScope() *CompilationScope {
	return &CompilationScope{
		instructions:        []byte{},
		lastInstruction:     nil,
		previousInstruction: nil,
	}
}

// Instructions returns the compiled bytecode instructions for the current compilation scope.
func (c *CompilationScope) Instructions() []byte {
	return c.instructions
}

// InstructionsLen returns the total number of instructions in the current compilation scope.
func (c *CompilationScope) InstructionsLen() int {
	return len(c.instructions)
}

// LastInstruction returns the most recently emitted bytecode instruction in the current compilation scope.
func (c *CompilationScope) LastInstruction() *EmittedInstruction {
	return c.lastInstruction
}

// PreviousInstruction returns the previously emitted instruction in the compilation scope.
func (c *CompilationScope) PreviousInstruction() *EmittedInstruction {
	return c.previousInstruction
}

// SetLastInstruction sets the last instruction in the current compilation scope to the provided emitted instruction.
func (c *CompilationScope) SetLastInstruction(instruction *EmittedInstruction) {
	c.lastInstruction = instruction
}

// SetPreviousInstruction assigns the provided EmittedInstruction as the previousInstruction in the current compilation scope.
func (c *CompilationScope) SetPreviousInstruction(instruction *EmittedInstruction) {
	c.previousInstruction = instruction
}

// InstructionsSet modifies the instruction at the specified position with the given byte value. Returns an error if the position is invalid.
func (c *CompilationScope) InstructionsSet(pos int, instruction byte) error {
	if pos < 0 || pos >= len(c.instructions) {
		return fmt.Errorf("invalid instruction position: %d", pos)
	}
	c.instructions[pos] = instruction
	return nil
}

// InstructionsGet retrieves the instruction byte at the specified position within the instructions slice.
// Returns an error if the position is out of bounds.
func (c *CompilationScope) InstructionsGet(pos int) (byte, error) {
	if pos < 0 || pos >= len(c.instructions) {
		return 0, fmt.Errorf("invalid instruction position: %d", pos)
	}
	return c.instructions[pos], nil
}

// InstructionsAppend appends the given byte instruction to the instructions of the compilation scope. It may return an error.
func (c *CompilationScope) InstructionsAppend(instruction []byte) error {
	c.instructions = append(c.instructions, instruction...)
	return nil
}

// InstructionsReplace replaces the instruction at the specified position with the given byte instruction.
func (c *CompilationScope) InstructionsReplace(pos int, newInstruction []byte) error {
	for i := 0; i < len(newInstruction); i++ {
		if err := c.InstructionsSet(pos+i, newInstruction[i]); err != nil {
			return err
		}
	}
	return nil
}
