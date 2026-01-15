package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJumpFalsy)
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding Opcode.
func NewOpJumpFalsy() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpFalsy{
		opcode: opcodes.NewOpcode(OpJumpFalsyId, operands, "OpJumpFalsy"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpFalsy) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpJumpFalsy) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(decoder *handler.Decoder) {
	obj := op.vm.StackPop()
	if obj.Falsy() {
		pos := decoder.Operand(0)
		op.vm.SetIp(uint(pos))
	}
}

// Compile generates the compiled representation of the OpJumpFalsy operation or returns an unimplemented error.
func (op *OpJumpFalsy) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
