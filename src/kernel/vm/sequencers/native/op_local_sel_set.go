package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalSelSet)
}

// OpLocalSelSet represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds Opcode to utilize its properties like opcode, name, and operands.
type OpLocalSelSet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpLocalSelSet creates and returns a new instance of the OpLocalSelSet operation executor.
func NewOpLocalSelSet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalSelSet{
		Opcode: op.Opcode(bytecode.OpLocalSelSet),
		vm:     vm,
	}
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpLocalSelSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	numSelectors := decoder.Read(0)
	localIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = op.vm.Stack().PeekOffset(-numSelectors + i)
	}
	val := op.vm.Stack().PeekOffset(-numSelectors - 1)
	op.vm.Stack().DecrementCount(numSelectors + 1)
	dst := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	if obj, ok := dst.(*objects.ObjectPointer); ok {
		dst = *obj.Value()
	}
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), dst, val, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
