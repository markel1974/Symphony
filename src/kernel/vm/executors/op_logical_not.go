package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpLNot represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpLNot struct {
	*bytecode.OpcodeDetails
}

// NewOpLNot creates a new instance of OpLNot, representing a logical NOT operation (!).
func NewOpLNot(op *bytecode.Opcodes) *OpLNot {
	return &OpLNot{OpcodeDetails: op.OpcodeToDetails(bytecode.OpLNot)}
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpLNot) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	operand := v.Stack().Pop()
	val := op.Factory().FalseValue()
	if operand.Boolean() {
		val = op.Factory().TrueValue()
	}
	v.Stack().Push(val)
}
