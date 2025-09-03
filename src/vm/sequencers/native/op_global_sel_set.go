package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpGlobalSelSet)
}

// OpGlobalSelSet represents an operation for setting a global variable's value using selectors for indexing or access.
type OpGlobalSelSet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalSelSet creates a new instance of OpGlobalSelSet with its corresponding Opcode initialized.
func NewOpGlobalSelSet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalSelSet{
		Opcode: op.Opcode(bytecode.OpGlobalSelSet),
		vm:     vmT,
	}, nil
}

// Execute performs the operation defined by OpGlobalSelSet, updating the VM state and handling global index assignment.
func (op *OpGlobalSelSet) Execute(decoder *core.Decoder) {
	// Operands Offset 3 (8-bit | 16bit)
	numSelectors := decoder.Read(0)
	globalIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = op.vm.Stack().PeekOffset(-numSelectors + i)
	}
	val := op.vm.Stack().PeekOffset(-numSelectors - 1)
	op.vm.Stack().DecrementCount(numSelectors + 1)
	glObj := op.vm.Globals().Get(uint(globalIndex))
	if err := op.vm.Factory().IndexAssign(op.vm.Frame().Id(), glObj, val, selectors); err != nil {
		op.vm.SetError(err)
		return
	}
}
