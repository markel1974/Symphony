package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpIndexSlice)
}

// OpIndexSlice represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds Opcode to inherit opcode, operand, and name information for execution and identification.
type OpIndexSlice struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIndexSlice{
		Opcode: op.Opcode(bytecode.OpIndexSlice),
		vm:     vmT,
	}, nil
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
	var val objects2.IObject = nil
	switch left := leftStack.(type) {
	case *objects2.Array:
		val = op.vm.Factory().NewArray(op.vm.Frame().Id(), left.Values()[lowIdx:highIdx])
	case *objects2.String:
		val = op.vm.Factory().NewString(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
	case *objects2.Bytes:
		val = op.vm.Factory().NewBytes(op.vm.Frame().Id(), left.Value()[lowIdx:highIdx])
	}
	if val != nil {
		op.vm.Stack().Push(val)
	}
}
