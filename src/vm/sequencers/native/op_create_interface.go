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
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateInterface creates a new instance of OpCreateInterface using the provided Opcodes instance.
// It returns an implementation of the core.IOpExecutor interface for bytecode execution.
func NewOpCreateInterface(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateInterface{
		Opcode: op.Opcode(opcodes.OpCreateInterface),
		vm:     vmT,
	}, nil
}

// Execute processes an `OpCreateInterface` operation, constructing an interface by combining methods and a concrete value.
func (op *OpCreateInterface) Execute(decoder *core.Decoder) {
	// Operands: Number of methods (8-bit)
	numMethods := decoder.Read(0)

	// 1. Build the iTable by popping method names and functions from the stack.
	iTable := make(map[string]objects.IObject, numMethods)
	for i := 0; i < numMethods; i++ {
		// Pop in reverse order: function first, then name.
		methodFunc := op.vm.Stack().Pop()
		methodNameObj := op.vm.Stack().Pop()
		methodName, ok := methodNameObj.(*objects.String)
		if !ok {
			op.vm.SetError(fmt.Errorf("interface method name must be a string, got %s", methodNameObj.TypeName()))
			return
		}
		iTable[methodName.Value()] = methodFunc
	}

	// 2. Pop the concrete value (the struct instance) that will be wrapped by the interface.
	concreteValue := op.vm.Stack().Pop()

	// 3. Create the new interface object.
	interfaceObj := op.vm.Factory().NewInterface(op.vm.Frame().Id(), concreteValue, iTable)

	// 4. Push the final interface object back onto the stack.
	op.vm.Stack().Push(interfaceObj)
}
