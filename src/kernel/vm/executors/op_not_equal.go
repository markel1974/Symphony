package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpNotEqual is a structure representing the "not equal (!=)" opcode operation in the virtual machine.
// It embeds OpcodeDetails to provide information about the opcode, including its identifier and operands.
type OpNotEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpNotEqual creates and returns a new instance of OpNotEqual with OpcodeDetails initialized from bytecode.
func NewOpNotEqual(op *bytecode.Opcodes) *OpNotEqual {
	return &OpNotEqual{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNotEqual)}
}

// Execute evaluates whether the top two stack elements are unequal and pushes the result as a boolean onto the stack.
func (op *OpNotEqual) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	right := v.Stack().Pop()
	left := v.Stack().Pop()
	val := op.Factory().FalseValue()
	if left.Equals(right) {
		val = op.Factory().TrueValue()
	}
	v.Stack().Push(val)
}
