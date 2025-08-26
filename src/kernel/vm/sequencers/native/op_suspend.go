package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpSuspend)
}

// OpSuspend represents an operation that suspends the execution of the virtual machine.
type OpSuspend struct {
	*bytecode.Opcode
}

// NewOpSuspend creates and returns a new OpSuspend instance with opcode details initialized for the suspend operation.
func NewOpSuspend(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpSuspend{Opcode: op.Opcode(bytecode.OpSuspend)}
}

// Execute performs the suspend operation on the given virtual machine by setting its shutdown state to true.
func (op *OpSuspend) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	v.Shutdown()
}
