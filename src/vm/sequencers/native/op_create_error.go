package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateError)
}

// OpCreateError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpCreateError struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCreateError creates and returns a new instance of OpCreateError with associated Opcode for the OpCreateError opcode.
func NewOpCreateError() handler.IOpExecutor {
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

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpCreateError) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute converts the top value on the Core stack into an error object and replaces it on the stack.
func (op *OpCreateError) Execute(_ *handler.Decoder) {
	value := op.vm.StackPeek()
	errObj := op.vm.CreateError(value)
	op.vm.StackSet(errObj)
}

// Compile generates the compiled representation of the OpCreateError operation or returns an unimplemented error.
func (op *OpCreateError) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
