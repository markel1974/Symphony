package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the unary addition operation with the sequencer system by appending NewOpUnaryAdd to the internal registry.
func init() {
	SequencerRegister(NewOpUnaryAdd)
}

// OpUnaryAdd represents an opcode that performs a unary addition operation in the virtual machine.
// It embeds the bytecode.Opcode type for opcode execution and uses core.IVMFullAccess for VM interaction.
type OpUnaryAdd struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpUnaryAdd initializes and returns an OpUnaryAdd executor, ensuring the provided VM supports full-access operations.
func NewOpUnaryAdd(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpUnaryAdd{
		Opcode: op.Opcode(bytecode.OpUnaryAdd),
		vm:     vmT,
	}, nil
}

// Execute executes the unary addition operation on the top operand of the stack and pushes the result back onto the stack.
func (op *OpUnaryAdd) Execute(_ *core.Decoder) {
	// Operands Offset 0
	operand := op.vm.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Float:
		res := op.vm.Factory().NewFloat(op.vm.Frame().Id(), +x.Value())
		op.vm.Stack().Push(res)
	default:
		res := op.vm.Factory().NewInt(op.vm.Frame().Id(), +x.AsInt64())
		op.vm.Stack().Push(res)
	}
}
