package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpCreateStruct)
}

// OpCreateStruct is a wrapper around bytecode.Opcode, representing a struct creation operation in bytecode execution.
type OpCreateStruct struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateStruct initializes and returns a new instance of OpCreateStruct with its Opcode set to OpCreateMap details.
func NewOpCreateStruct(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateStruct{
		Opcode: op.Opcode(bytecode.OpCreateStruct),
		vm:     vmT,
	}, nil
}

// Execute processes the OpCreateMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpCreateStruct) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	typeNameObj := op.vm.Stack().Pop()
	mElem := op.vm.Stack().PopMapElements(numElements)
	structObj := op.vm.Factory().NewStruct(op.vm.Frame().Id(), typeNameObj.AsString(), mElem)
	op.vm.Stack().Push(structObj)
}
