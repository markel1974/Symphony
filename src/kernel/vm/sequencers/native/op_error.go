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
	*bytecode.Opcode
	vm *core.VM
}

// NewOpError creates and returns a new instance of OpError with associated Opcode for the OpError opcode.
func NewOpError(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpError{
		Opcode: op.Opcode(bytecode.OpError),
		vm:     vm,
	}
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpError) Execute(_ *core.Decoder) {
	// Operands Offset  0
	value := op.vm.Stack().Peek()
	e := op.vm.Factory().NewError(op.vm.Frame().Id(), value.String())
	op.vm.Stack().Set(e)
}
