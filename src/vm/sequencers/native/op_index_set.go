package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
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
	vm     handler.IVMFullAccess
}

// NewOpIndexSet creates a new instance of OpIndexGet, binding it to the virtual machine and initializing with the given opcode.
// Returns an implementation of handler.IOpExecutor or an error if the Core does not support IVMFullAccess.
func NewOpIndexSet() handler.IOpExecutor {
	operands := _noOperands
	return &OpIndexSet{
		opcode: opcodes.NewOpcode(OpIndexSetId, operands, "OpIndexSet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIndexSet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIndexSet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute modifies a container's index with a new value, setting an error in the virtual machine if the operation fails.
func (op *OpIndexSet) Execute(_ *handler.Decoder) {
	value := op.vm.StackPop()
	index := op.vm.StackPop()
	container := op.vm.StackPop()
	if err := container.IndexSet(index, value); err != nil {
		op.vm.Shutdown(objects.ComputeIndexSetError(err, container.TypeName(), index.TypeName()))
		return
	}
}

// Compile generates the compiled representation of the OpIndexSet operation or returns an unimplemented error.
func (op *OpIndexSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
