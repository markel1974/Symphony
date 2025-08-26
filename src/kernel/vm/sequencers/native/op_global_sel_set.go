package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpGlobalSelSet)
}

// OpGlobalSelSet represents an operation for setting a global variable's value using selectors for indexing or access.
type OpGlobalSelSet struct {
	*bytecode.OpcodeDetails
}

// NewOpGlobalSelSet creates a new instance of OpGlobalSelSet with its corresponding OpcodeDetails initialized.
func NewOpGlobalSelSet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpGlobalSelSet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGlobalSelSet)}
}

// Execute performs the operation defined by OpGlobalSelSet, updating the VM state and handling global index assignment.
func (op *OpGlobalSelSet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 3 (8-bit | 16bit)
	numSelectors := decoder.Read(0)
	globalIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.Stack().PeekOffset(-numSelectors + i)
	}
	val := v.Stack().PeekOffset(-numSelectors - 1)
	v.Stack().DecrementCount(numSelectors + 1)
	glObj := v.Globals().Get(uint(globalIndex))
	if err := op.Factory().IndexAssign(v.FrameID(), glObj, val, selectors); err != nil {
		v.SetError(err)
		return
	}
}
