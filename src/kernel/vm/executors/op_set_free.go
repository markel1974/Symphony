package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpSetFree represents an operation to set the value of a free variable within a closure's environment.
type OpSetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpSetFree creates and returns a new instance of OpSetFree initialized with its corresponding OpcodeDetails.
func NewOpSetFree(op *bytecode.Opcodes) *OpSetFree {
	return &OpSetFree{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetFree)}
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpSetFree) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	o := v.Stack().Pop()
	v.FreeVarsIndex(freeIndex).SetValue(o)
}
