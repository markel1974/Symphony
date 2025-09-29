package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpConstant)
}

// OpConstant represents an operation used to load a constant onto the stack.
type OpConstant struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpConstant creates a new OpConstant instance with opcode details initialized for the OpConstant operation.
func NewOpConstant() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpConstant{
		opcode: opcodes.NewOpcode(OpConstantId, operands, "OpConstant"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpConstant) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpConstant) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute executes the OpConstant instruction in the virtual machine, pushing a global constant onto the stack.
func (op *OpConstant) Execute(decoder *handler.Decoder) {
	cIdx := decoder.Operand(0)
	glObj := op.vm.ConstantsGet(uint(cIdx))
	op.vm.StackPush(glObj)
}

// Compile generates the compiled representation of the OpConstant operation or returns an unimplemented error.
func (op *OpConstant) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
