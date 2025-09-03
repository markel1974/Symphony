package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpNoOp)
}

// OpNoOp represents a no-operation opcode, typically used as a placeholder or for alignment purposes.
type OpNoOp struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpNoOp initializes and returns a new OpNoOp instance using the given Opcodes configuration.
func NewOpNoOp(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpNoOp{
		Opcode: op.Opcode(bytecode.OpNoOp),
		vm:     vmT,
	}, nil
}

// Execute performs a no-operation (NOP) for the virtual machine, advancing the instruction pointer without side effects.
func (op *OpNoOp) Execute(_ *core.Decoder) {
	// Operands Offset 0
}
