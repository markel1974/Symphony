package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers the NewOpDerefSet operation with the sequencer system.
func init() {
	SequencerRegister(NewOpDerefSet)
}

// OpDerefSet represents an operation for dereferencing a pointer and setting its value in the virtual machine.
type OpDerefSet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpDerefSet creates a new OpDerefSet instance with the specified opcode for the dereference and set operation.
func NewOpDerefSet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpDerefSet{
		Opcode: op.Opcode(bytecode.OpDerefSet),
		vm:     vmT,
	}, nil
}

// Execute performs a dereference-and-set operation on the stack, assigning a value to the object pointed by a pointer.
func (op *OpDerefSet) Execute(_ *core.Decoder) {
	// Stack (dall'alto verso il basso): [valore, puntatore, ...]
	valueToSet := op.vm.Stack().Pop()
	pointerObj := op.vm.Stack().Pop()

	ptr, ok := pointerObj.(*objects.ObjectPointer)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid operation: cannot assign to a non-pointer type %s", pointerObj.TypeName()))
		return
	}

	// Assegna il nuovo valore all'oggetto puntato.
	ptr.SetValue(valueToSet)

	// Lascia il valore assegnato sullo stack per possibili assegnazioni a catena,
	// o può essere rimosso con un OpPop successivo se necessario.
	op.vm.Stack().Push(valueToSet)
}
