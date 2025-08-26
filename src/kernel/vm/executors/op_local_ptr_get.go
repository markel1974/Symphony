package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpLocalPtrGet retrieves a local variable as a pointer using its index within the current frame.
type OpLocalPtrGet struct {
	*bytecode.OpcodeDetails
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its OpcodeDetails.
func NewOpLocalPtrGet(op *bytecode.Opcodes) *OpLocalPtrGet {
	return &OpLocalPtrGet{OpcodeDetails: op.OpcodeToDetails(bytecode.OpLocalPtrGet)}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpLocalPtrGet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	sp := v.BasePointer() + localIndex
	val := v.Stack().PeekAbsolute(sp)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		v.Stack().Push(obj)
		return
	}
	freeVar := op.Factory().NewObjectPointer(v.FrameID(), &val)
	v.Stack().SetAbsolute(sp, freeVar)
	v.Stack().Push(freeVar)
}
