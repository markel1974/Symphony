package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers the NewOpIndexSet operation with the sequencer system by adding it to the internal register container.
func init() {
	SequencerRegister(NewOpIndexSet)
}

// OpIndexSet represents an operation for setting a value at a specified index in a container within a virtual machine.
// It holds a reference to the Opcode and the IVMFullAccess interfaces needed for execution.
type OpIndexSet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIndexSet creates a new instance of OpIndexGet, binding it to the virtual machine and initializing with the given opcode.
// Returns an implementation of core.IOpExecutor or an error if the VM does not support IVMFullAccess.
func NewOpIndexSet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIndexGet{
		Opcode: op.Opcode(bytecode.OpIndexSet),
		vm:     vmT,
	}, nil
}

// Execute modifies a container's index with a new value, setting an error in the virtual machine if the operation fails.
func (op *OpIndexSet) Execute(_ *core.Decoder) {
	value := op.vm.Stack().Pop()
	index := op.vm.Stack().Pop()
	container := op.vm.Stack().Pop()
	if err := container.IndexSet(index, value); err != nil {
		op.vm.SetError(objects.ComputeIndexSetError(err, container.TypeName(), index.TypeName()))
		return
	}
}
