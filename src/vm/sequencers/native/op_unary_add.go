package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the unary addition operation with the sequencer system by appending NewOpUnaryAdd to the internal registry.
func init() {
	SequencerRegister(NewOpUnaryAdd)
}

// OpUnaryAdd represents an opcode that performs a unary addition operation in the virtual machine.
// It embeds the bytecode.Opcode type for opcode execution and uses core.IVMFullAccess for VM interaction.
type OpUnaryAdd struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpUnaryAdd initializes and returns an OpUnaryAdd executor, ensuring the provided VM supports full-access operations.
func NewOpUnaryAdd() core.IOpExecutor {
	operands := _noOperands
	return &OpUnaryAdd{
		opcode: opcodes.NewOpcode(OpUnaryAddId, operands, "OpUnaryAdd"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpUnaryAdd) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute executes the unary addition operation on the top operand of the stack and pushes the result back onto the stack.
func (op *OpUnaryAdd) Execute(_ *core.Decoder) {
	// Operands Offset 0
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

// Opcode returns the opcode associated with the instance.
func (op *OpUnaryAdd) Opcode() *opcodes.Opcode {
	return op.opcode
}
