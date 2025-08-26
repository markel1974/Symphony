package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpNoOp represents a no-operation opcode, typically used as a placeholder or for alignment purposes.
type OpNoOp struct {
	*bytecode.OpcodeDetails
}

// NewOpNoOp initializes and returns a new OpNoOp instance using the given Opcodes configuration.
func NewOpNoOp(op *bytecode.Opcodes) *OpNoOp {
	return &OpNoOp{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNoOp)}
}

// Execute performs a no-operation (NOP) for the virtual machine, advancing the instruction pointer without side effects.
func (op *OpNoOp) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
}
