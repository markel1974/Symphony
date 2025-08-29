package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpMap)
}

// OpMap is a wrapper around bytecode.Opcode, representing a map creation operation in bytecode execution.
type OpMap struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpMap initializes and returns a new instance of OpMap with its Opcode set to OpMap details.
func NewOpMap(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpMap{
		Opcode: op.Opcode(bytecode.OpMap),
		vm:     vm,
	}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpMap) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	mElem := op.vm.Stack().PopMapElements(numElements)
	op.vm.Stack().Push(op.vm.Factory().NewMap(op.vm.Frame().Id(), mElem))
}
