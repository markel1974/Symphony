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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpUnaryBitwiseComplement initializes and returns an OpUnaryBitwiseComplement instance with the corresponding Opcode configuration.
func NewOpUnaryBitwiseComplement() core.IOpExecutor {
	operands := _noOperands
	return &OpUnaryBitwiseComplement{
		opcode: opcodes.NewOpcode(OpUnaryBitwiseComplementId, operands, "OpUnaryBitwiseComplement"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpUnaryBitwiseComplement) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
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

// Opcode returns the opcode associated with the instance.
func (op *OpUnaryBitwiseComplement) Opcode() *opcodes.Opcode {
	return op.opcode
}
