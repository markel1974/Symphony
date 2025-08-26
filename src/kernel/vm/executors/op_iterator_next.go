package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext(op *bytecode.Opcodes) *OpIteratorNext {
	return &OpIteratorNext{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorNext)}
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := v.Stack().PeekAbsolute(v.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	if iterator.Next() {
		v.Stack().Push(op.Factory().TrueValue())
	} else {
		v.Stack().Push(op.Factory().FalseValue())
	}
}
