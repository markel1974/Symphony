package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpIteratorInit)
}

// OpIteratorInit represents an operation that initializes an iterator over an iterable object.
// It embeds Opcode for additional opcode-specific metadata.
type OpIteratorInit struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIteratorInit{
		Opcode: op.Opcode(bytecode.OpIteratorInit),
		vm:     vm,
	}
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iterable := op.vm.Stack().Pop()
	if !iterable.CanIterate() {
		op.vm.SetError(fmt.Errorf("not iterable: %s", iterable.TypeName()))
		return
	}
	iterator := iterable.Iterate(op.vm.Frame().Id())
	destSlot := op.vm.Frame().BasePointer() + localIndex
	op.vm.Stack().SetAbsolute(destSlot, iterator)
}
