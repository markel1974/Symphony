package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpGlobalSet represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpGlobalSet struct {
	*bytecode.OpcodeDetails
}

// NewOpGlobalSet creates and returns a new instance of OpGlobalSet with initialized OpcodeDetails.
func NewOpGlobalSet(op *bytecode.Opcodes) *OpGlobalSet {
	return &OpGlobalSet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGlobalSet)}
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpGlobalSet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	pos := decoder.Read(0)
	val := v.Stack().Peek()
	v.Globals().Set(uint(pos), val)
}
