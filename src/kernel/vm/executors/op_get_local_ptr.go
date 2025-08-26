package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpGetLocalPtr retrieves a local variable as a pointer using its index within the current frame.
type OpGetLocalPtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocalPtr creates and returns a new instance of OpGetLocalPtr, initializing its OpcodeDetails.
func NewOpGetLocalPtr(op *bytecode.Opcodes) *OpGetLocalPtr {
	return &OpGetLocalPtr{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetLocalPtr)}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpGetLocalPtr) Execute(v *core.VM, decoder *core.Decoder) {
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
