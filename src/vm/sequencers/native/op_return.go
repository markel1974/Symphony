package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpReturn)
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.Opcode.
type OpReturn struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpReturn creates a new instance of OpReturn with its Opcode initialized for the OpReturn operation.
func NewOpReturn() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpReturn{
		opcode: opcodes.NewOpcode(OpReturnId, operands, "OpReturn"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpReturn) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpReturn) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the Core.
func (op *OpReturn) Execute(decoder *handler.Decoder) {
	var ret []objects.IObject
	if nRet := decoder.Operand(0); nRet > 0 {
		ret = make([]objects.IObject, nRet)
		for i := 0; i < nRet; i++ {
			ret[i] = op.vm.StackPop()
		}
	}
	op.vm.Return(ret)
}

// Compile generates the compiled representation of the OpReturn operation or returns an unimplemented error.
func (op *OpReturn) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
