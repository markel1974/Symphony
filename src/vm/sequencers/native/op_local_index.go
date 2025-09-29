package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalIndex)
}

// OpLocalIndex represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds Opcode to utilize its properties like opcode, name, and operands.
type OpLocalIndex struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLocalIndex creates and returns a new instance of the OpLocalIndex operation executor.
func NewOpLocalIndex() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.SzUint16}
	return &OpLocalIndex{
		opcode: opcodes.NewOpcode(OpLocalIndexId, operands, "OpLocalIndex"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalIndex) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLocalIndex) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpLocalIndex) Execute(decoder *handler.Decoder) {
	localIndex := decoder.Operand(0)
	selCount := decoder.Operand(1)
	dstObj := op.vm.StackPeekBP(uint(localIndex))
	selectors := make([]objects.IObject, selCount)
	for i := 0; i < selCount; i++ {
		offset := selCount - i
		selectors[i] = op.vm.StackPeekSP(uint(offset))
	}
	offset := selCount + 1
	srcObj := op.vm.StackPeekSP(uint(offset))
	op.vm.StackDecrementCount(uint(offset))
	if err := op.vm.Factory().IndexAssign(op.vm.FrameId(), dstObj, srcObj, selectors); err != nil {
		op.vm.Shutdown(err)
		return
	}
}

// Compile generates the compiled representation of the OpLocalIndex operation or returns an unimplemented error.
func (op *OpLocalIndex) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
