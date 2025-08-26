package compiler

import "github.com/markel1974/c64emu/src/kernel/vm/bytecode"

// EmittedInstruction represents a bytecode instruction emitted during compilation with its opcode and position metadata.
type EmittedInstruction struct {
	opcode   bytecode.OpcodeId
	position int
}

// NewEmittedInstruction creates a new EmittedInstruction with the provided opcode and position.
func NewEmittedInstruction(opcode bytecode.OpcodeId, position int) *EmittedInstruction {
	return &EmittedInstruction{
		opcode:   opcode,
		position: position,
	}
}
