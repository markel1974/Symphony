package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpDefineLocal represents the opcode for defining a new local variable within the current frame's scope.
type OpDefineLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpDefineLocal creates a new instance of OpDefineLocal with its associated opcode details.
func NewOpDefineLocal(op *bytecode.Opcodes) *OpDefineLocal {
	return &OpDefineLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpDefineLocal)}
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpDefineLocal) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := v.Stack().Peek()
	destSlot := v.BasePointer() + localIndex
	v.Stack().SetAbsolute(destSlot, val)
}
