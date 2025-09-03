package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpSuspend)
}

// OpSuspend represents an operation that suspends the execution of the virtual machine.
type OpSuspend struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpSuspend creates and returns a new OpSuspend instance with opcode details initialized for the suspend operation.
func NewOpSuspend(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpSuspend{
		Opcode: op.Opcode(bytecode.OpSuspend),
		vm:     vmT,
	}, nil
}

// Execute performs the suspend operation on the given virtual machine by setting its shutdown state to true.
func (op *OpSuspend) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.Shutdown()
}
