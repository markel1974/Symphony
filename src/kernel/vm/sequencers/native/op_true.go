package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpTrue)
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpTrue{
		Opcode: op.Opcode(bytecode.OpTrue),
		vm:     vm,
	}
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(_ *core.Decoder) {
	// Operands Offset 0
	val := op.vm.Factory().TrueValue()
	op.vm.Stack().Push(val)
}
