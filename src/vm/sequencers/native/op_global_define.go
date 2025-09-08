package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

// init initializes the sequencer by registering the NewOpGlobalDefine operation with the SequencerRegister function.
func init() {
	SequencerRegister(NewOpGlobalDefine)
}

// OpGlobalDefine represents an operation to define a global variable in the VM environment.
// It binds a value from the stack to a global index specified by the decoder.
// Embeds `bytecode.Opcode` for opcode-related operations and uses `core.IVMFullAccess` for VM interactions.
type OpGlobalDefine struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalDefine creates a new OpGlobalDefine executor for the given virtual machine and opcode configuration.
// Returns an IOpExecutor instance or an error if the VM does not implement IVMFullAccess.
func NewOpGlobalDefine(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalDefine{
		Opcode: op.Opcode(bytecode.OpGlobalDefine),
		vm:     vmT,
	}, nil
}

// Execute sets a value from the stack into the global variables using the operand index from the decoder.
func (op *OpGlobalDefine) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	index := decoder.Read(0)
	val := op.vm.Stack().Peek()
	op.vm.Globals().Set(uint(index), val)
}
