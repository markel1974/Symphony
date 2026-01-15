package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpTrue)
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue() handler.IOpExecutor {
	operands := _noOperands
	return &OpTrue{
		opcode: opcodes.NewOpcode(OpTrueId, operands, "OpTrue"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpTrue) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpTrue) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(_ *handler.Decoder) {
	val := op.vm.Factory().TrueValue()
	op.vm.StackPush(val)
}

// Compile generates the compiled representation of the OpTrue operation or returns an unimplemented error.
func (op *OpTrue) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
