package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpNoOp)
}

// OpNoOp represents a no-operation opcode, typically used as a placeholder or for alignment purposes.
type OpNoOp struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpNoOp initializes and returns a new OpNoOp instance using the given Opcodes configuration.
func NewOpNoOp() core.IOpExecutor {
	operands := _noOperands
	return &OpNoOp{
		opcode: opcodes.NewOpcode(OpNoOpId, operands, "OpNoOp"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpNoOp) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpNoOp) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs a no-operation (NOP) for the virtual machine, advancing the instruction pointer without side effects.
func (op *OpNoOp) Execute(_ *core.Decoder) {
	// Operands Offset 0
}

// Compile generates the compiled representation of the OpLogical operation or returns an unimplemented error.
func (op *OpNoOp) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
