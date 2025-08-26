package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIndexSlice)
}

// OpIndexSlice represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds OpcodeDetails to inherit opcode, operand, and name information for execution and identification.
type OpIndexSlice struct {
	*bytecode.OpcodeDetails
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIndexSlice{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIndexSlice)}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	highStack := v.Stack().Pop()
	lowStack := v.Stack().Pop()
	leftStack := v.Stack().Pop()
	lowIdx, highIdx, err := op.Factory().BoundsCheck(lowStack, highStack, int64(leftStack.Length()))
	if err != nil {
		v.SetError(err)
		return
	}
	var val objects.IObject = nil
	switch left := leftStack.(type) {
	case *objects.Array:
		val = op.Factory().NewArray(v.FrameID(), left.Values()[lowIdx:highIdx])
	case *objects.ArrayImmutable:
		val = op.Factory().NewArray(v.FrameID(), left.Values()[lowIdx:highIdx])
	case *objects.String:
		val = op.Factory().NewString(v.FrameID(), left.Value()[lowIdx:highIdx])
	case *objects.Bytes:
		val = op.Factory().NewBytes(v.FrameID(), left.Value()[lowIdx:highIdx])
	}
	if val != nil {
		v.Stack().Push(val)
	}
}
