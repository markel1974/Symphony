package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeGetPtr)
}

// OpFreePtrGet represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds OpcodeDetails, which provides opcode metadata such as identifier, operands, and name.
type OpFreePtrGet struct {
	*bytecode.OpcodeDetails
}

// NewOpFreeGetPtr creates a new instance of OpFreePtrGet initialized with the corresponding OpcodeDetails.
func NewOpFreeGetPtr(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreePtrGet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpFreePtrGet)}
}

// Execute executes the OpFreePtrGet operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpFreePtrGet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := v.Frame().FreeVarsIndex(freeIndex)
	v.Stack().Push(val)
}
