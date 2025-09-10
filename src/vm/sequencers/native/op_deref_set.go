package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpDerefSet operation with the sequencer system.
func init() {
	SequencerRegister(NewOpDerefSet)
}

// OpDerefSet represents an operation for dereferencing a pointer and setting its value in the virtual machine.
type OpDerefSet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpDerefSet creates a new OpDerefSet instance with the specified opcode for the dereference and set operation.
func NewOpDerefSet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpDerefSet{
		Opcode: op.Opcode(opcodes.OpDerefSet),
		vm:     vmT,
	}, nil
}

// Execute performs a dereference-and-set operation on the stack, assigning a value to the object pointed by a pointer.
func (op *OpDerefSet) Execute(_ *core.Decoder) {
	pointerObj := op.vm.Stack().Pop()
	valueToSet := op.vm.Stack().Pop()
	ptr, ok := pointerObj.(*objects.ObjectPointer)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid operation: cannot assign to a non-pointer type %s", pointerObj.TypeName()))
		return
	}
	if err := ptr.AssignValue(valueToSet); err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.Stack().Push(valueToSet)
}
