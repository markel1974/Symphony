package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateArray)
}

// OpCreateArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base Opcode for opcode, operands, and name information.
type OpCreateArray struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCreateArray creates and returns a new instance of OpCreateArray, initialized with details for the OpCreateArray operation.
func NewOpCreateArray() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpCreateArray{
		opcode: opcodes.NewOpcode(OpCreateArrayId, operands, "OpCreateArray"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateArray) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpCreateArray) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCreateArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpCreateArray) Execute(decoder *handler.Decoder) {
	numElements := decoder.Operand(0)
	arrObj := op.vm.StackPopArray(uint(numElements))
	op.vm.StackPush(arrObj)
}

// Compile generates the compiled representation of the OpCreateArray operation or returns an unimplemented error.
func (op *OpCreateArray) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
