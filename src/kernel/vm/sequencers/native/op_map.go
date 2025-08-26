package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpMap)
}

// OpMap is a wrapper around bytecode.OpcodeDetails, representing a map creation operation in bytecode execution.
type OpMap struct {
	*bytecode.OpcodeDetails
}

// NewOpMap initializes and returns a new instance of OpMap with its OpcodeDetails set to OpMap details.
func NewOpMap(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpMap{OpcodeDetails: op.OpcodeToDetails(bytecode.OpMap)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpMap) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	mElem := v.Stack().PopMapElements(numElements)
	v.Stack().Push(op.Factory().NewMap(v.Frame().Id(), mElem))
}
