package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpGlobalCopy function as a sequencer operation during package initialization.
func init() {
	SequencerRegister(NewOpGlobalCopy)
}

// OpGlobalCopy represents an operation that copies a value from one global variable index to another in the VM's global state.
// It embeds the bytecode.Opcode for opcode-related metadata and uses core.IVMFullAccess for VM interactions.
type OpGlobalCopy struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpGlobalCopy initializes an OpGlobalCopy executor with the provided VM and Opcodes instance or returns an error.
func NewOpGlobalCopy() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpGlobalCopy) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation to copy a global variable from sourceIndex to destIndex in the virtual machine.
func (op *OpGlobalCopy) Execute(decoder *core.Decoder) {
	sourceIndex := decoder.Operand(1)
	destIndex := decoder.Operand(0)
	value := op.vm.Globals().Get(uint(sourceIndex))
	op.vm.Globals().Set(uint(destIndex), value)
}

// Compile generates the compiled representation of the OpGlobalCopy operation or returns an unimplemented error.
func (op *OpGlobalCopy) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
