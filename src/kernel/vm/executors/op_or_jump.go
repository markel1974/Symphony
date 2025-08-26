package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpOrJump represents an operation that performs a logical OR and jumps based on the result.
type OpOrJump struct {
	*bytecode.OpcodeDetails
}

// NewOpOrJump creates and returns a new instance of OpOrJump, associated with the OpOrJump opcode and its details.
func NewOpOrJump(op *bytecode.Opcodes) *OpOrJump {
	return &OpOrJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpOrJump)}
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpOrJump) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := v.Stack().Peek()
	if obj.Boolean() {
		v.Stack().Decrement()
	} else {
		pos := decoder.Read(0)
		v.SetIp(pos - 1)
	}
}
