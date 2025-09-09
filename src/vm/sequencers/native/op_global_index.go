package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpGlobalIndex)
}

// OpGlobalIndex represents an operation for setting a global variable's value using selectors for indexing or access.
type OpGlobalIndex struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalIndex creates a new instance of OpGlobalIndex with its corresponding Opcode initialized.
func NewOpGlobalIndex(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalIndex{
		Opcode: op.Opcode(bytecode.OpGlobalIndex),
		vm:     vmT,
	}, nil
}

// Execute performs the operation defined by OpGlobalIndex, updating the VM state and handling global index assignment.
func (op *OpGlobalIndex) Execute(decoder *core.Decoder) {
	// Operands Offset 3 (8-bit | 16bit)
	globalIndex := decoder.Read(0)
	selCount := decoder.Read(1)
	dstObj := op.vm.Globals().Get(uint(globalIndex))
	//if obj, ok := dstObj.(*objects.ObjectPointer); ok {
	//	dstObj = *obj.Value()
	//}
	selectors := make([]objects.IObject, selCount)
	for i := 0; i < selCount; i++ {
		selectors[i] = op.vm.Stack().PeekOffset(-selCount + i)
	}
	srcObj := op.vm.Stack().PeekOffset(-selCount - 1)
	op.vm.Stack().DecrementCount(selCount + 1)
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), dstObj, srcObj, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
