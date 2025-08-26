package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpOr)
}

// OpJumpOr represents an operation that performs a logical OR and jumps based on the result.
type OpJumpOr struct {
	*bytecode.OpcodeDetails
}

// NewOpJumpOr creates and returns a new instance of OpJumpOr, associated with the OpJumpOr opcode and its details.
func NewOpJumpOr(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJumpOr{OpcodeDetails: op.OpcodeToDetails(bytecode.OpJumpOr)}
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpJumpOr) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := v.Stack().Peek()
	if obj.Boolean() {
		v.Stack().Decrement()
	} else {
		pos := decoder.Read(0)
		v.SetIp(pos - 1)
	}
}
