package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpSetGlobal represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpSetGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetGlobal creates and returns a new instance of OpSetGlobal with initialized OpcodeDetails.
func NewOpSetGlobal(op *bytecode.Opcodes) *OpSetGlobal {
	return &OpSetGlobal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetGlobal)}
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpSetGlobal) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	pos := decoder.Read(0)
	val := v.Stack().Peek()
	v.Globals().Set(uint(pos), val)
}
