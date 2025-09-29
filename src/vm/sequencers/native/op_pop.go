package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpPop)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop() handler.IOpExecutor {
	operands := _noOperands
	return &OpPop{
		opcode: opcodes.NewOpcode(OpPopId, operands, "OpPop"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpPop) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpPop) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the Core.
func (op *OpPop) Execute(_ *handler.Decoder) {
	op.vm.StackDecrement()
}

// Compile generates the compiled representation of the OpPop operation or returns an unimplemented error.
func (op *OpPop) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
