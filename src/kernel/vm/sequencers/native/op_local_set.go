package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalSet)
}

// OpLocalSet represents an operation to set the value of a local variable within the current frame.
// It embeds Opcode for opcode-specific information such as name, operands, and code.
type OpLocalSet struct {
	*bytecode.Opcode
}

// NewOpLocalSet initializes and returns a new instance of OpLocalSet with associated opcode details.
func NewOpLocalSet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalSet{Opcode: op.Opcode(bytecode.OpLocalSet)}
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpLocalSet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := v.Stack().Peek()
	destSlot := v.Frame().BasePointer() + localIndex
	existingValue := v.Stack().PeekAbsolute(destSlot)
	if obj, ok := existingValue.(*objects.ObjectPointer); ok {
		obj.SetValue(val)
	} else {
		v.Stack().SetAbsolute(destSlot, val)
	}
}
