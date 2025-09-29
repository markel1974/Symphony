package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the unary addition operation with the sequencer system by appending NewOpUnaryAdd to the internal registry.
func init() {
	SequencerRegister(NewOpUnaryAdd)
}

// OpUnaryAdd represents an opcode that performs a unary addition operation in the virtual machine.
// It embeds the bytecode.Opcode type for opcode execution and uses handler.IVMFullAccess for Core interaction.
type OpUnaryAdd struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpUnaryAdd initializes and returns an OpUnaryAdd executor, ensuring the provided Core supports full-access operations.
func NewOpUnaryAdd() handler.IOpExecutor {
	operands := _noOperands
	return &OpUnaryAdd{
		opcode: opcodes.NewOpcode(OpUnaryAddId, operands, "OpUnaryAdd"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnaryAdd) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpUnaryAdd) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute executes the unary addition operation on the top operand of the stack and pushes the result back onto the stack.
func (op *OpUnaryAdd) Execute(_ *handler.Decoder) {
	operand := op.vm.StackPop()
	switch x := operand.(type) {
	case *objects.Float:
		res := op.vm.Factory().NewFloat(op.vm.FrameId(), +x.Value())
		op.vm.StackPush(res)
	default:
		res := op.vm.Factory().NewInt(op.vm.FrameId(), +x.AsInt64())
		op.vm.StackPush(res)
	}
}

// Compile generates the compiled representation of the OpUnaryAdd operation or returns an unimplemented error.
func (op *OpUnaryAdd) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
