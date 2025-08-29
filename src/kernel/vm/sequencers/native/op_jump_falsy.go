package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpFalsy)
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding Opcode.
func NewOpJumpFalsy(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJumpFalsy{
		Opcode: op.Opcode(bytecode.OpJumpFalsy),
		vm:     vm,
	}
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := op.vm.Stack().Pop()
	if obj.Falsy() {
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	}
}
