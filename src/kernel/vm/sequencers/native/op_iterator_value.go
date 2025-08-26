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
// It embeds OpcodeDetails, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated OpcodeDetails initialized.
func NewOpIteratorValue(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIteratorValue{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorValue)}
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := v.Stack().PeekAbsolute(v.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.Stack().Push(iterator.Value(v.FrameID()))
}
