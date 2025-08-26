package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	*bytecode.OpcodeDetails
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue(op *bytecode.Opcodes) *OpTrue {
	return &OpTrue{OpcodeDetails: op.OpcodeToDetails(bytecode.OpTrue)}
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	val := op.Factory().TrueValue()
	v.Stack().Push(val)
}
