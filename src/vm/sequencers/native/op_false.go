package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFalse)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse() handler.IOpExecutor {
	operands := _noOperands
	return &OpFalse{
		opcode: opcodes.NewOpcode(OpFalseId, operands, "OpFalse"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpFalse) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpFalse) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(_ *handler.Decoder) {
	val := op.vm.Factory().FalseValue()
	op.vm.StackPush(val)
}

// Compile generates the compiled representation of the OpFalse operation or returns an unimplemented error.
func (op *OpFalse) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
