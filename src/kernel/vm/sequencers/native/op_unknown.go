package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpUnknown)
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding Opcode configuration set.
func NewOpUnknown(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpUnknown{
		Opcode: op.Opcode(bytecode.OpUnknown),
		vm:     vm,
	}
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.SetError(fmt.Errorf("unknown opcode at: %d", op.vm.GetIp()))
}
