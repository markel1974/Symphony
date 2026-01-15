package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalDefine)
}

// OpLocalDefine represents the opcode for defining a new local variable within the current frame's scope.
type OpLocalDefine struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLocalDefine creates a new instance of OpLocalDefine with its associated opcode details.
func NewOpLocalDefine() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalDefine{
		opcode: opcodes.NewOpcode(OpLocalDefineId, operands, "OpLocalDefine"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalDefine) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLocalDefine) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpLocalDefine) Execute(decoder *handler.Decoder) {
	localIndex := decoder.Operand(0)
	val := op.vm.StackPeek()
	op.vm.StackSetBP(uint(localIndex), val)
}

// Compile generates the compiled representation of the OpLocalDefine operation or returns an unimplemented error.
func (op *OpLocalDefine) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
