package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorKey)
}

// OpIteratorKey wraps bytecode.Opcode to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorKey{
		opcode: opcodes.NewOpcode(OpIteratorKeyId, operands, "OpIteratorKey"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorKey) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIteratorKey) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the "iterator key" operation, retrieves the iterator key, and pushes it onto the Core stack.
func (op *OpIteratorKey) Execute(decoder *handler.Decoder) {
	stableIdx := decoder.Operand(0)
	itObj := op.vm.StackPeekBP(uint(stableIdx))
	it, ok := itObj.(objects.IIterator)
	if !ok {
		op.vm.Shutdown(objects.ComputeIteratorError(objects.ErrNotIterator, itObj.TypeName()))
		return
	}
	op.vm.StackPush(it.Key(op.vm.FrameId()))
}

// Compile generates the compiled representation of the OpIteratorKey operation or returns an unimplemented error.
func (op *OpIteratorKey) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
