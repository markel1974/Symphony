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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice() core.IOpExecutor {
	operands := _noOperands
	return &OpIndexSlice{
		opcode: opcodes.NewOpcode(OpIndexSliceId, operands, "OpIndexSlice"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIndexSlice) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIndexSlice) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(_ *core.Decoder) {
	highObj := op.vm.StackPop()
	lowObj := op.vm.StackPop()
	leftObj := op.vm.StackPop()
	lowIdx, highIdx, err := op.vm.Factory().BoundsCheck(lowObj, highObj, int64(leftObj.Length()))
	if err != nil {
		op.vm.SetError(err)
		return
	}
	switch left := leftObj.(type) {
	case *objects.Array:
		val := op.vm.Factory().NewArray(op.vm.FrameId(), left.Values()[lowIdx:highIdx])
		op.vm.StackPush(val)
	case *objects.String:
		val := op.vm.Factory().NewString(op.vm.FrameId(), left.Value()[lowIdx:highIdx])
		op.vm.StackPush(val)
	case *objects.Bytes:
		val := op.vm.Factory().NewBytes(op.vm.FrameId(), left.Value()[lowIdx:highIdx])
		op.vm.StackPush(val)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: %s[%d:%d]", left.TypeName(), lowIdx, highIdx))
		return
	}
}

// Compile generates the compiled representation of the OpIndexSlice operation or returns an unimplemented error.
func (op *OpIndexSlice) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
