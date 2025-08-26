package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpIndex represents the operation for performing an indexing operation on a value.
type OpIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpIndex creates and returns a new instance of OpIndex initialized with its associated OpcodeDetails.
func NewOpIndex(op *bytecode.Opcodes) *OpIndex {
	return &OpIndex{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIndex)}
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndex) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset  0
	index := v.Stack().Pop()
	left := v.Stack().Pop()
	val, err := left.IndexGet(v.FrameID(), index)
	if err != nil {
		if objects.Is(err, objects.ErrNotIndexable) {
			v.SetError(fmt.Errorf("not indexable: %s", index.TypeName()))
			return
		}
		if objects.Is(err, objects.ErrInvalidIndexType) {
			v.SetError(fmt.Errorf("invalid index type: %s", index.TypeName()))
			return
		}
		v.SetError(err)
		return
	}
	if val == nil {
		val = op.Factory().UndefinedValue()
	}
	v.Stack().Push(val)
}
