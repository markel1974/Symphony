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
	vm *core.VM
}

// NewOpNotEqual creates and returns a new instance of OpNotEqual with Opcode initialized from bytecode.
func NewOpNotEqual(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNotEqual{
		Opcode: op.Opcode(bytecode.OpNotEqual),
		vm:     vm,
	}
}

// Execute evaluates whether the top two stack elements are unequal and pushes the result as a boolean onto the stack.
func (op *OpNotEqual) Execute(_ *core.Decoder) {
	// Operands Offset  0
	right := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	var val objects.IObject
	if left.Equals(right) {
		val = op.vm.Factory().FalseValue()
		//val = v.Factory().TrueValue()
	} else {
		val = op.vm.Factory().TrueValue()
		//val = v.Factory().FalseValue()
	}
	op.vm.Stack().Push(val)
}
