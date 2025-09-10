package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalGet)
}

// OpGlobalGet represents an operation to retrieve a global variable in the virtual machine.
// It embeds Opcode for detailed opcode information.
type OpGlobalGet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalGet creates a new instance of OpGlobalGet with its associated opcode details.
func NewOpGlobalGet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalGet{
		Opcode: op.Opcode(opcodes.OpGlobalGet),
		vm:     vmT,
	}, nil
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGlobalGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	index := decoder.Read(0)
	obj := op.vm.Globals().Get(uint(index))
	op.vm.Stack().Push(obj)
}
