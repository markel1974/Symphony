package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJumpOr)
}

// OpJumpOr represents an operation that performs a logical OR and jumps based on the result.
type OpJumpOr struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpOr creates and returns a new instance of OpJumpOr, associated with the OpJumpOr opcode and its details.
func NewOpJumpOr(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpOr{
		Opcode: op.Opcode(opcodes.OpJumpOr),
		vm:     vmT,
	}, nil
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpJumpOr) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := op.vm.Stack().Peek()
	if obj.Falsy() {
		op.vm.Stack().Decrement()
	} else {
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	}
}
