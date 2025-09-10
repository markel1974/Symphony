package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorInit)
}

// OpIteratorInit represents an operation that initializes an iterator over an iterable object.
// It embeds Opcode for additional opcode-specific metadata.
type OpIteratorInit struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIteratorInit{
		Opcode: op.Opcode(opcodes.OpIteratorInit),
		vm:     vmT,
	}, nil
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iterable := op.vm.Stack().Pop()
	if !iterable.Iterable() {
		op.vm.SetError(fmt.Errorf("not iterable: %s", iterable.TypeName()))
		return
	}
	iterator := iterable.Iterate(op.vm.Frame().Id())
	destSlot := op.vm.Frame().BasePointer() + localIndex
	op.vm.Stack().SetAbsolute(destSlot, iterator)
}
