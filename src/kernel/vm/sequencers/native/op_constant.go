package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpConstant)
}

// OpConstant represents an operation used to load a constant onto the stack.
type OpConstant struct {
	*bytecode.OpcodeDetails
}

// NewOpConstant creates a new OpConstant instance with opcode details initialized for the OpConstant operation.
func NewOpConstant(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpConstant{OpcodeDetails: op.OpcodeToDetails(bytecode.OpConstant)}
}

// Execute executes the OpConstant instruction in the virtual machine, pushing a global constant onto the stack.
func (op *OpConstant) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	cIdx := decoder.Read(0)
	glObj := v.Constants().Get(uint(cIdx))
	v.Stack().Push(glObj)
}
