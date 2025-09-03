package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpIteratorKey)
}

// OpIteratorKey wraps bytecode.Opcode to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIteratorKey{
		Opcode: op.Opcode(bytecode.OpIteratorKey),
		vm:     vmT,
	}, nil
}

// Execute processes the "iterator key" operation, retrieves the iterator key, and pushes it onto the VM stack.
func (op *OpIteratorKey) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	op.vm.Stack().Push(iterator.Key(op.vm.Frame().Id()))
}
