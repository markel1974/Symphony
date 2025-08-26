package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFalse)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	*bytecode.Opcode
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFalse{Opcode: op.Opcode(bytecode.OpFalse)}
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	val := v.Factory().FalseValue()
	v.Stack().Push(val)
}
