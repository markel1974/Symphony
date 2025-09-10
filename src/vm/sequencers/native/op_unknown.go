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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding Opcode configuration set.
func NewOpUnknown() core.IOpExecutor {
	operands := _noOperands
	return &OpUnknown{
		opcode: opcodes.NewOpcode(OpUnknownId, operands, "OpUnknown"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpUnknown) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.SetError(fmt.Errorf("unknown opcode at: %d", op.vm.GetIp()))
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnknown) Opcode() *opcodes.Opcode {
	return op.opcode
}
