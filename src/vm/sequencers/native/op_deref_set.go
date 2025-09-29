package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpDerefSet operation with the sequencer system.
func init() {
	SequencerRegister(NewOpDerefSet)
}

// OpDerefSet represents an operation for dereferencing a pointer and setting its value in the virtual machine.
type OpDerefSet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpDerefSet creates a new OpDerefSet instance with the specified opcode for the dereference and set operation.
func NewOpDerefSet() handler.IOpExecutor {
	operands := _noOperands
	return &OpDerefSet{
		opcode: opcodes.NewOpcode(OpDerefSetId, operands, "OpDerefSet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpDerefSet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Opcode returns the opcode associated with the instance.
func (op *OpDerefSet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Execute performs a dereference-and-set operation on the stack, assigning a value to the object pointed by a pointer.
func (op *OpDerefSet) Execute(_ *handler.Decoder) {
	pointerObj := op.vm.StackPop()
	valueToSet := op.vm.StackPop()
	ptr, ok := pointerObj.(*objects.ObjectPointer)
	if !ok {
		op.vm.Shutdown(fmt.Errorf("invalid operation: cannot assign to a non-pointer type %s", pointerObj.TypeName()))
		return
	}
	if err := ptr.AssignValue(valueToSet); err != nil {
		op.vm.Shutdown(err)
		return
	}
	op.vm.StackPush(valueToSet)
}

// Compile generates the compiled representation of the OpDerefSet operation or returns an unimplemented error.
func (op *OpDerefSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
