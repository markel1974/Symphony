package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the NewOpIntLogical function with the sequencer system by appending it to the internal registry.
func init() {
	SequencerRegister(NewOpIntLogical)
}

// OpIntLogical represents an executor for performing logical operations on integer operands within a virtual machine.
// It extends bytecode.Opcode to utilize its opcode properties and depends on the IVMFullAccess interface for VM interactions.
type OpIntLogical struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpIntLogical creates a new instance of OpIntLogical, validating the provided virtual machine and opcode inputs.
func NewOpIntLogical(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIntLogical{
		Opcode: op.Opcode(bytecode.OpIntLogical),
		vm:     vmT,
	}, nil
}

// Execute performs the logical operation between two integers on the stack and stores the result in the destination object.
func (op *OpIntLogical) Execute(decoder *core.Decoder) {
	dstObj := op.vm.Stack().PeekAbsolute(decoder.Read(0))
	binaryOp := objects2.LogicalOperator(decoder.Read(1))
	dst, ok := dstObj.(*objects2.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	rhsObj := op.vm.Stack().Pop()
	rhs, ok := rhsObj.(*objects2.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("rhs expected int, got %s", rhsObj.TypeName()))
		return
	}
	lhsObj := op.vm.Stack().Pop()
	lhs, ok := lhsObj.(*objects2.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("lhs expected int, got %s", lhsObj.TypeName()))
		return
	}
	result, err := op.vm.Factory().LogicalOpInt64(binaryOp, lhs.Value(), rhs.Value())
	if err != nil {
		op.vm.SetError(err)
	}
	if result {
		dst.SetValue(1)
	} else {
		dst.SetValue(0)
	}
}
