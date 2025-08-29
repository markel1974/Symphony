package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIteratorKey)
}

// OpIteratorKey wraps bytecode.Opcode to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIteratorKey{
		Opcode: op.Opcode(bytecode.OpIteratorKey),
		vm:     vm,
	}
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
