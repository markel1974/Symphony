package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpGlobalSet)
}

// OpGlobalSet represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpGlobalSet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpGlobalSet creates and returns a new instance of OpGlobalSet with initialized Opcode.
func NewOpGlobalSet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpGlobalSet{
		Opcode: op.Opcode(bytecode.OpGlobalSet),
		vm:     vm,
	}
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpGlobalSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	pos := decoder.Read(0)
	val := op.vm.Stack().Peek()
	op.vm.Globals().Set(uint(pos), val)
}
