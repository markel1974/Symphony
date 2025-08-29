package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpFreeGet)
}

// OpFreeGet represents an operation to retrieve a free variable in a closure during execution.
type OpFreeGet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpFreeGet creates and returns a new instance of OpFreeGet, initializing its Opcode using bytecode metadata.
func NewOpFreeGet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpFreeGet{
		Opcode: op.Opcode(bytecode.OpFreeGet),
		vm:     vm,
	}
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpFreeGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	freeIndex := decoder.Read(0)
	val := *op.vm.Frame().FreeVarsIndex(freeIndex).Value()
	op.vm.Stack().Push(val)
}
