package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateMap)
}

// OpCreateMap is a wrapper around bytecode.Opcode, representing a map creation operation in bytecode execution.
type OpCreateMap struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCreateMap initializes and returns a new instance of OpCreateMap with its Opcode set to OpCreateMap details.
func NewOpCreateMap() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateMap) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCreateMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpCreateMap) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	numElements := decoder.Operand(0)
	mElem := op.vm.StackPopMap(uint(numElements))
	op.vm.StackPush(op.vm.Factory().NewMap(op.vm.FrameId(), mElem))
}

// Compile generates the compiled representation of the OpCreateMap operation or returns an unimplemented error.
func (op *OpCreateMap) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
