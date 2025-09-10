package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJumpFalsy)
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding Opcode.
func NewOpJumpFalsy() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpFalsy{
		opcode: opcodes.NewOpcode(OpJumpFalsyId, operands, "OpJumpFalsy"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpJumpFalsy) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := op.vm.Stack().Pop()
	if obj.Falsy() {
		pos := decoder.Operand(0)
		op.vm.SetIp(pos - 1)
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpFalsy) Opcode() *opcodes.Opcode {
	return op.opcode
}
