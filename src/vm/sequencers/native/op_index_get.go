package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIndexGet)
}

// OpIndexGet represents the operation for performing an indexing operation on a value.
type OpIndexGet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIndexGet creates and returns a new instance of OpIndexGet initialized with its associated Opcode.
func NewOpIndexGet() handler.IOpExecutor {
	operands := _noOperands
	return &OpIndexGet{
		opcode: opcodes.NewOpcode(OpIndexGetId, operands, "OpIndexGet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIndexGet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Opcode returns the opcode associated with the instance.
func (op *OpIndexGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndexGet) Execute(_ *handler.Decoder) {
	index := op.vm.StackPop()
	left := op.vm.StackPop()
	val, err := left.IndexGet(op.vm.FrameId(), index)
	if err != nil {
		op.vm.Shutdown(objects.ComputeIndexGetError(err, index.TypeName(), index.TypeName()))
		return
	}
	if val == nil {
		val = op.vm.Factory().UndefinedValue()
	}
	op.vm.StackPush(val)
}

// Compile generates the compiled representation of the OpIndexGet operation or returns an unimplemented error.
func (op *OpIndexGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
