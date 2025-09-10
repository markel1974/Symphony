package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpIntLogical function with the sequencer system by appending it to the internal registry.
func init() {
	SequencerRegister(NewOpIntLogical)
}

// OpIntLogical represents an executor for performing logical operations on integer operands within a virtual machine.
// It extends bytecode.Opcode to utilize its opcode properties and depends on the IVMFullAccess interface for VM interactions.
type OpIntLogical struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpIntLogical creates a new instance of OpIntLogical, validating the provided virtual machine and opcode inputs.
func NewOpIntLogical(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpIntLogical{
		Opcode: op.Opcode(opcodes.OpIntLogical),
		vm:     vmT,
	}, nil
}

// Execute performs the logical operation between two integers on the stack and stores the result in the destination object.
func (op *OpIntLogical) Execute(decoder *core.Decoder) {
	logicalOp := objects.LogicalOperator(decoder.Read(0))
	lhsObj := op.vm.Stack().PeekAbsolute(decoder.Read(1))
	rhsObj := op.vm.Stack().PeekAbsolute(decoder.Read(2))
	dstObj := op.vm.Stack().PeekAbsolute(decoder.Read(3))
	dst, ok := dstObj.(*objects.Int)
	if !ok {
		op.vm.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	result, err := op.vm.Factory().LogicalOpInt64(logicalOp, lhsObj.AsInt64(), rhsObj.AsInt64())
	if err != nil {
		op.vm.SetError(err)
		return
	}
	if result {
		dst.SetValue(1)
	} else {
		dst.SetValue(0)
	}
}
