package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorNext)
}

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorNext{
		opcode: opcodes.NewOpcode(OpIteratorNextId, operands, "OpIteratorNext"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorNext) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIteratorNext) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(decoder *handler.Decoder) {
	stableIdx := decoder.Operand(0)
	itObj := op.vm.StackPeekBP(uint(stableIdx))
	it, ok := itObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(objects.ComputeIteratorError(objects.ErrNotIterator, itObj.TypeName()))
		return
	}
	if it.Next() {
		op.vm.StackPush(op.vm.Factory().TrueValue())
	} else {
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Compile generates the compiled representation of the OpIteratorNext operation or returns an unimplemented error.
func (op *OpIteratorNext) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
