package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpImmutable)
}

// OpImmutable represents an operation that creates immutable objects, inheriting details from OpcodeDetails.
type OpImmutable struct {
	*bytecode.OpcodeDetails
}

// NewOpImmutable creates a new instance of OpImmutable with details loaded from bytecode.OpcodeToDetails.
func NewOpImmutable(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpImmutable{OpcodeDetails: op.OpcodeToDetails(bytecode.OpImmutable)}
}

// Execute processes the top element on the stack and converts it into an immutable version if it's an array or map.
func (op *OpImmutable) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	val := v.Stack().Peek()
	switch value := val.(type) {
	case *objects.Array:
		obj := op.Factory().NewArrayImmutable(v.FrameID(), value.Values())
		v.Stack().Set(obj)
	case *objects.Map:
		obj := op.Factory().NewMapImmutable(v.FrameID(), value.Values())
		v.Stack().Set(obj)
	}
}
