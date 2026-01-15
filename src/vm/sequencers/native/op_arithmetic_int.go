package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the OpLogicalInt operation with the Sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpArithmeticInt)
}

// OpArithmeticInt represents an arithmetic operation executor for integer-based operations within the virtual machine.
// It embeds bytecode.Opcode and utilizes the full access interface provided by handler.IVMFullAccess.
type OpArithmeticInt struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpArithmeticInt creates a new OpArithmeticInt executor for integer arithmetic operations in the virtual machine context.
// It requires a virtual machine implementing handler.IVMFullAccess and associates the operation with the provided opcode.
// Returns an instance of handler.IOpExecutor or an error if the virtual machine type assertion fails.
func NewOpArithmeticInt() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint8}
	return &OpArithmeticInt{
		opcode: opcodes.NewOpcode(OpArithmeticIntId, operands, "OpArithmeticInt"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpArithmeticInt) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpArithmeticInt) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs an integer arithmetic operation on the stack using the provided decoder to read operands and operator.
func (op *OpArithmeticInt) Execute(decoder *handler.Decoder) {
	arithmeticOp := objects.ArithmeticOperator(decoder.Operand(0))
	lhs := decoder.Operand(1)
	rhs := decoder.Operand(2)
	dst := decoder.Operand(3)
	lhsObj := op.vm.StackPeekBP(uint(lhs))
	rhsObj := op.vm.StackPeekBP(uint(rhs))
	dstObj := op.vm.StackPeekBP(uint(dst))
	out, ok := dstObj.(*objects.Int)
	if !ok {
		op.vm.Shutdown(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	result, err := op.vm.Factory().ArithmeticOpInt64(arithmeticOp, lhsObj.AsInt64(), rhsObj.AsInt64())
	if err != nil {
		op.vm.Shutdown(err)
	}
	out.SetValue(result)
}

// Compile generates the compiled representation of the OpArithmeticInt operation or returns an unimplemented error.
func (op *OpArithmeticInt) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
