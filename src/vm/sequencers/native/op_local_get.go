package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalGet)
}

// OpLocalGet represents an operation to retrieve a local variable from the stack using its index.
type OpLocalGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpLocalGet creates a new OpLocalGet instance and initializes it with details for the OpLocalGet opcode.
func NewOpLocalGet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalGet{
		opcode: opcodes.NewOpcode(OpLocalGetId, operands, "OpLocalGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLocalGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpLocalGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	localIndex := decoder.Operand(0)
	val := op.vm.StackPeekOffsetBP(uint(localIndex))
	op.vm.StackPush(val)
}

// Compile generates the compiled representation of the OpLocalGet operation or returns an unimplemented error.
func (op *OpLocalGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
