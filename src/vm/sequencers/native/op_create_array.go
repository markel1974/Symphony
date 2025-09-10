package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateArray)
}

// OpCreateArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base Opcode for opcode, operands, and name information.
type OpCreateArray struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateArray creates and returns a new instance of OpCreateArray, initialized with details for the OpCreateArray operation.
func NewOpCreateArray(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmFullAccess, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateArray{
		Opcode: op.Opcode(opcodes.OpCreateArray),
		vm:     vmFullAccess,
	}, nil
}

// Execute processes the OpCreateArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpCreateArray) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	elements := op.vm.Stack().PopArrayElements(numElements)
	arr := op.vm.Factory().NewArray(op.vm.Frame().Id(), elements)
	op.vm.Stack().Push(arr)
}
