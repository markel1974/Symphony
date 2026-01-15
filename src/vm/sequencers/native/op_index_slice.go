package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIndexSlice)
}

// OpIndexSlice represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds Opcode to inherit opcode, operand, and name information for execution and identification.
type OpIndexSlice struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIndexSlice creates a new instance of OpIndexSlice containing details for the slice indexing bytecode operation.
func NewOpIndexSlice() handler.IOpExecutor {
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

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIndexSlice) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpIndexSlice) Execute(_ *handler.Decoder) {
	highObj := op.vm.StackPop()
	lowObj := op.vm.StackPop()
	targetObj := op.vm.StackPop()
	var high int
	var low int
	if _, ok := highObj.(*objects.Undefined); ok {
		high = targetObj.Length()
	} else {
		high = int(highObj.AsInt64())
	}
	if _, ok := lowObj.(*objects.Undefined); ok {
		low = 0
	} else {
		low = int(lowObj.AsInt64())
	}
	slice := op.vm.CreateSlice(high, low, targetObj)
	op.vm.StackPush(slice)
}

// Compile generates the compiled representation of the OpIndexSlice operation or returns an unimplemented error.
func (op *OpIndexSlice) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
