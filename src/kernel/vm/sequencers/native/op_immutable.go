package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpImmutable)
}

// OpImmutable represents an operation that creates immutable objects, inheriting details from Opcode.
type OpImmutable struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpImmutable creates a new instance of OpImmutable with details loaded from bytecode.OpcodeToDetails.
func NewOpImmutable(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpImmutable{
		Opcode: op.Opcode(bytecode.OpImmutable),
		vm:     vm,
	}
}

// Execute processes the top element on the stack and converts it into an immutable version if it's an array or map.
func (op *OpImmutable) Execute(_ *core.Decoder) {
	// Operands Offset  0
	val := op.vm.Stack().Peek()
	switch value := val.(type) {
	case *objects.Array:
		obj := op.vm.Factory().NewArrayImmutable(op.vm.Frame().Id(), value.Values())
		op.vm.Stack().Set(obj)
	case *objects.Map:
		obj := op.vm.Factory().NewMapImmutable(op.vm.Frame().Id(), value.Values())
		op.vm.Stack().Set(obj)
	}
}
