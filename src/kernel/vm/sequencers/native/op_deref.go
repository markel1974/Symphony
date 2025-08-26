package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpDeref)
}

// OpDeref represents an operation for dereferencing a pointer.
type OpDeref struct {
	*bytecode.OpcodeDetails
}

// NewOpDeref creates a new OpDeref instance.
func NewOpDeref(op *bytecode.Opcodes) core.IOpExecutor {
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
	// Replace the stack pointer with the value it points to.
	v.Stack().Push(*ptr.Value())
}
