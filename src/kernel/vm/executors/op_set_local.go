package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpSetLocal represents an operation to set the value of a local variable within the current frame.
// It embeds OpcodeDetails for opcode-specific information such as name, operands, and code.
type OpSetLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetLocal initializes and returns a new instance of OpSetLocal with associated opcode details.
func NewOpSetLocal(op *bytecode.Opcodes) *OpSetLocal {
	return &OpSetLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetLocal)}
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpSetLocal) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := v.Stack().Peek()
	destSlot := v.BasePointer() + localIndex
	existingValue := v.Stack().PeekAbsolute(destSlot)
	if obj, ok := existingValue.(*objects.ObjectPointer); ok {
		obj.SetValue(val)
	} else {
		v.Stack().SetAbsolute(destSlot, val)
	}
}
