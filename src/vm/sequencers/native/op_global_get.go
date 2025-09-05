package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpGlobalGet)
}

// OpGlobalGet represents an operation to retrieve a global variable in the virtual machine.
// It embeds Opcode for detailed opcode information.
type OpGlobalGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalGet creates a new instance of OpGlobalGet with its associated opcode details.
func NewOpGlobalGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalGet{
		Opcode: op.Opcode(bytecode.OpGlobalGet),
		vm:     vmT,
	}, nil
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGlobalGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	glIndex := decoder.Read(0)
	glObj := op.vm.Globals().Get(uint(glIndex))
	if glObj == nil {
		op.vm.SetError(fmt.Errorf("undefined global: %d", glIndex))
		return
	}
	op.vm.Stack().Push(glObj)
}
