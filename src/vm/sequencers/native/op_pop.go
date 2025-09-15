package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpPop)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpPop) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(_ *core.Decoder) {
	// Operands Offset 0
	op.vm.StackDecrement()
}

// Compile generates the compiled representation of the OpPop operation or returns an unimplemented error.
func (op *OpPop) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
