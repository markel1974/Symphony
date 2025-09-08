package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpStruct)
}

// OpStruct is a wrapper around bytecode.Opcode, representing a struct creation operation in bytecode execution.
type OpStruct struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpStruct initializes and returns a new instance of OpStruct with its Opcode set to OpMap details.
func NewOpStruct(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpStruct{
		Opcode: op.Opcode(bytecode.OpStruct),
		vm:     vmT,
	}, nil
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpStruct) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	typeNameObj := op.vm.Stack().Pop()
	mElem := op.vm.Stack().PopMapElements(numElements)
	structObj := op.vm.Factory().NewStruct(op.vm.Frame().Id(), typeNameObj.AsString(), mElem)
	op.vm.Stack().Push(structObj)
}
