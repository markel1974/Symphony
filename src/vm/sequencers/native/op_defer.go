package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpDefer)
}

type OpDefer struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

func NewOpDefer() handler.IOpExecutor {
	operands := _noOperands
	return &OpDefer{
		opcode: opcodes.NewOpcode(OpDeferId, operands, "OpDefer"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpDefer) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpDefer) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the top stack item, deferring its execution if it is a compiled function, otherwise sets an error.
func (op *OpDefer) Execute(_ *handler.Decoder) {
	obj := op.vm.StackPop()
	op.vm.FrameDeferredAdd(obj)
}

// Compile generates the compiled representation of the OpDefer operation or returns an unimplemented error.
func (op *OpDefer) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
