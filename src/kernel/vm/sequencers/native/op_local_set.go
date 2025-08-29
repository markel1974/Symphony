package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalSet)
}

// OpLocalSet represents an operation to set the value of a local variable within the current frame.
// It embeds Opcode for opcode-specific information such as name, operands, and code.
type OpLocalSet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpLocalSet initializes and returns a new instance of OpLocalSet with associated opcode details.
func NewOpLocalSet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalSet{
		Opcode: op.Opcode(bytecode.OpLocalSet),
		vm:     vm,
	}
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpLocalSet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := op.vm.Stack().Peek()
	destSlot := op.vm.Frame().BasePointer() + localIndex
	existingValue := op.vm.Stack().PeekAbsolute(destSlot)
	if obj, ok := existingValue.(*objects.ObjectPointer); ok {
		obj.SetValue(val)
	} else {
		op.vm.Stack().SetAbsolute(destSlot, val)
	}
}
