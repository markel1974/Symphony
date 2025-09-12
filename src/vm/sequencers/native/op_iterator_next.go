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
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Operand(0)
	iteratorObj := op.vm.StackPeekOffsetBP(uint(localIndex))
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		op.vm.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	if iterator.Next() {
		op.vm.StackPush(op.vm.Factory().TrueValue())
	} else {
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorNext) Opcode() *opcodes.Opcode {
	return op.opcode
}
