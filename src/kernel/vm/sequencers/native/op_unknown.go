package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpUnknown)
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*bytecode.Opcode
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding Opcode configuration set.
func NewOpUnknown(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpUnknown{Opcode: op.Opcode(bytecode.OpUnknown)}
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 0
	v.SetError(fmt.Errorf("unknown opcode at: %d", v.GetIp()))
}
