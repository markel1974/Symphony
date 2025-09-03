package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpError)
}

// OpError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpError struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpError creates and returns a new instance of OpError with associated Opcode for the OpError opcode.
func NewOpError(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpError{
		Opcode: op.Opcode(bytecode.OpError),
		vm:     vmT,
	}, nil
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpError) Execute(_ *core.Decoder) {
	// Operands Offset  0
	value := op.vm.Stack().Peek()
	e := op.vm.Factory().NewError(op.vm.Frame().Id(), value.AsString())
	op.vm.Stack().Set(e)
}
