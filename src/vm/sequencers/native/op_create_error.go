package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateError)
}

// OpCreateError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpCreateError struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateError creates and returns a new instance of OpCreateError with associated Opcode for the OpCreateError opcode.
func NewOpCreateError(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateError{
		Opcode: op.Opcode(opcodes.OpCreateError),
		vm:     vmT,
	}, nil
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpCreateError) Execute(_ *core.Decoder) {
	// Operands Offset  0
	value := op.vm.Stack().Peek()
	e := op.vm.Factory().NewError(op.vm.Frame().Id(), value.AsString())
	op.vm.Stack().Set(e)
}
