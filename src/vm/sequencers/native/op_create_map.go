package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateMap)
}

// OpCreateMap is a wrapper around bytecode.Opcode, representing a map creation operation in bytecode execution.
type OpCreateMap struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateMap initializes and returns a new instance of OpCreateMap with its Opcode set to OpCreateMap details.
func NewOpCreateMap(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateMap{
		Opcode: op.Opcode(opcodes.OpCreateMap),
		vm:     vmT,
	}, nil
}

// Execute processes the OpCreateMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpCreateMap) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	mElem := op.vm.Stack().PopMapElements(numElements)
	op.vm.Stack().Push(op.vm.Factory().NewMap(op.vm.Frame().Id(), mElem))
}
