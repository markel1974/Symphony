package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpDerefGet)
}

// OpDerefGet represents an operation for dereferencing a pointer.
type OpDerefGet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpDerefGet creates a new OpDerefGet instance.
func NewOpDerefGet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpDerefGet{
		Opcode: op.Opcode(opcodes.OpDerefGet),
		vm:     vmT,
	}, nil
}

// Execute performs the dereference operation. It takes a pointer from the
// stack and replaces it with the value it points to.
func (op *OpDerefGet) Execute(_ *core.Decoder) {
	operand := op.vm.Stack().Pop()
	ptr, ok := operand.(*objects.ObjectPointer)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid operation: cannot dereference non-pointer type %s", operand.TypeName()))
		return
	}
	op.vm.Stack().Push(*ptr.Value())
}
