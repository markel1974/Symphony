package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpLocalDefine)
}

// OpLocalDefine represents the opcode for defining a new local variable within the current frame's scope.
type OpLocalDefine struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpLocalDefine creates a new instance of OpLocalDefine with its associated opcode details.
func NewOpLocalDefine(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalDefine{
		Opcode: op.Opcode(bytecode.OpLocalDefine),
		vm:     vm,
	}
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpLocalDefine) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := op.vm.Stack().Peek()
	destSlot := op.vm.Frame().BasePointer() + localIndex
	op.vm.Stack().SetAbsolute(destSlot, val)
}
