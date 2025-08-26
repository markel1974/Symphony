package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpArray)
}

// OpArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base OpcodeDetails for opcode, operands, and name information.
type OpArray struct {
	*bytecode.OpcodeDetails
}

// NewOpArray creates and returns a new instance of OpArray, initialized with details for the OpArray operation.
func NewOpArray(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpArray{OpcodeDetails: op.OpcodeToDetails(bytecode.OpArray)}
}

// Execute processes the OpArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpArray) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Read(0)
	elements := v.Stack().PopArrayElements(numElements)
	arr := op.Factory().NewArray(v.FrameID(), elements)
	v.Stack().Push(arr)
}
