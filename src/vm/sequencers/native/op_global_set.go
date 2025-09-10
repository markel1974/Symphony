package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalSet)
}

// OpGlobalSet represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpGlobalSet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalSet creates and returns a new instance of OpGlobalSet with initialized Opcode.
func NewOpGlobalSet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalSet{
		Opcode: op.Opcode(opcodes.OpGlobalSet),
		vm:     vmT,
	}, nil
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpGlobalSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	index := decoder.Read(0)
	val := op.vm.Stack().Peek()
	op.vm.Globals().Set(uint(index), val)
}
