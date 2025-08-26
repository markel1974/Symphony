package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpNotLogical)
}

// OpNotLogical represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpNotLogical struct {
	*bytecode.OpcodeDetails
}

// NewOpNotLogical creates a new instance of OpNotLogical, representing a logical NOT operation (!).
func NewOpNotLogical(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNotLogical{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNotLogical)}
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpNotLogical) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	operand := v.Stack().Pop()
	val := op.Factory().FalseValue()
	if operand.Boolean() {
		val = op.Factory().TrueValue()
	}
	v.Stack().Push(val)
}
