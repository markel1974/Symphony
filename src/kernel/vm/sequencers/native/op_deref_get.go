package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpDerefGet)
}

// OpDerefGet represents an operation for dereferencing a pointer.
type OpDerefGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpDerefGet creates a new OpDerefGet instance.
func NewOpDerefGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpDerefGet{
		Opcode: op.Opcode(bytecode.OpDerefGet),
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
	// Replace the stack pointer with the value it points to.
	op.vm.Stack().Push(*ptr.Value())
}
