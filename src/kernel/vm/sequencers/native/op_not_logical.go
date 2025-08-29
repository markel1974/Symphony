package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpNotLogical)
}

// OpNotLogical represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpNotLogical struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpNotLogical creates a new instance of OpNotLogical, representing a logical NOT operation (!).
func NewOpNotLogical(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpNotLogical{
		Opcode: op.Opcode(bytecode.OpNotLogical),
		vm:     vmT,
	}, nil
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpNotLogical) Execute(_ *core.Decoder) {
	// Operands Offset  0
	operand := op.vm.Stack().Pop()
	var val objects.IObject
	if operand.Falsy() {
		val = op.vm.Factory().TrueValue()
	} else {
		val = op.vm.Factory().FalseValue()
	}
	op.vm.Stack().Push(val)
}
