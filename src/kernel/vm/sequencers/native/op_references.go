package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpReferences)
}

// OpReferences extends OpcodeDetails to represent operations specifically related to reference handling in the bytecode.
type OpReferences struct {
	*bytecode.OpcodeDetails
}

// NewOpReferences initializes a new OpReferences instance with corresponding OpcodeDetails from the bytecode package.
func NewOpReferences(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpReferences{OpcodeDetails: op.OpcodeToDetails(bytecode.OpReferences)}
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpReferences) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	nameIndex := decoder.Read(0)
	symbol := v.References().Get(uint(nameIndex))
	v.Stack().Push(symbol)
}
