package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init registers the NewOpInterface function with the sequencer using SequencerRegister.
func init() {
	SequencerRegister(NewOpInterface)
}

// OpInterface represents an executor for handling interface-related bytecode operations in a virtual machine.
// It extends the bytecode.Opcode structure to access opcode details like id, operands, and name.
type OpInterface struct {
	*bytecode.Opcode
}

// NewOpInterface creates a new instance of OpInterface using the provided Opcodes instance.
// It returns an implementation of the core.IOpExecutor interface for bytecode execution.
func NewOpInterface(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpInterface{Opcode: op.Opcode(bytecode.OpInterface)}
}

// Execute processes an `OpInterface` operation, constructing an interface by combining methods and a concrete value.
func (op *OpInterface) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands: Number of methods (8-bit)
	numMethods := decoder.Read(0)

	// 1. Build the iTable by popping method names and functions from the stack.
	iTable := make(map[string]objects.IObject, numMethods)
	for i := 0; i < numMethods; i++ {
		// Pop in reverse order: function first, then name.
		methodFunc := v.Stack().Pop()
		methodNameObj := v.Stack().Pop()

		methodName, ok := methodNameObj.(*objects.String)
		if !ok {
			v.SetError(fmt.Errorf("interface method name must be a string, got %s", methodNameObj.TypeName()))
			return
		}
		iTable[methodName.Value()] = methodFunc
	}

	// 2. Pop the concrete value (the struct instance) that will be wrapped by the interface.
	concreteValue := v.Stack().Pop()

	// 3. Create the new interface object.
	interfaceObj := v.Factory().NewInterface(v.Frame().Id(), concreteValue, iTable)

	// 4. Push the final interface object back onto the stack.
	v.Stack().Push(interfaceObj)
}
