package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
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
	vm     core.IVMFullAccess
}

// NewOpLocalIndex creates and returns a new instance of the OpLocalIndex operation executor.
func NewOpLocalIndex() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.SzUint16}
	return &OpLocalIndex{
		opcode: opcodes.NewOpcode(OpLocalIndexId, operands, "OpLocalIndex"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLocalIndex) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpLocalIndex) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit|8-bit)
	localIndex := decoder.Operand(0)
	selCount := decoder.Operand(1)
	dstObj := op.vm.StackPeekOffsetBP(uint(localIndex))
	//if obj, ok := dstObj.(*objects.ObjectPointer); ok {
	//	dstObj = *obj.Value()
	//}
	selectors := make([]objects.IObject, selCount)
	for i := 0; i < selCount; i++ {
		offset := selCount - i
		selectors[i] = op.vm.StackPeekOffsetSP(uint(offset))
	}
	offset := selCount + 1
	srcObj := op.vm.StackPeekOffsetSP(uint(offset))
	op.vm.StackDecrementCount(uint(offset))
	if err := op.vm.Factory().IndexAssign(op.vm.FrameId(), dstObj, srcObj, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalIndex) Opcode() *opcodes.Opcode {
	return op.opcode
}
