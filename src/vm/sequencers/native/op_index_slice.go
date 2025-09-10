package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	objects "github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIndexSlice)
}

// OpIndexSlice represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds Opcode to inherit opcode, operand, and name information for execution and identification.
type OpIndexSlice struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIndexSlice{
		Opcode: op.Opcode(opcodes.OpIndexSlice),
		vm:     vmT,
	}, nil
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(_ *core.Decoder) {
	// Operands Offset  0
	highStack := op.vm.Stack().Pop()
	lowStack := op.vm.Stack().Pop()
	leftObj := op.vm.Stack().Pop()
	lowIdx, highIdx, err := op.vm.Factory().BoundsCheck(lowStack, highStack, int64(leftObj.Length()))
	if err != nil {
		op.vm.SetError(err)
		return
	}
	switch left := leftObj.(type) {
	case *objects.Array:
		val := op.vm.Factory().NewArray(op.vm.Frame().Id(), left.Values()[lowIdx:highIdx])
		op.vm.Stack().Push(val)
	case *objects.String:
		val := op.vm.Factory().NewString(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
		op.vm.Stack().Push(val)
	case *objects.Bytes:
		val := op.vm.Factory().NewBytes(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
		op.vm.Stack().Push(val)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: %s[%d:%d]", left.TypeName(), lowIdx, highIdx))
		return
	}
}
