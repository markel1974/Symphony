package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJumpAnd)
}

// OpJumpAnd represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpJumpAnd struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpJumpAnd creates and returns a new instance of OpJumpAnd, initializing it with details for the OpJumpAnd opcode.
func NewOpJumpAnd() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpAnd{
		opcode: opcodes.NewOpcode(OpJumpAndId, operands, "OpJumpAnd"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpAnd) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpJumpAnd) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpJumpAnd) Execute(decoder *core.Decoder) {
	obj := op.vm.StackPeek()
	if obj.Falsy() {
		pos := decoder.Operand(0)
		op.vm.SetIp(pos)
	} else {
		op.vm.StackDecrement()
	}
}

// Compile generates the compiled representation of the OpJumpAnd operation or returns an unimplemented error.
func (op *OpJumpAnd) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
