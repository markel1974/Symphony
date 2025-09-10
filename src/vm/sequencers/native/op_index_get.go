package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIndexGet)
}

// OpIndexGet represents the operation for performing an indexing operation on a value.
type OpIndexGet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpIndexGet creates and returns a new instance of OpIndexGet initialized with its associated Opcode.
func NewOpIndexGet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIndexGet{
		Opcode: op.Opcode(opcodes.OpIndexGet),
		vm:     vmT,
	}, nil
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndexGet) Execute(_ *core.Decoder) {
	// Operands Offset  0
	index := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	val, err := left.IndexGet(op.vm.Frame().Id(), index)
	if err != nil {
		op.vm.SetError(objects.ComputeIndexGetError(err, index.TypeName(), index.TypeName()))
		return
	}
	if val == nil {
		val = op.vm.Factory().UndefinedValue()
	}
	op.vm.Stack().Push(val)
}
