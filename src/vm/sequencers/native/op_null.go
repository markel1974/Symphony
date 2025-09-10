package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpNull)
}

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpNull{
		Opcode: op.Opcode(opcodes.OpNull),
		vm:     vmT,
	}, nil
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(_ *core.Decoder) {
	// Operands Offset 0
	val := op.vm.Factory().UndefinedValue()
	op.vm.Stack().Push(val)
}
