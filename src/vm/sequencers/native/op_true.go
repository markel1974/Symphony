package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpTrue)
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue() core.IOpExecutor {
	operands := _noOperands
	return &OpTrue{
		opcode: opcodes.NewOpcode(OpTrueId, operands, "OpTrue"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpTrue) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(_ *core.Decoder) {
	// Operands Offset 0
	val := op.vm.Factory().TrueValue()
	op.vm.Stack().Push(val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpTrue) Opcode() *opcodes.Opcode {
	return op.opcode
}
