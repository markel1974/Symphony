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
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated Opcode initialized.
func NewOpIteratorValue(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIteratorValue{
		Opcode: op.Opcode(opcodes.OpIteratorValue),
		vm:     vmT,
	}, nil
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	op.vm.Stack().Push(iterator.Value(op.vm.Frame().Id()))
}
