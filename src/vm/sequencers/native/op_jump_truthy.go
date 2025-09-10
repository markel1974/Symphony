package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the necessary operation for the sequencer by registering NewOpJumpTruthy with the SequencerRegister function.
func init() {
	SequencerRegister(NewOpJumpTruthy)
}

// OpJumpTruthy defines a structure for executing a truthy-check jump operation in a virtual machine context.
type OpJumpTruthy struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpTruthy creates a new OpJumpTruthy executor, validating the provided virtual machine for required full access capabilities.
func NewOpJumpTruthy(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpTruthy{
		Opcode: op.Opcode(opcodes.OpJumpTruthy),
		vm:     vmT,
	}, nil
}

// Execute evaluates the top object on the stack and updates the instruction pointer if it is truthy.
func (op *OpJumpTruthy) Execute(decoder *core.Decoder) {
	obj := op.vm.Stack().Pop() // Prendiamo il valore dalla stack
	if !obj.Falsy() {          // Controlliamo se è 'truthy'
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	}
}
