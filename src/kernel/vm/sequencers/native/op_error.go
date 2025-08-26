package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpError)
}

// OpError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpError struct {
	*bytecode.OpcodeDetails
}

// NewOpError creates and returns a new instance of OpError with associated OpcodeDetails for the OpError opcode.
func NewOpError(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpError{OpcodeDetails: op.OpcodeToDetails(bytecode.OpError)}
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpError) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	value := v.Stack().Peek()
	e := op.Factory().NewError(v.Frame().Id(), value.String())
	v.Stack().Set(e)
}
