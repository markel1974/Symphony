package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers the OpIntLogical operation with the Sequencer system during package initialization.
func init() {
	SequencerRegister(NewOpIntArithmetic)
}

// OpIntArithmetic represents an arithmetic operation executor for integer-based operations within the virtual machine.
// It embeds bytecode.Opcode and utilizes the full access interface provided by core.IVMFullAccess.
type OpIntArithmetic struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIntArithmetic creates a new OpIntArithmetic executor for integer arithmetic operations in the virtual machine context.
// It requires a virtual machine implementing core.IVMFullAccess and associates the operation with the provided opcode.
// Returns an instance of core.IOpExecutor or an error if the virtual machine type assertion fails.
func NewOpIntArithmetic(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIntArithmetic{
		Opcode: op.Opcode(bytecode.OpIntArithmetic),
		vm:     vmT,
	}, nil
}

// Execute performs an integer arithmetic operation on the stack using the provided decoder to read operands and operator.
func (op *OpIntArithmetic) Execute(decoder *core.Decoder) {
	dstObj := op.vm.Stack().PeekAbsolute(decoder.Read(0))
	binaryOp := objects.ArithmeticOperator(decoder.Read(1))
	dst, ok := dstObj.(*objects.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	rhsObj := op.vm.Stack().Pop()
	rhs, ok := rhsObj.(*objects.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("rhs expected int, got %s", rhsObj.TypeName()))
		return
	}
	lhsObj := op.vm.Stack().Pop()
	lhs, ok := lhsObj.(*objects.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("lhs expected int, got %s", lhsObj.TypeName()))
		return
	}
	result, err := op.vm.Factory().ArithmeticOpInt64(binaryOp, lhs.Value(), rhs.Value())
	if err != nil {
		op.vm.SetError(err)
	}
	dst.SetValue(result)
}
