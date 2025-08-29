package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpNoOp)
}

// OpNoOp represents a no-operation opcode, typically used as a placeholder or for alignment purposes.
type OpNoOp struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpNoOp initializes and returns a new OpNoOp instance using the given Opcodes configuration.
func NewOpNoOp(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpNoOp{
		Opcode: op.Opcode(bytecode.OpNoOp),
		vm:     vm,
	}
}

// Execute performs a no-operation (NOP) for the virtual machine, advancing the instruction pointer without side effects.
func (op *OpNoOp) Execute(_ *core.Decoder) {
	// Operands Offset 0
}
