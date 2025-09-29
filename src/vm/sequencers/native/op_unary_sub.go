package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
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
	vm     handler.IVMFullAccess
}

// NewOpUnarySub creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpUnarySub() handler.IOpExecutor {
	operands := _noOperands
	return &OpUnarySub{
		opcode: opcodes.NewOpcode(OpUnarySubId, operands, "OpUnarySub"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnarySub) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpUnarySub) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpUnarySub) Execute(_ *handler.Decoder) {
	obj := op.vm.StackPop()
	switch x := obj.(type) {
	case *objects.Float:
		res := op.vm.Factory().NewFloat(op.vm.FrameId(), -x.Value())
		op.vm.StackPush(res)
	default:
		res := op.vm.Factory().NewInt(op.vm.FrameId(), -x.AsInt64())
		op.vm.StackPush(res)
	}
}

// Compile generates the compiled representation of the OpUnarySub operation or returns an unimplemented error.
func (op *OpUnarySub) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
