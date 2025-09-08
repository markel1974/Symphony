package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

// init registers the NewOpGlobalCopy function as a sequencer operation during package initialization.
func init() {
	SequencerRegister(NewOpGlobalCopy)
}

// OpGlobalCopy represents an operation that copies a value from one global variable index to another in the VM's global state.
// It embeds the bytecode.Opcode for opcode-related metadata and uses core.IVMFullAccess for VM interactions.
type OpGlobalCopy struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalCopy initializes an OpGlobalCopy executor with the provided VM and Opcodes instance or returns an error.
func NewOpGlobalCopy(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalCopy{
		Opcode: op.Opcode(bytecode.OpGlobalCopy),
		vm:     vmT,
	}, nil
}

// Execute performs the operation to copy a global variable from sourceIndex to destIndex in the virtual machine.
func (op *OpGlobalCopy) Execute(decoder *core.Decoder) {
	sourceIndex := decoder.Read(1)
	destIndex := decoder.Read(0)
	value := op.vm.Globals().Get(uint(sourceIndex))
	op.vm.Globals().Set(uint(destIndex), value)
}
