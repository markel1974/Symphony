package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpAnd)
}

// OpJumpAnd represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpJumpAnd struct {
	*bytecode.Opcode
}

// NewOpJumpAnd creates and returns a new instance of OpJumpAnd, initializing it with details for the OpJumpAnd opcode.
func NewOpJumpAnd(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJumpAnd{Opcode: op.Opcode(bytecode.OpJumpAnd)}
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpJumpAnd) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	obj := v.Stack().Peek()
	if obj.Boolean() {
		pos := decoder.Read(0)
		v.SetIp(pos - 1)
	} else {
		v.Stack().Decrement()
	}
}
