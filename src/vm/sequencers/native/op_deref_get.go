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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpDerefGet creates a new OpDerefGet instance.
func NewOpDerefGet() core.IOpExecutor {
	operands := _noOperands
	return &OpDerefGet{
		opcode: opcodes.NewOpcode(OpDerefGetId, operands, "OpDerefGet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpDerefGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
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

// Opcode returns the opcode associated with the instance.
func (op *OpDerefGet) Opcode() *opcodes.Opcode {
	return op.opcode
}
