package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIndex)
}

// OpIndex represents the operation for performing an indexing operation on a value.
type OpIndex struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIndex creates and returns a new instance of OpIndex initialized with its associated Opcode.
func NewOpIndex(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIndex{
		Opcode: op.Opcode(bytecode.OpIndex),
		vm:     vmT,
	}, nil
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndex) Execute(_ *core.Decoder) {
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
