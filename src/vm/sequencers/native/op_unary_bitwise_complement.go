package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
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
	vm     handler.IVMFullAccess
}

// NewOpUnaryBitwiseComplement initializes and returns an OpUnaryBitwiseComplement instance with the corresponding Opcode configuration.
func NewOpUnaryBitwiseComplement() handler.IOpExecutor {
	operands := _noOperands
	return &OpUnaryBitwiseComplement{
		opcode: opcodes.NewOpcode(OpUnaryBitwiseComplementId, operands, "OpUnaryBitwiseComplement"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnaryBitwiseComplement) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpUnaryBitwiseComplement) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpUnaryBitwiseComplement) Execute(_ *handler.Decoder) {
	obj := op.vm.StackPop()
	res := op.vm.Factory().NewInt(op.vm.FrameId(), ^obj.AsInt64())
	op.vm.StackPush(res)
}

// Compile generates the compiled representation of the OpUnaryBitwiseComplement operation or returns an unimplemented error.
func (op *OpUnaryBitwiseComplement) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
