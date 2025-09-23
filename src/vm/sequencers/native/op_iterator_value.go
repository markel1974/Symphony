package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorValue)
}

// OpIteratorValue retrieves the value from the current iterator position.
// It embeds Opcode, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated Opcode initialized.
func NewOpIteratorValue() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIteratorValue) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(decoder *core.Decoder) {
	stableIdx := decoder.Operand(0)
	itObj := op.vm.StackPeekBP(uint(stableIdx))
	it, ok := itObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(objects.ComputeIteratorError(objects.ErrNotIterator, itObj.TypeName()))
		return
	}
	op.vm.StackPush(it.Value(op.vm.FrameId()))
}

// Compile generates the compiled representation of the OpIteratorValue operation or returns an unimplemented error.
func (op *OpIteratorValue) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
