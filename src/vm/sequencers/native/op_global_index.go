package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalIndex)
}

// OpGlobalIndex represents an operation for setting a global variable's value using selectors for indexing or access.
type OpGlobalIndex struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpGlobalIndex creates a new instance of OpGlobalIndex with its corresponding Opcode initialized.
func NewOpGlobalIndex() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.Relocatable}
	return &OpGlobalIndex{
		opcode: opcodes.NewOpcode(OpGlobalIndexId, operands, "OpGlobalIndex"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpGlobalIndex) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation defined by OpGlobalIndex, updating the VM state and handling global index assignment.
func (op *OpGlobalIndex) Execute(decoder *core.Decoder) {
	// Operands Offset 3 (8-bit | 16bit)
	globalIndex := decoder.Operand(0)
	selCount := decoder.Operand(1)
	dstObj := op.vm.Globals().Get(uint(globalIndex))
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
	op.vm.StackDecrementCount(offset)
	if err := op.vm.Factory().IndexAssign(op.vm.FrameId(), dstObj, srcObj, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalIndex) Opcode() *opcodes.Opcode {
	return op.opcode
}
