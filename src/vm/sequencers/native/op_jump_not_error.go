// File: vm/sequencers/native/op_jump_if_not_error.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpJumpNotError operation with the sequencer system, enabling its functionality in the virtual machine.
func init() {
	SequencerRegister(NewOpJumpNotError)
}

// OpJumpNotError represents an operation that conditionally jumps if the top stack value is not a valid error.
// It uses VM control flow and stack operations to determine execution flow based on error type validity.
type OpJumpNotError struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpNotError creates a new OpJumpNotError executor, verifying that the provided VM implements IVMFullAccess.
// It returns an instance of IOpExecutor or an error if the VM does not support full access functionality.
func NewOpJumpNotError(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpNotError{
		Opcode: op.Opcode(opcodes.OpJumpNotError),
		vm:     vmT,
	}, nil
}

// Execute evaluates the top stack element and conditionally updates the instruction pointer if it is not a valid error.
func (op *OpJumpNotError) Execute(decoder *core.Decoder) {
	// The object to be checked is at the top of the stack. We don't remove it yet.
	obj := op.vm.Stack().Peek()

	// The Error type's Falsy() logic is perfect: it returns 'true' only for a valid error.
	// Here we use inverse logic: jump if the object is NOT a valid error.
	// objects.Error.Falsy() returns true, so we need to verify that the type is Error and Falsy() is true.
	err, isErr := obj.(*objects.Error)

	if !isErr || !err.Falsy() {
		// It's not an error, or it's a null error (Falsy() returns false).
		// In both cases, skip the if block.
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	}
	// If it's a valid error, execution continues in the if block.
	// The compiler will handle popping the error from the stack.
}
