package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJump)
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJump{
		opcode: opcodes.NewOpcode(OpJumpId, operands, "OpJump"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpJump) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`VM`) to a calculated position in the frame.
func (op *OpJump) Execute(decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	pos := decoder.Operand(0)
	op.vm.SetIp(pos - 1)
}

// Opcode returns the opcode associated with the instance.
func (op *OpJump) Opcode() *opcodes.Opcode {
	return op.opcode
}
