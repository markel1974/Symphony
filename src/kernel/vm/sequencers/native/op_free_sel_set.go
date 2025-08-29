package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpFreeSelSet)
}

// OpFreeSelSet represents an operation to set a free variable's value using selectors.
type OpFreeSelSet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpFreeSelSet creates a new instance of OpFreeSelSet with initialized Opcode referencing OpFreeSelSet.
func NewOpFreeSelSet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeSelSet{
		Opcode: op.Opcode(bytecode.OpFreeSelSet),
		vm:     vm,
	}
}

// Execute updates the instruction pointer, retrieves operands, processes selectors, and performs indexed assignment in the VM.
func (op *OpFreeSelSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	numSelectors := decoder.Read(0)
	freeIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = op.vm.Stack().PeekOffset(-numSelectors + i)
	}
	val := op.vm.Stack().PeekOffset(-numSelectors - 1)
	op.vm.Stack().DecrementCount(numSelectors + 1)
	fvi := op.vm.Frame().FreeVarsIndex(freeIndex)
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), *fvi.Value(), val, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
