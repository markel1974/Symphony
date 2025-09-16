package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorNext)
}

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorNext{
		opcode: opcodes.NewOpcode(OpIteratorNextId, operands, "OpIteratorNext"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorNext) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIteratorNext) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(decoder *core.Decoder) {
	localIndex := decoder.Operand(0)
	iteratorObj := op.vm.StackPeekOffsetBP(uint(localIndex))
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		err := objects.ComputeIteratorError(objects.ErrNotIterator, iteratorObj.TypeName())
		op.vm.SetError(err)
		return
	}
	if iterator.Next() {
		op.vm.StackPush(op.vm.Factory().TrueValue())
	} else {
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Compile generates the compiled representation of the OpIteratorNext operation or returns an unimplemented error.
func (op *OpIteratorNext) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
