package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalGet)
}

// OpLocalGet represents an operation to retrieve a local variable from the stack using its index.
type OpLocalGet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpLocalGet creates a new OpLocalGet instance and initializes it with details for the OpLocalGet opcode.
func NewOpLocalGet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalGet{
		Opcode: op.Opcode(bytecode.OpLocalGet),
		vm:     vm,
	}
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpLocalGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	op.vm.Stack().Push(val)
}
