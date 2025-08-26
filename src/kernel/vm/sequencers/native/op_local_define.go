package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpLocalDefine)
}

// OpLocalDefine represents the opcode for defining a new local variable within the current frame's scope.
type OpLocalDefine struct {
	*bytecode.Opcode
}

// NewOpLocalDefine creates a new instance of OpLocalDefine with its associated opcode details.
func NewOpLocalDefine(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalDefine{Opcode: op.Opcode(bytecode.OpLocalDefine)}
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpLocalDefine) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := v.Stack().Peek()
	destSlot := v.Frame().BasePointer() + localIndex
	v.Stack().SetAbsolute(destSlot, val)
}
