package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpIndexSet operation with the sequencer system by adding it to the internal register container.
func init() {
	SequencerRegister(NewOpIndexSet)
}

// OpIndexSet represents an operation for setting a value at a specified index in a container within a virtual machine.
// It holds a reference to the Opcode and the IVMFullAccess interfaces needed for execution.
type OpIndexSet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIndexSet creates a new instance of OpIndexGet, binding it to the virtual machine and initializing with the given opcode.
// Returns an implementation of core.IOpExecutor or an error if the VM does not support IVMFullAccess.
func NewOpIndexSet() core.IOpExecutor {
	operands := _noOperands
	return &OpIndexSet{
		opcode: opcodes.NewOpcode(OpIndexSetId, operands, "OpIndexSet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIndexSet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute modifies a container's index with a new value, setting an error in the virtual machine if the operation fails.
func (op *OpIndexSet) Execute(_ *core.Decoder) {
	value := op.vm.StackPop()
	index := op.vm.StackPop()
	container := op.vm.StackPop()
	if err := container.IndexSet(index, value); err != nil {
		op.vm.SetError(objects.ComputeIndexSetError(err, container.TypeName(), index.TypeName()))
		return
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIndexSet) Opcode() *opcodes.Opcode {
	return op.opcode
}
