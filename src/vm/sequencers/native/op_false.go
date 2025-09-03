package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpFalse)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFalse{
		Opcode: op.Opcode(bytecode.OpFalse),
		vm:     vmT,
	}, nil
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(_ *core.Decoder) {
	// Operands Offset  0
	val := op.vm.Factory().FalseValue()
	op.vm.Stack().Push(val)
}
