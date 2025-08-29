package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIntOp)
}

// OpIntOp extends Opcode and represents integer operations performed on a virtual machine.
type OpIntOp struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIntOp initializes and returns a new instance of OpIntOp with relevant opcode details provided by bytecode.Opcodes.
func NewOpIntOp(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIntOp{
		Opcode: op.Opcode(bytecode.OpIntOp),
		vm:     vmT,
	}, nil
}

// Execute performs a specified binary operation between two integers from the stack and stores the result in a destination slot.
// It retrieves operands, validates types, and executes the operation, setting an error on unsupported cases or type mismatches.
func (op *OpIntOp) Execute(decoder *core.Decoder) {
	dstObj := op.vm.Stack().PeekAbsolute(decoder.Read(0))
	binaryOp := objects.Operator(decoder.Read(1))
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
	result, err := op.vm.Factory().BinaryOpInt64(binaryOp, lhs.Value(), rhs.Value())
	if err != nil {
		op.vm.SetError(err)
	}
	dst.SetValue(result)
}
