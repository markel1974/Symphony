package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpOr)
}

// OpJumpOr represents an operation that performs a logical OR and jumps based on the result.
type OpJumpOr struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpJumpOr creates and returns a new instance of OpJumpOr, associated with the OpJumpOr opcode and its details.
func NewOpJumpOr(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJumpOr{
		Opcode: op.Opcode(bytecode.OpJumpOr),
		vm:     vm,
	}
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
