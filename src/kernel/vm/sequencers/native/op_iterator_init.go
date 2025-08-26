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
// It embeds OpcodeDetails for additional opcode-specific metadata.
type OpIteratorInit struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIteratorInit{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorInit)}
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iterable := v.Stack().Pop()
	if !iterable.CanIterate() {
		v.SetError(fmt.Errorf("not iterable: %s", iterable.TypeName()))
		return
	}
	iterator := iterable.Iterate(v.FrameID())
	destSlot := v.BasePointer() + localIndex
	v.Stack().SetAbsolute(destSlot, iterator)
}
