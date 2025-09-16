package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the OpIntLogical operation with the Sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpIntArithmetic)
}

// OpIntArithmetic represents an arithmetic operation executor for integer-based operations within the virtual machine.
// It embeds bytecode.Opcode and utilizes the full access interface provided by core.IVMFullAccess.
type OpIntArithmetic struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIntArithmetic creates a new OpIntArithmetic executor for integer arithmetic operations in the virtual machine context.
// It requires a virtual machine implementing core.IVMFullAccess and associates the operation with the provided opcode.
// Returns an instance of core.IOpExecutor or an error if the virtual machine type assertion fails.
func NewOpIntArithmetic() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint8}
	return &OpIntArithmetic{
		opcode: opcodes.NewOpcode(OpIntArithmeticId, operands, "OpIntArithmetic"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIntArithmetic) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIntArithmetic) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs an integer arithmetic operation on the stack using the provided decoder to read operands and operator.
func (op *OpIntArithmetic) Execute(decoder *core.Decoder) {
	arithmeticOp := objects.ArithmeticOperator(decoder.Operand(0))
	lhs := decoder.Operand(1)
	rhs := decoder.Operand(2)
	dst := decoder.Operand(3)
	lhsObj := op.vm.StackPeekBP(uint(lhs))
	rhsObj := op.vm.StackPeekBP(uint(rhs))
	dstObj := op.vm.StackPeekBP(uint(dst))
	out, ok := dstObj.(*objects.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	result, err := op.vm.Factory().ArithmeticOpInt64(arithmeticOp, lhsObj.AsInt64(), rhsObj.AsInt64())
	if err != nil {
		op.vm.SetError(err)
	}
	out.SetValue(result)
}

// Compile generates the compiled representation of the OpIntArithmetic operation or returns an unimplemented error.
func (op *OpIntArithmetic) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
