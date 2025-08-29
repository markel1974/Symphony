package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeGetPtr)
}

// OpFreePtrGet represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds Opcode, which provides opcode metadata such as identifier, operands, and name.
type OpFreePtrGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpFreeGetPtr creates a new instance of OpFreePtrGet initialized with the corresponding Opcode.
func NewOpFreeGetPtr(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFreePtrGet{
		Opcode: op.Opcode(bytecode.OpFreePtrGet),
		vm:     vmT,
	}, nil
}

// Execute executes the OpFreePtrGet operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpFreePtrGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := op.vm.Frame().FreeVarsIndex(freeIndex)
	op.vm.Stack().Push(val)
}
