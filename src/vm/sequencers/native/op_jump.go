package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJump)
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJump{
		opcode: opcodes.NewOpcode(OpJumpId, operands, "OpJump"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJump) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpJump) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`Core`) to a calculated position in the frame.
func (op *OpJump) Execute(decoder *handler.Decoder) {
	pos := decoder.Operand(0)
	op.vm.SetIp(uint(pos))
}

// Compile generates the compiled representation of the OpJump operation or returns an unimplemented error.
func (op *OpJump) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
