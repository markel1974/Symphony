package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes and registers the NewOpBinary operation with the sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpLogical)
}

// OpLogical represents a logical operation bytecode execution handler within the virtual machine context.
type OpLogical struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpLogical creates a new instance of OpLogical executor for logical bytecode operations using the provided VM and opcode.
// It returns an IOpExecutor implementation or an error if the VM does not support IVMFullAccess.
func NewOpLogical() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpLogical{
		opcode: opcodes.NewOpcode(OpLogicalId, operands, "OpLogical"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLogical) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the logical operation by decoding the opcode, applying the binary operation, and updating the stack.
func (op *OpLogical) Execute(decoder *core.Decoder) {
	// Operands Offset  1 (8 bits)
	opcode := decoder.Operand(0)
	right := op.vm.StackPop()
	left := op.vm.StackPop()
	operator := objects.LogicalOperator(opcode)
	res, err := left.LogicalOp(op.vm.FrameId(), operator, right)
	if err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.StackPush(res)
}

// Opcode returns the opcode associated with the instance.
func (op *OpLogical) Opcode() *opcodes.Opcode {
	return op.opcode
}
