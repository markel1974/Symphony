package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpPop)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpPop{
		Opcode: op.Opcode(bytecode.OpPop),
		vm:     vmT,
	}, nil
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.Stack().Decrement()
}
