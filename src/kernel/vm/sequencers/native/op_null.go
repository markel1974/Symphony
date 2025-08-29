package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpNull)
}

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNull{
		Opcode: op.Opcode(bytecode.OpNull),
		vm:     vm,
	}
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(_ *core.Decoder) {
	// Operands Offset 0
	val := op.vm.Factory().UndefinedValue()
	op.vm.Stack().Push(val)
}
