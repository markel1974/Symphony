package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpUnarySub)
}

// OpUnarySub represents an operation for negating a numeric value.
// It embeds Opcode, providing details such as the opcode, operands, and name.
type OpUnarySub struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpUnarySub creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpUnarySub() core.IOpExecutor {
	operands := _noOperands
	return &OpUnarySub{
		opcode: opcodes.NewOpcode(OpUnarySubId, operands, "OpUnarySub"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpUnarySub) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpUnarySub) Execute(_ *core.Decoder) {
	// Operands Offset 0
	operand := op.vm.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Float:
		res := op.vm.Factory().NewFloat(op.vm.Frame().Id(), -x.Value())
		op.vm.Stack().Push(res)
	default:
		res := op.vm.Factory().NewInt(op.vm.Frame().Id(), -x.AsInt64())
		op.vm.Stack().Push(res)
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnarySub) Opcode() *opcodes.Opcode {
	return op.opcode
}
