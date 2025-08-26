package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpStruct)
}

// OpStruct is a wrapper around bytecode.Opcode, representing a struct creation operation in bytecode execution.
type OpStruct struct {
	*bytecode.Opcode
}

// NewOpStruct initializes and returns a new instance of OpStruct with its Opcode set to OpMap details.
func NewOpStruct(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpStruct{Opcode: op.Opcode(bytecode.OpStruct)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpStruct) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	mElem := v.Stack().PopMapElements(numElements)
	v.Stack().Push(v.Factory().NewStruct(v.Frame().Id(), mElem))
}
