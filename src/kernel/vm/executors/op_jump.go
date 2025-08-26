package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJump)
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	*bytecode.OpcodeDetails
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpJump)}
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`VM`) to a calculated position in the frame.
func (op *OpJump) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	pos := decoder.Read(0)
	v.SetIp(pos - 1)
}
