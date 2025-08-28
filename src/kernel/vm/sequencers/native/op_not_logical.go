package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpNotLogical)
}

// OpNotLogical represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpNotLogical struct {
	*bytecode.Opcode
}

// NewOpNotLogical creates a new instance of OpNotLogical, representing a logical NOT operation (!).
func NewOpNotLogical(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNotLogical{Opcode: op.Opcode(bytecode.OpNotLogical)}
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpNotLogical) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	operand := v.Stack().Pop()
	var val objects.IObject
	if operand.Falsy() {
		val = v.Factory().TrueValue()
	} else {
		val = v.Factory().FalseValue()
	}
	v.Stack().Push(val)
}
