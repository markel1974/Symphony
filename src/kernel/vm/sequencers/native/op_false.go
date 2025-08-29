package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFalse)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFalse{
		Opcode: op.Opcode(bytecode.OpFalse),
		vm:     vm,
	}
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(_ *core.Decoder) {
	// Operands Offset  0
	val := op.vm.Factory().FalseValue()
	op.vm.Stack().Push(val)
}
