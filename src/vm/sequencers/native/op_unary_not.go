package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpUnaryNot)
}

// OpUnaryNot represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpUnaryNot struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpUnaryNot creates a new instance of OpUnaryNot, representing a logical NOT operation (!).
func NewOpUnaryNot() core.IOpExecutor {
	operands := _noOperands
	return &OpUnaryNot{
		opcode: opcodes.NewOpcode(OpUnaryNotId, operands, "OpUnaryNot"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpUnaryNot) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpUnaryNot) Execute(_ *core.Decoder) {
	// Operands Offset  0
	operand := op.vm.Stack().Pop()
	var val objects.IObject
	if operand.Falsy() {
		val = op.vm.Factory().TrueValue()
	} else {
		val = op.vm.Factory().FalseValue()
	}
	op.vm.Stack().Push(val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnaryNot) Opcode() *opcodes.Opcode {
	return op.opcode
}
