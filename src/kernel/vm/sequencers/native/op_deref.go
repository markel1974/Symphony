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
	*bytecode.Opcode
	vm *core.VM
}

// NewOpDeref creates a new OpDeref instance.
func NewOpDeref(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpDeref{
		Opcode: op.Opcode(bytecode.OpDeref),
		vm:     vm,
	}
}

// Execute performs the dereference operation. It takes a pointer from the
// stack and replaces it with the value it points to.
func (op *OpDeref) Execute(_ *core.Decoder) {
	operand := op.vm.Stack().Pop()
	ptr, ok := operand.(*objects.ObjectPointer)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid operation: cannot dereference non-pointer type %s", operand.TypeName()))
		return
	}
	// Replace the stack pointer with the value it points to.
	op.vm.Stack().Push(*ptr.Value())
}
