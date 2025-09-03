package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpLocalSelSet)
}

// OpLocalSelSet represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds Opcode to utilize its properties like opcode, name, and operands.
type OpLocalSelSet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpLocalSelSet creates and returns a new instance of the OpLocalSelSet operation executor.
func NewOpLocalSelSet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpLocalSelSet{
		Opcode: op.Opcode(bytecode.OpLocalSelSet),
		vm:     vmT,
	}, nil
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpLocalSelSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	numSelectors := decoder.Read(0)
	localIndex := decoder.Read(1)
	selectors := make([]objects2.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		target := -numSelectors + i
		selectors[i] = op.vm.Stack().PeekOffset(target)
	}
	val := op.vm.Stack().PeekOffset(-numSelectors - 1)
	op.vm.Stack().DecrementCount(numSelectors + 1)
	dst := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	if obj, ok := dst.(*objects2.ObjectPointer); ok {
		dst = *obj.Value()
	}
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), dst, val, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
