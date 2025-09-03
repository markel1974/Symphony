package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpConstant)
}

// OpConstant represents an operation used to load a constant onto the stack.
type OpConstant struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpConstant creates a new OpConstant instance with opcode details initialized for the OpConstant operation.
func NewOpConstant(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpConstant{
		Opcode: op.Opcode(bytecode.OpConstant),
		vm:     vmT,
	}, nil
}

// Execute executes the OpConstant instruction in the virtual machine, pushing a global constant onto the stack.
func (op *OpConstant) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	cIdx := decoder.Read(0)
	glObj := op.vm.Constants().Get(uint(cIdx))
	op.vm.Stack().Push(glObj)
}
