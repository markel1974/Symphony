package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpUnknown)
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding Opcode configuration set.
func NewOpUnknown() handler.IOpExecutor {
	operands := _noOperands
	return &OpUnknown{
		opcode: opcodes.NewOpcode(OpUnknownId, operands, "OpUnknown"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpUnknown) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpUnknown) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(_ *handler.Decoder) {
	op.vm.Shutdown(fmt.Errorf("unknown opcode at: %d", op.vm.GetIp()))
}

// Compile generates the compiled representation of the OpUnknown operation or returns an unimplemented error.
func (op *OpUnknown) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
