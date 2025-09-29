package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the necessary operation for the sequencer by registering NewOpJumpTruthy with the SequencerRegister function.
func init() {
	SequencerRegister(NewOpJumpTruthy)
}

// OpJumpTruthy defines a structure for executing a truthy-check jump operation in a virtual machine context.
type OpJumpTruthy struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpJumpTruthy creates a new OpJumpTruthy executor, validating the provided virtual machine for required full access capabilities.
func NewOpJumpTruthy() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpJumpTruthy{
		opcode: opcodes.NewOpcode(OpJumpTruthyId, operands, "OpJumpTruthy"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpTruthy) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpJumpTruthy) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute evaluates the top object on the stack and updates the instruction pointer if it is truthy.
func (op *OpJumpTruthy) Execute(decoder *handler.Decoder) {
	obj := op.vm.StackPop()
	if !obj.Falsy() {
		pos := decoder.Operand(0)
		op.vm.SetIp(uint(pos))
	}
}

// Compile generates the compiled representation of the OpJumpTruthy operation or returns an unimplemented error.
func (op *OpJumpTruthy) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
