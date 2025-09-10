package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpUnaryBitwiseComplement)
}

// OpUnaryBitwiseComplement represents an operation for performing a bitwise complement on an operand.
// It extends Opcode, inheriting its metadata and behaviors.
type OpUnaryBitwiseComplement struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpUnaryBitwiseComplement initializes and returns an OpUnaryBitwiseComplement instance with the corresponding Opcode configuration.
func NewOpUnaryBitwiseComplement(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpUnaryBitwiseComplement{
		Opcode: op.Opcode(opcodes.OpUnaryBitwiseComplement),
		vm:     vmT,
	}, nil
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpUnaryBitwiseComplement) Execute(_ *core.Decoder) {
	// Operands Offset 0
	operand := op.vm.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.vm.Factory().NewInt(op.vm.Frame().Id(), ^x.Value())
		op.vm.Stack().Push(res)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: ^%s", operand.TypeName()))
		return
	}
}
