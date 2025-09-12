package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorKey)
}

// OpIteratorKey wraps bytecode.Opcode to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorKey{
		opcode: opcodes.NewOpcode(OpIteratorKeyId, operands, "OpIteratorKey"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIteratorKey) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the "iterator key" operation, retrieves the iterator key, and pushes it onto the VM stack.
func (op *OpIteratorKey) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Operand(0)
	iteratorObj := op.vm.StackPeekOffsetBP(uint(localIndex))
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	op.vm.StackPush(iterator.Key(op.vm.FrameId()))
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorKey) Opcode() *opcodes.Opcode {
	return op.opcode
}
