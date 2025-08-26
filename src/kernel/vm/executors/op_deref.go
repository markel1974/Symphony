package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpDeref represents an operation for dereferencing a pointer.
type OpDeref struct {
	*bytecode.OpcodeDetails
}

// NewOpDeref creates a new OpDeref instance.
func NewOpDeref(op *bytecode.Opcodes) *OpDeref {
	return &OpDeref{OpcodeDetails: op.OpcodeToDetails(bytecode.OpDeref)}
}

// Execute performs the dereference operation. It takes a pointer from the
// stack and replaces it with the value it points to.
func (op *OpDeref) Execute(v *core.VM, _ *core.Decoder) {
	operand := v.Stack().Pop()

	ptr, ok := operand.(*objects.ObjectPointer)
	if !ok {
		v.SetError(fmt.Errorf("invalid operation: cannot dereference non-pointer type %s", operand.TypeName()))
		return
	}

	// Sostituisce il puntatore sullo stack con il valore puntato.
	v.Stack().Push(*ptr.Value())
}
