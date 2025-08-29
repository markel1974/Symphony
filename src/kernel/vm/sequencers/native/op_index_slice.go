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
	vm *core.VM
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIndexSlice{
		Opcode: op.Opcode(bytecode.OpIndexSlice),
		vm:     vm,
	}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(_ *core.Decoder) {
	// Operands Offset  0
	highStack := op.vm.Stack().Pop()
	lowStack := op.vm.Stack().Pop()
	leftStack := op.vm.Stack().Pop()
	lowIdx, highIdx, err := op.vm.Factory().BoundsCheck(lowStack, highStack, int64(leftStack.Length()))
	if err != nil {
		op.vm.SetError(err)
		return
	}
	var val objects.IObject = nil
	switch left := leftStack.(type) {
	case *objects.Array:
		val = op.vm.Factory().NewArray(op.vm.Frame().Id(), left.Values()[lowIdx:highIdx])
	case *objects.ArrayImmutable:
		val = op.vm.Factory().NewArray(op.vm.Frame().Id(), left.Values()[lowIdx:highIdx])
	case *objects.String:
		val = op.vm.Factory().NewString(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
	case *objects.Bytes:
		val = op.vm.Factory().NewBytes(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
	}
	if val != nil {
		op.vm.Stack().Push(val)
	}
}
