package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the NewOpGlobalCopy function as a sequencer operation during package initialization.
func init() {
	SequencerRegister(NewOpGlobalCopy)
}

// OpGlobalCopy represents an operation that copies a value from one global variable index to another in the Core's global state.
// It embeds the bytecode.Opcode for opcode-related metadata and uses handler.IVMFullAccess for Core interactions.
type OpGlobalCopy struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpGlobalCopy initializes an OpGlobalCopy executor with the provided Core and Opcodes instance or returns an error.
func NewOpGlobalCopy() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable, opcodes.Relocatable}
	return &OpGlobalCopy{
		opcode: opcodes.NewOpcode(OpGlobalCopyId, operands, "OpGlobalCopy"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalCopy) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpGlobalCopy) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation to copy a global variable from sourceIndex to destIndex in the virtual machine.
func (op *OpGlobalCopy) Execute(decoder *handler.Decoder) {
	sourceIndex := decoder.Operand(1)
	destIndex := decoder.Operand(0)
	value, err := op.vm.GlobalsGet(uint(sourceIndex))
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	if err = op.vm.GlobalsSet(uint(destIndex), value); err != nil {
		op.vm.Shutdown(err)
		return
	}
}

// Compile generates the compiled representation of the OpGlobalCopy operation or returns an unimplemented error.
func (op *OpGlobalCopy) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
