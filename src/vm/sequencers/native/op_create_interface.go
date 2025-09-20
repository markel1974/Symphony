package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpCreateInterface function with the sequencer using SequencerRegister.
func init() {
	SequencerRegister(NewOpCreateInterface)
}

// OpCreateInterface represents an executor for handling interface-related bytecode operations in a virtual machine.
// It extends the bytecode.Opcode structure to access opcode details like id, operands, and name.
type OpCreateInterface struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCreateInterface creates a new instance of OpCreateInterface using the provided Opcodes instance.
// It returns an implementation of the core.IOpExecutor interface for bytecode execution.
func NewOpCreateInterface() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateInterface) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes an `OpCreateInterface` operation, constructing an interface by combining methods and a concrete value.
func (op *OpCreateInterface) Execute(decoder *core.Decoder) {
	numMethods := decoder.Operand(0)
	iTable := make(map[string]objects.IObject, numMethods)
	for i := 0; i < numMethods; i++ {
		methodFunc := op.vm.StackPop()
		methodNameObj := op.vm.StackPop()
		iTable[methodNameObj.AsString()] = methodFunc
	}
	concreteValue := op.vm.StackPop()
	interfaceObj := op.vm.Factory().NewInterface(op.vm.FrameId(), concreteValue, iTable)
	op.vm.StackPush(interfaceObj)
}

// Compile generates the compiled representation of the OpCreateInterface operation or returns an unimplemented error.
func (op *OpCreateInterface) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
