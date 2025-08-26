package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIndexSlice)
}

// OpIndexSlice represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds Opcode to inherit opcode, operand, and name information for execution and identification.
type OpIndexSlice struct {
	*bytecode.Opcode
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIndexSlice{Opcode: op.Opcode(bytecode.OpIndexSlice)}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	highStack := v.Stack().Pop()
	lowStack := v.Stack().Pop()
	leftStack := v.Stack().Pop()
	lowIdx, highIdx, err := v.Factory().BoundsCheck(lowStack, highStack, int64(leftStack.Length()))
	if err != nil {
		v.SetError(err)
		return
	}
	var val objects.IObject = nil
	switch left := leftStack.(type) {
	case *objects.Array:
		val = v.Factory().NewArray(v.Frame().Id(), left.Values()[lowIdx:highIdx])
	case *objects.ArrayImmutable:
		val = v.Factory().NewArray(v.Frame().Id(), left.Values()[lowIdx:highIdx])
	case *objects.String:
		val = v.Factory().NewString(v.Frame().Id(), left.Value()[lowIdx:highIdx])
	case *objects.Bytes:
		val = v.Factory().NewBytes(v.Frame().Id(), left.Value()[lowIdx:highIdx])
	}
	if val != nil {
		v.Stack().Push(val)
	}
}
