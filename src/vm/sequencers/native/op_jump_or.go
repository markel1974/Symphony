package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpJumpOr)
}

// OpJumpOr represents an operation that performs a logical OR and jumps based on the result.
type OpJumpOr struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpJumpOr creates and returns a new instance of OpJumpOr, associated with the OpJumpOr opcode and its details.
func NewOpJumpOr() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpOr{
		opcode: opcodes.NewOpcode(OpJumpOrId, operands, "OpJumpOr"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpOr) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpJumpOr) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpJumpOr) Execute(decoder *core.Decoder) {
	obj := op.vm.StackPeek()
	if obj.Falsy() {
		op.vm.StackDecrement()
	} else {
		pos := decoder.Operand(0)
		op.vm.SetIp(pos)
	}
}

// Compile generates the compiled representation of the OpJumpOr operation or returns an unimplemented error.
func (op *OpJumpOr) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
