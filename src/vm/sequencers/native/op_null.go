package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpNull)
}

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull() core.IOpExecutor {
	operands := _noOperands
	return &OpNull{
		opcode: opcodes.NewOpcode(OpNullId, operands, "OpNull"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpNull) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpNull) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(_ *core.Decoder) {
	op.vm.StackPush(op.vm.Factory().UndefinedValue())
}

// Compile generates the compiled representation of the OpNull operation or returns an unimplemented error.
func (op *OpNull) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
