package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	*bytecode.OpcodeDetails
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull(op *bytecode.Opcodes) *OpNull {
	return &OpNull{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNull)}
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	val := op.Factory().UndefinedValue()
	v.Stack().Push(val)
}
