package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init initializes and registers the NewOpBinary operation with the sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpLogical)
}

// OpLogical represents a logical operation bytecode execution handler within the virtual machine context.
type OpLogical struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpLogical creates a new instance of OpLogical executor for logical bytecode operations using the provided VM and opcode.
// It returns an IOpExecutor implementation or an error if the VM does not support IVMFullAccess.
func NewOpLogical(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpLogical{
		Opcode: op.Opcode(bytecode.OpLogical),
		vm:     vmT,
	}, nil
}

// Execute processes the logical operation by decoding the opcode, applying the binary operation, and updating the stack.
func (op *OpLogical) Execute(decoder *core.Decoder) {
	// Operands Offset  1 (8 bits)
	opcode := decoder.Read(0)
	right := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	operator := objects.LogicalOperator(opcode)
	res, err := left.LogicalOp(op.vm.Frame().Id(), operator, right)
	if err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.Stack().Push(res)
}
