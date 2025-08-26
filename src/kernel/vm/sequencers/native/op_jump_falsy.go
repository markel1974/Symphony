package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpFalsy)
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	*bytecode.Opcode
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding Opcode.
func NewOpJumpFalsy(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJumpFalsy{Opcode: op.Opcode(bytecode.OpJumpFalsy)}
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := v.Stack().Pop()
	if obj.Boolean() {
		pos := decoder.Read(0)
		v.SetIp(pos - 1)
	}
}
