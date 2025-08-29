package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpEqual)
}

// OpEqual represents an operation that checks if two values are equal and updates the stack accordingly.
type OpEqual struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpEqual creates and returns an instance of OpEqual, initialized with its corresponding opcode details.
func NewOpEqual(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpEqual{
		Opcode: op.Opcode(bytecode.OpEqual),
		vm:     vm,
	}
}

// Execute performs the equality comparison between the top two stack values and pushes the result (true or false) back onto the stack.
func (op *OpEqual) Execute(_ *core.Decoder) {
	// Operands Offset 0
	right := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	var val objects.IObject

	if left.Equals(right) {
		val = op.vm.Factory().TrueValue()
		//val = v.Factory().FalseValue()
	} else {
		val = op.vm.Factory().FalseValue()
		//val = v.Factory().TrueValue()
	}
	op.vm.Stack().Push(val)
}
