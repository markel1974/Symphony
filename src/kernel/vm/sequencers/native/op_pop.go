package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpPop)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpPop{
		Opcode: op.Opcode(bytecode.OpPop),
		vm:     vm,
	}
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.Stack().Decrement()
}
