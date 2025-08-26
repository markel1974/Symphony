package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIteratorValue)
}

// OpIteratorValue retrieves the value from the current iterator position.
// It embeds Opcode, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	*bytecode.Opcode
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated Opcode initialized.
func NewOpIteratorValue(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIteratorValue{Opcode: op.Opcode(bytecode.OpIteratorValue)}
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := v.Stack().PeekAbsolute(v.Frame().BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.Stack().Push(iterator.Value(v.Frame().Id()))
}
