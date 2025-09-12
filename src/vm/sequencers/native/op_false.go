package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFalse)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse() core.IOpExecutor {
	operands := _noOperands
	return &OpFalse{
		opcode: opcodes.NewOpcode(OpFalseId, operands, "OpFalse"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpFalse) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(_ *core.Decoder) {
	// Operands Offset  0
	val := op.vm.Factory().FalseValue()
	op.vm.StackPush(val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpFalse) Opcode() *opcodes.Opcode {
	return op.opcode
}
