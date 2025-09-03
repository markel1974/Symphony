package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the NewOpInterface function with the sequencer using SequencerRegister.
func init() {
	SequencerRegister(NewOpInterface)
}

// OpInterface represents an executor for handling interface-related bytecode operations in a virtual machine.
// It extends the bytecode.Opcode structure to access opcode details like id, operands, and name.
type OpInterface struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpInterface creates a new instance of OpInterface using the provided Opcodes instance.
// It returns an implementation of the core.IOpExecutor interface for bytecode execution.
func NewOpInterface(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpInterface{
		Opcode: op.Opcode(bytecode.OpInterface),
		vm:     vmT,
	}, nil
}

// Execute processes an `OpInterface` operation, constructing an interface by combining methods and a concrete value.
func (op *OpInterface) Execute(decoder *core.Decoder) {
	// Operands: Number of methods (8-bit)
	numMethods := decoder.Read(0)

	// 1. Build the iTable by popping method names and functions from the stack.
	iTable := make(map[string]objects2.IObject, numMethods)
	for i := 0; i < numMethods; i++ {
		// Pop in reverse order: function first, then name.
		methodFunc := op.vm.Stack().Pop()
		methodNameObj := op.vm.Stack().Pop()

		methodName, ok := methodNameObj.(*objects2.String)
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
