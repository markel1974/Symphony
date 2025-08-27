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
}

// NewOpDerefSet creates a new OpDerefSet instance with the specified opcode for the dereference and set operation.
func NewOpDerefSet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpDerefSet{Opcode: op.Opcode(bytecode.OpDerefSet)}
}

// Execute performs a dereference-and-set operation on the stack, assigning a value to the object pointed by a pointer.
func (op *OpDerefSet) Execute(v *core.VM, _ *core.Decoder) {
	// Stack (dall'alto verso il basso): [valore, puntatore, ...]
	valueToSet := v.Stack().Pop()
	pointerObj := v.Stack().Pop()

	ptr, ok := pointerObj.(*objects.ObjectPointer)
	if !ok {
		v.SetError(fmt.Errorf("invalid operation: cannot assign to a non-pointer type %s", pointerObj.TypeName()))
		return
	}

	// Assegna il nuovo valore all'oggetto puntato.
	ptr.SetValue(valueToSet)

	// Lascia il valore assegnato sullo stack per possibili assegnazioni a catena,
	// o può essere rimosso con un OpPop successivo se necessario.
	v.Stack().Push(valueToSet)
}
