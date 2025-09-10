package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpSuspend)
}

// OpSuspend represents an operation that suspends the execution of the virtual machine.
type OpSuspend struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpSuspend creates and returns a new OpSuspend instance with opcode details initialized for the suspend operation.
func NewOpSuspend() core.IOpExecutor {
	operands := _noOperands
	return &OpSuspend{
		opcode: opcodes.NewOpcode(OpSuspendId, operands, "OpSuspend"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpSuspend) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the suspend operation on the given virtual machine by setting its shutdown state to true.
func (op *OpSuspend) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.Shutdown()
}

// Opcode returns the opcode associated with the instance.
func (op *OpSuspend) Opcode() *opcodes.Opcode {
	return op.opcode
}
