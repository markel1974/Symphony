package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpNotEqual)
}

// OpNotEqual is a structure representing the "not equal (!=)" opcode operation in the virtual machine.
// It embeds Opcode to provide information about the opcode, including its identifier and operands.
type OpNotEqual struct {
	*bytecode.Opcode
}

// NewOpNotEqual creates and returns a new instance of OpNotEqual with Opcode initialized from bytecode.
func NewOpNotEqual(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNotEqual{Opcode: op.Opcode(bytecode.OpNotEqual)}
}

// Execute evaluates whether the top two stack elements are unequal and pushes the result as a boolean onto the stack.
func (op *OpNotEqual) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	right := v.Stack().Pop()
	left := v.Stack().Pop()
	var val objects.IObject
	if left.Equals(right) {
		val = v.Factory().TrueValue()
	} else {
		val = v.Factory().FalseValue()
	}
	v.Stack().Push(val)
}
