package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpArray)
}

// OpArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base Opcode for opcode, operands, and name information.
type OpArray struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpArray creates and returns a new instance of OpArray, initialized with details for the OpArray operation.
func NewOpArray(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmFullAccess, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpArray{
		Opcode: op.Opcode(bytecode.OpArray),
		vm:     vmFullAccess,
	}, nil
}

// Execute processes the OpArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpArray) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	elements := op.vm.Stack().PopArrayElements(numElements)
	arr := op.vm.Factory().NewArray(op.vm.Frame().Id(), elements)
	op.vm.Stack().Push(arr)
}
