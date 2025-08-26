package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpGetLocal represents an operation to retrieve a local variable from the stack using its index.
type OpGetLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocal creates a new OpGetLocal instance and initializes it with details for the OpGetLocal opcode.
func NewOpGetLocal(op *bytecode.Opcodes) *OpGetLocal {
	return &OpGetLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetLocal)}
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpGetLocal) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := v.Stack().PeekAbsolute(v.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.Stack().Push(val)
}
