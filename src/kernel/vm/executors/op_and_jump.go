package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpAndJump represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpAndJump struct {
	*bytecode.OpcodeDetails
}

// NewOpAndJump creates and returns a new instance of OpAndJump, initializing it with details for the OpAndJump opcode.
func NewOpAndJump(op *bytecode.Opcodes) *OpAndJump {
	return &OpAndJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpAndJump)}
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpAndJump) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	obj := v.Stack().Peek()
	if obj.Boolean() {
		pos := decoder.Read(0)
		v.SetIp(pos - 1)
	} else {
		v.Stack().Decrement()
	}
}
