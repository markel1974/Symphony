package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalPtrGet)
}

// OpLocalPtrGet retrieves a local variable as a pointer using its index within the current frame.
type OpLocalPtrGet struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its Opcode.
func NewOpLocalPtrGet(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalPtrGet{
		Opcode: op.Opcode(bytecode.OpLocalPtrGet),
		vm:     vm,
	}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpLocalPtrGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	sp := op.vm.Frame().BasePointer() + localIndex
	val := op.vm.Stack().PeekAbsolute(sp)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		op.vm.Stack().Push(obj)
		return
	}
	freeVar := op.vm.Factory().NewObjectPointer(op.vm.Frame().Id(), &val)
	op.vm.Stack().SetAbsolute(sp, freeVar)
	op.vm.Stack().Push(freeVar)
}
