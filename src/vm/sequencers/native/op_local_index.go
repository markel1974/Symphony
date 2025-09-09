package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalIndex)
}

// OpLocalIndex represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds Opcode to utilize its properties like opcode, name, and operands.
type OpLocalIndex struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpLocalIndex creates and returns a new instance of the OpLocalIndex operation executor.
func NewOpLocalIndex(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpLocalIndex{
		Opcode: op.Opcode(bytecode.OpLocalIndex),
		vm:     vmT,
	}, nil
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpLocalIndex) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	localIndex := decoder.Read(0)
	selCount := decoder.Read(1)
	dstObj := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	//if obj, ok := dstObj.(*objects.ObjectPointer); ok {
	//	dstObj = *obj.Value()
	//}
	selectors := make([]objects.IObject, selCount)
	for i := 0; i < selCount; i++ {
		target := -selCount + i
		selectors[i] = op.vm.Stack().PeekOffset(target)
	}
	srcObj := op.vm.Stack().PeekOffset(-selCount - 1)
	op.vm.Stack().DecrementCount(selCount + 1)
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), dstObj, srcObj, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
