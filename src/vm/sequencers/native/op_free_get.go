package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeGet)
}

// OpFreeGet represents an operation to retrieve a free variable in a closure during execution.
type OpFreeGet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpFreeGet creates and returns a new instance of OpFreeGet, initializing its Opcode using bytecode metadata.
func NewOpFreeGet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFreeGet{
		Opcode: op.Opcode(opcodes.OpFreeGet),
		vm:     vmT,
	}, nil
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpFreeGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	freeIndex := decoder.Read(0)
	freeVar := op.vm.Frame().FreeVarsIndex(uint(freeIndex))
	if freeVar == nil {
		op.vm.SetError(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	z := *freeVar.Value()
	op.vm.Stack().Push(z)
}
