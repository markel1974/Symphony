package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpSliceIndex represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds OpcodeDetails to inherit opcode, operand, and name information for execution and identification.
type OpSliceIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpSliceIndex creates a new instance of OpSliceIndex containing details for the slice indexing bytecode operation.
func NewOpSliceIndex(op *bytecode.Opcodes) *OpSliceIndex {
	return &OpSliceIndex{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSliceIndex)}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpSliceIndex) Execute(v *core.VM, _ *core.Decoder) {
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
