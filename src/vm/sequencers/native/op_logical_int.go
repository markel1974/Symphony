package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpLogicalInt function with the sequencer system by appending it to the internal registry.
func init() {
	SequencerRegister(NewOpLogicalInt)
}

// OpLogicalInt represents an executor for performing logical operations on integer operands within a virtual machine.
// It extends bytecode.Opcode to utilize its opcode properties and depends on the IVMFullAccess interface for Core interactions.
type OpLogicalInt struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLogicalInt creates a new instance of OpLogicalInt, validating the provided virtual machine and opcode inputs.
func NewOpLogicalInt() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint16, opcodes.SzUint8}
	return &OpLogicalInt{
		opcode: opcodes.NewOpcode(OpLogicalIntId, operands, "OpLogicalInt"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLogicalInt) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLogicalInt) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the logical operation between two integers on the stack and stores the result in the destination object.
func (op *OpLogicalInt) Execute(decoder *handler.Decoder) {
	logicalOp := objects.LogicalOperator(decoder.Operand(0))
	lhsIndex := decoder.Operand(1)
	rhsIndex := decoder.Operand(2)
	dstIndex := decoder.Operand(3)
	lhsObj := op.vm.StackPeekBP(uint(lhsIndex))
	rhsObj := op.vm.StackPeekBP(uint(rhsIndex))
	dstObj := op.vm.StackPeekBP(uint(dstIndex))
	result, err := op.vm.Factory().LogicalOpInt64(logicalOp, lhsObj.AsInt64(), rhsObj.AsInt64())
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	if err = op.vm.Factory().AssignBool(result, dstObj); err != nil {
		op.vm.Shutdown(err)
		return
	}
}

// Compile generates the compiled representation of the OpLogicalInt operation or returns an unimplemented error.
func (op *OpLogicalInt) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
