package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpTrue)
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpTrue{
		Opcode: op.Opcode(opcodes.OpTrue),
		vm:     vmT,
	}, nil
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(_ *core.Decoder) {
	// Operands Offset 0
	val := op.vm.Factory().TrueValue()
	op.vm.Stack().Push(val)
}
