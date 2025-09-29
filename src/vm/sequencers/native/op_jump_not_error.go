// File: vm/sequencers/native/op_jump_if_not_error.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpJumpNotError operation with the sequencer system, enabling its functionality in the virtual machine.
func init() {
	SequencerRegister(NewOpJumpNotError)
}

// OpJumpNotError represents an operation that conditionally jumps if the top stack value is not a valid error.
// It uses Core control flow and stack operations to determine execution flow based on error type validity.
type OpJumpNotError struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpJumpNotError creates a new OpJumpNotError executor, verifying that the provided Core implements IVMFullAccess.
// It returns an instance of IOpExecutor or an error if the Core does not support full access functionality.
func NewOpJumpNotError() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpNotError{
		opcode: opcodes.NewOpcode(OpJumpNotErrorId, operands, "OpJumpNotError"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpNotError) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpJumpNotError) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute evaluates the top stack element and conditionally updates the instruction pointer if it is not a valid error.
func (op *OpJumpNotError) Execute(decoder *handler.Decoder) {
	obj := op.vm.StackPeek()
	err, isErr := obj.(*objects.Error)
	if !isErr || !err.Falsy() {
		pos := decoder.Operand(0)
		op.vm.SetIp(uint(pos))
	}
}

// Compile generates the compiled representation of the OpJumpNotError operation or returns an unimplemented error.
func (op *OpJumpNotError) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
