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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIntLogical creates a new instance of OpIntLogical, validating the provided virtual machine and opcode inputs.
func NewOpIntLogical() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint8}
	return &OpIntLogical{
		opcode: opcodes.NewOpcode(OpIntLogicalId, operands, "OpIntLogical"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIntLogical) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the logical operation between two integers on the stack and stores the result in the destination object.
func (op *OpIntLogical) Execute(decoder *core.Decoder) {
	logicalOp := objects.LogicalOperator(decoder.Operand(0))
	lhs := decoder.Operand(1)
	rhs := decoder.Operand(2)
	dst := decoder.Operand(3)
	lhsObj := op.vm.StackPeekOffsetBP(uint(lhs))
	rhsObj := op.vm.StackPeekOffsetBP(uint(rhs))
	dstObj := op.vm.StackPeekOffsetBP(uint(dst))
	out, ok := dstObj.(*objects.Int)
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
		out.SetValue(1)
	} else {
		out.SetValue(0)
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIntLogical) Opcode() *opcodes.Opcode {
	return op.opcode
}
