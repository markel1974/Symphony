package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeGet)
}

// OpFreeGet represents an operation to retrieve a free variable in a closure during execution.
type OpFreeGet struct {
	*bytecode.OpcodeDetails
}

// NewOpFreeGet creates and returns a new instance of OpFreeGet, initializing its OpcodeDetails using bytecode metadata.
func NewOpFreeGet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeGet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpFreeGet)}
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpFreeGet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := *v.FreeVarsIndex(freeIndex).Value()
	v.Stack().Push(val)
}
