package tables

import (
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// EmittedInstruction represents a single instruction emitted during compilation, including its opcode and position.
type EmittedInstruction struct {
	opcode   opcodes.OpcodeId
	position int
}

// NewEmittedInstruction creates a new instance of EmittedInstruction with the given opcode and position.
func NewEmittedInstruction(opcode opcodes.OpcodeId, position int) *EmittedInstruction {
	return &EmittedInstruction{
		opcode:   opcode,
		position: position,
	}
}

// Opcode returns the operation code (OpcodeId) of the emitted instruction.
func (ei *EmittedInstruction) Opcode() opcodes.OpcodeId {
	return ei.opcode
}

// Position returns the index or location of the emitted instruction within its containing sequence or bytecode stream.
func (ei *EmittedInstruction) Position() int {
	return ei.position
}
