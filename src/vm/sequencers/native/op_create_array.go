package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
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
	vm     core.IVMFullAccess
}

// NewOpCreateArray creates and returns a new instance of OpCreateArray, initialized with details for the OpCreateArray operation.
func NewOpCreateArray() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateArray) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCreateArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpCreateArray) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Operand(0)
	elements := op.vm.StackPopArray(uint(numElements))
	arr := op.vm.Factory().NewArray(op.vm.FrameId(), elements)
	op.vm.StackPush(arr)
}

// Compile generates the compiled representation of the OpCreateArray operation or returns an unimplemented error.
func (op *OpCreateArray) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
