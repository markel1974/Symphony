package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpIteratorNext)
}

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIteratorNext{
		Opcode: op.Opcode(bytecode.OpIteratorNext),
		vm:     vmT,
	}, nil
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	iteratorObj := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	if iterator.Next() {
		op.vm.Stack().Push(op.vm.Factory().TrueValue())
	} else {
		op.vm.Stack().Push(op.vm.Factory().FalseValue())
	}
}
