package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpUnknown)
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding Opcode configuration set.
func NewOpUnknown(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpUnknown{
		Opcode: op.Opcode(opcodes.OpUnknown),
		vm:     vmT,
	}, nil
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.SetError(fmt.Errorf("unknown opcode at: %d", op.vm.GetIp()))
}
