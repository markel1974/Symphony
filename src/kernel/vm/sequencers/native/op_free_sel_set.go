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
	*bytecode.OpcodeDetails
}

// NewOpFreeSelSet creates a new instance of OpFreeSelSet with initialized OpcodeDetails referencing OpFreeSelSet.
func NewOpFreeSelSet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeSelSet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpFreeSelSet)}
}

// Execute updates the instruction pointer, retrieves operands, processes selectors, and performs indexed assignment in the VM.
func (op *OpFreeSelSet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	numSelectors := decoder.Read(0)
	freeIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.Stack().PeekOffset(-numSelectors + i)
	}
	val := v.Stack().PeekOffset(-numSelectors - 1)
	v.Stack().DecrementCount(numSelectors + 1)
	fvi := v.FreeVarsIndex(freeIndex)
	if err := op.Factory().IndexAssign(v.FrameID(), *fvi.Value(), val, selectors); err != nil {
		v.SetError(err)
		return
	}
}
