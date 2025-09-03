package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalPtrGet)
}

// OpLocalPtrGet retrieves a local variable as a pointer using its index within the current frame.
type OpLocalPtrGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its Opcode.
func NewOpLocalPtrGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpLocalPtrGet{
		Opcode: op.Opcode(bytecode.OpLocalPtrGet),
		vm:     vmT,
	}, nil
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
	op.vm.Stack().Push(freeVar)
}
