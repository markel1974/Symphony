package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateError)
}

// OpCreateError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpCreateError struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCreateError creates and returns a new instance of OpCreateError with associated Opcode for the OpCreateError opcode.
func NewOpCreateError() core.IOpExecutor {
	operands := _noOperands
	return &OpCreateError{
		opcode: opcodes.NewOpcode(OpCreateErrorId, operands, "OpCreateError"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateError) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateError) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpCreateError) Execute(_ *core.Decoder) {
	value := op.vm.StackPeek()
	e := op.vm.Factory().NewError(op.vm.FrameId(), value.AsString())
	op.vm.StackSet(e)
}

// Compile generates the compiled representation of the OpCreateError operation or returns an unimplemented error.
func (op *OpCreateError) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
