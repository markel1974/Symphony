package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpPop)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	*bytecode.OpcodeDetails
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpPop{OpcodeDetails: op.OpcodeToDetails(bytecode.OpPop)}
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	v.Stack().Decrement()
}
