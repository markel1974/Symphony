package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateMap)
}

// OpCreateMap is a wrapper around bytecode.Opcode, representing a map creation operation in bytecode execution.
type OpCreateMap struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCreateMap initializes and returns a new instance of OpCreateMap with its Opcode set to OpCreateMap details.
func NewOpCreateMap() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpCreateMap{
		opcode: opcodes.NewOpcode(OpCreateMapId, operands, "OpCreateMap"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateMap) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpCreateMap) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCreateMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpCreateMap) Execute(decoder *handler.Decoder) {
	numElements := decoder.Operand(0)
	mObj := op.vm.StackPopMap(uint(numElements))
	op.vm.StackPush(mObj)
}

// Compile generates the compiled representation of the OpCreateMap operation or returns an unimplemented error.
func (op *OpCreateMap) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
