package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init initializes and registers the NewOpBinary operation with the sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpLogical)
}

// OpLogical represents a logical operation bytecode execution handler within the virtual machine context.
type OpLogical struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLogical creates a new instance of OpLogical executor for logical bytecode operations using the provided Core and opcode.
// It returns an IOpExecutor implementation or an error if the Core does not support IVMFullAccess.
func NewOpLogical() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpLogical{
		opcode: opcodes.NewOpcode(OpLogicalId, operands, "OpLogical"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLogical) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLogical) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the logical operation by decoding the opcode, applying the binary operation, and updating the stack.
func (op *OpLogical) Execute(decoder *handler.Decoder) {
	opcode := decoder.Operand(0)
	rhs := op.vm.StackPop()
	lhs := op.vm.StackPop()
	operator := objects.LogicalOperator(opcode)
	res, err := lhs.LogicalOp(op.vm.FrameId(), operator, rhs)
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	op.vm.StackPush(res)
}

// Compile generates the compiled representation of the OpLogical operation or returns an unimplemented error.
func (op *OpLogical) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
