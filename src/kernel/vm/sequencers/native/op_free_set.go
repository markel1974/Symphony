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
	*bytecode.OpcodeDetails
}

// NewOpFreeSet creates and returns a new instance of OpFreeSet initialized with its corresponding OpcodeDetails.
func NewOpFreeSet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeSet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpFreeSet)}
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpFreeSet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	o := v.Stack().Pop()
	v.Frame().FreeVarsIndex(freeIndex).SetValue(o)
}
