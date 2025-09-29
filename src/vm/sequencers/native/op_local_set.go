package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalSet)
}

// OpLocalSet represents an operation to set the value of a local variable within the current frame.
// It embeds Opcode for opcode-specific information such as name, operands, and code.
type OpLocalSet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLocalSet initializes and returns a new instance of OpLocalSet with associated opcode details.
func NewOpLocalSet() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalSet{
		opcode: opcodes.NewOpcode(OpLocalSetId, operands, "OpLocalSet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalSet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLocalSet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpLocalSet) Execute(decoder *handler.Decoder) {
	localIndex := decoder.Operand(0)
	val := op.vm.StackPeek()
	obj := op.vm.StackPeekBP(uint(localIndex))
	if freeObj, ok := obj.(*objects.ObjectPointer); ok {
		op.vm.Factory().SetPointer(freeObj, val)
	} else {
		op.vm.StackSetBP(uint(localIndex), val)
	}
}

// Compile generates the compiled representation of the OpLocalSet operation or returns an unimplemented error.
func (op *OpLocalSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
