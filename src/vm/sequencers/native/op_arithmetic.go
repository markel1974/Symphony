package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the NewOpLogical operation with the sequencer system by adding it to the _registerContainer.
func init() {
	SequencerRegister(NewOpArithmetic)
}

// OpArithmetic represents an operation responsible for performing arithmetic computations within the virtual machine.
type OpArithmetic struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpArithmetic creates a new arithmetic operation executor with the given virtual machine and opcode.
// It ensures the virtual machine implements the IVMFullAccess interface before creating the executor.
// Returns an IOpExecutor for arithmetic operations or an error if the provided Core does not support full access.
func NewOpArithmetic() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpArithmetic{
		opcode: opcodes.NewOpcode(OpArithmeticId, operands, "OpArithmetic"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpArithmetic) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the OpArithmetic instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpArithmetic) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the arithmetic operation by decoding the operator and applying it to the top two stack elements.
// It retrieves the operator from the bytecode, performs the calculation, and pushes the result back onto the stack.
// In case of an operation error, it sets the virtual machine's error state.
func (op *OpArithmetic) Execute(decoder *handler.Decoder) {
	opcode := decoder.Operand(0)
	right := op.vm.StackPop()
	left := op.vm.StackPop()
	operator := objects.ArithmeticOperator(opcode)
	res, err := left.ArithmeticOp(op.vm.FrameId(), operator, right)
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	op.vm.StackPush(res)
}

// Compile generates the compiled representation of the OpArithmetic operation or returns an unimplemented error.
func (op *OpArithmetic) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
