package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpGetFreePtr represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds OpcodeDetails, which provides opcode metadata such as identifier, operands, and name.
type OpGetFreePtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFreePtr creates a new instance of OpGetFreePtr initialized with the corresponding OpcodeDetails.
func NewOpGetFreePtr(op *bytecode.Opcodes) *OpGetFreePtr {
	return &OpGetFreePtr{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetFreePtr)}
}

// Execute executes the OpGetFreePtr operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpGetFreePtr) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := v.FreeVarsIndex(freeIndex)
	v.Stack().Push(val)
}
