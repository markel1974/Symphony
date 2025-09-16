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
	targetObj := op.vm.StackPop()
	ret, err := op.vm.Factory().CreateSlice(op.vm.FrameId(), lowObj, highObj, targetObj)
	if err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.StackPush(ret)
}

// Compile generates the compiled representation of the OpIndexSlice operation or returns an unimplemented error.
func (op *OpIndexSlice) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
