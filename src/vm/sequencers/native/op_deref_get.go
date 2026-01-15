package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpDerefGet)
}

// OpDerefGet represents an operation for dereferencing a pointer.
type OpDerefGet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpDerefGet creates a new OpDerefGet instance.
func NewOpDerefGet() handler.IOpExecutor {
	operands := _noOperands
	return &OpDerefGet{
		opcode: opcodes.NewOpcode(OpDerefGetId, operands, "OpDerefGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpDerefGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpDerefGet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the dereference operation. It takes a pointer from the
// stack and replaces it with the value it points to.
func (op *OpDerefGet) Execute(_ *handler.Decoder) {
	operand := op.vm.StackPop()
	ptr, ok := operand.(*objects.ObjectPointer)
	if !ok {
		op.vm.Shutdown(fmt.Errorf("invalid operation: cannot dereference non-pointer type %s", operand.TypeName()))
		return
	}
	op.vm.StackPush(*ptr.Value())
}

// Compile generates the compiled representation of the OpDerefGet operation or returns an unimplemented error.
func (op *OpDerefGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
