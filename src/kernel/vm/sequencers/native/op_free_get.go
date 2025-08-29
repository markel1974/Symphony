package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeGet)
}

// OpFreeGet represents an operation to retrieve a free variable in a closure during execution.
type OpFreeGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpFreeGet creates and returns a new instance of OpFreeGet, initializing its Opcode using bytecode metadata.
func NewOpFreeGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFreeGet{
		Opcode: op.Opcode(bytecode.OpFreeGet),
		vm:     vmT,
	}, nil
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpFreeGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := *op.vm.Frame().FreeVarsIndex(freeIndex).Value()
	op.vm.Stack().Push(val)
}
