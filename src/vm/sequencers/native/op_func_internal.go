package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

// init initializes the sequencer system by registering the internal function operation for bytecode processing.
func init() {
	SequencerRegister(NewOpFuncInternal)
}

// OpFuncInternal represents an operation handler for internal functions within the virtual machine environment.
// It embeds a bytecode.Opcode instance and provides access to the virtual machine through a core.IVMFullAccess interface.
type OpFuncInternal struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpFuncInternal creates an internal operation executor for the given virtual machine and opcode configuration.
// Returns an implementation of core.IOpExecutor or an error if the VM doesn't support full access.
func NewOpFuncInternal(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFuncInternal{
		Opcode: op.Opcode(bytecode.OpFuncInternal),
		vm:     vmT,
	}, nil
}

// Execute executes the operation defined by OpFuncInternal using the provided decoder to process operands and VM context.
func (op *OpFuncInternal) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	nameIndex := decoder.Read(0)
	symbol := op.vm.Internals().Get(uint(nameIndex))
	op.vm.Stack().Push(symbol)
}
