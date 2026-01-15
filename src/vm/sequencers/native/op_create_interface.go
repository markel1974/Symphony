package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the NewOpCreateInterface function with the sequencer using SequencerRegister.
func init() {
	SequencerRegister(NewOpCreateInterface)
}

// OpCreateInterface represents an executor for handling interface-related bytecode operations in a virtual machine.
// It extends the bytecode.Opcode structure to access opcode details like id, operands, and name.
type OpCreateInterface struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCreateInterface creates a new instance of OpCreateInterface using the provided Opcodes instance.
// It returns an implementation of the handler.IOpExecutor interface for bytecode execution.
func NewOpCreateInterface() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpCreateInterface{
		opcode: opcodes.NewOpcode(OpCreateInterfaceId, operands, "OpCreateInterface"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateInterface) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpCreateInterface) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes an `OpCreateInterface` operation, constructing an interface by combining methods and a concrete value.
func (op *OpCreateInterface) Execute(decoder *handler.Decoder) {
	numMethods := decoder.Operand(0)
	interfaceObj := op.vm.StackPopInterface(numMethods)
	op.vm.StackPush(interfaceObj)
}

// Compile generates the compiled representation of the OpCreateInterface operation or returns an unimplemented error.
func (op *OpCreateInterface) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
