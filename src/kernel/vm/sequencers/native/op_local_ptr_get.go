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
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its Opcode.
func NewOpLocalPtrGet(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpLocalPtrGet{Opcode: op.Opcode(bytecode.OpLocalPtrGet)}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpLocalPtrGet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	sp := v.Frame().BasePointer() + localIndex
	val := v.Stack().PeekAbsolute(sp)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		v.Stack().Push(obj)
		return
	}
	freeVar := v.Factory().NewObjectPointer(v.Frame().Id(), &val)
	v.Stack().SetAbsolute(sp, freeVar)
	v.Stack().Push(freeVar)
}
