package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpReturn)
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.OpcodeDetails.
type OpReturn struct {
	*bytecode.OpcodeDetails
}

// NewOpReturn creates a new instance of OpReturn with its OpcodeDetails initialized for the OpReturn operation.
func NewOpReturn(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpReturn{OpcodeDetails: op.OpcodeToDetails(bytecode.OpReturn)}
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the VM.
func (op *OpReturn) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	// collect return values from the stack using Pop(),
	// this is necessary to uncover the underlying values.
	var returnValues []objects.IObject
	if numReturnVals := decoder.Read(0); numReturnVals > 0 {
		returnValues = make([]objects.IObject, decoder.Read(0))
		for i := 0; i < numReturnVals; i++ {
			returnValues[i] = v.Stack().Pop()
		}
	}
	v.Return(returnValues)
}
