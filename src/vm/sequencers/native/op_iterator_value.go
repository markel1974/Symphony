package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorValue)
}

// OpIteratorValue retrieves the value from the current iterator position.
// It embeds Opcode, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated Opcode initialized.
func NewOpIteratorValue() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorValue{
		opcode: opcodes.NewOpcode(OpIteratorValueId, operands, "OpIteratorValue"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorValue) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIteratorValue) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(decoder *handler.Decoder) {
	stableIdx := decoder.Operand(0)
	itObj := op.vm.StackPeekBP(uint(stableIdx))
	it, ok := itObj.(objects.IIterator)
	if !ok {
		op.vm.Shutdown(objects.ComputeIteratorError(objects.ErrNotIterator, itObj.TypeName()))
		return
	}
	op.vm.StackPush(it.Value(op.vm.FrameId()))
}

// Compile generates the compiled representation of the OpIteratorValue operation or returns an unimplemented error.
func (op *OpIteratorValue) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
