package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeSet)
}

// OpFreeSet represents an operation to set the value of a free variable within a closure's environment.
type OpFreeSet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpFreeSet creates and returns a new instance of OpFreeSet initialized with its corresponding Opcode.
func NewOpFreeSet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeSet{
		Opcode: op.Opcode(bytecode.OpFreeSet),
		vm:     vm,
	}
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpFreeSet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	o := op.vm.Stack().Pop()
	op.vm.Frame().FreeVarsIndex(freeIndex).SetValue(o)
}
