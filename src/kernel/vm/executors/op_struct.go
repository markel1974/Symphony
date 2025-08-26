package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpStruct is a wrapper around bytecode.OpcodeDetails, representing a struct creation operation in bytecode execution.
type OpStruct struct {
	*bytecode.OpcodeDetails
}

// NewOpStruct initializes and returns a new instance of OpStruct with its OpcodeDetails set to OpMap details.
func NewOpStruct(op *bytecode.Opcodes) *OpStruct {
	return &OpStruct{OpcodeDetails: op.OpcodeToDetails(bytecode.OpStruct)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpStruct) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	mElem := v.Stack().PopMapElements(numElements)
	v.Stack().Push(op.Factory().NewStruct(v.FrameID(), mElem))
}
