package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpCallMethod)
}

// OpCallMethod represents a bytecode operation for invoking a method on an interface or object with dynamic dispatch.
type OpCallMethod struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpCallMethod creates and returns a new instance of OpCallMethod with initialized Opcode for the OpCallMethod opcode.
func NewOpCallMethod(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCallMethod{
		Opcode: op.Opcode(bytecode.OpCallMethod),
		vm:     vmT,
	}, nil
}

// Execute performs the dynamic dispatch logic.
func (op *OpCallMethod) Execute(decoder *core.Decoder) {
	// Operands Definition: Method name index (16-bit), Number of arguments (8-bit) -> [2, 1]
	// Decoder Logic (Reversed):
	// - decoder.Read(0) reads the LAST operand (numArgs, 1 byte)
	// - decoder.Read(1) reads the FIRST operand (methodNameIndex, 2 bytes)

	numArgs := decoder.Read(0)
	methodNameIndex := decoder.Read(1)

	// 1. Get method name from constants table.
	methodNameObj := op.vm.Constants().Get(uint(methodNameIndex))
	methodName, ok := methodNameObj.(*objects2.String)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid method name constant: not a string"))
		return
	}

	// 2. Get the interface object from stack. It's located below the arguments.
	interfaceObj := op.vm.Stack().PeekOffset(-1 - numArgs)
	io, ok := interfaceObj.(*objects2.Interface)
	if !ok {
		// If not an interface object, it could be a direct method call on a struct.
		// For now, we only handle the interface case. We could extend this.
		op.vm.SetError(fmt.Errorf("method call on non-interface object type: %s", interfaceObj.TypeName()))
		return
	}

	// 3. Perform Dynamic Dispatch: lookup method in 'ITable'.
	method, found := io.ITable()[methodName.Value()]
	if !found {
		op.vm.SetError(fmt.Errorf("undefined method '%s' for type '%s'", methodName.Value(), io.Value().TypeName()))
		return
	}

	callee, ok := method.(objects2.IObject)
	if !ok {
		op.vm.SetError(fmt.Errorf("method '%s' is not a callable function", methodName.Value()))
		return
	}

	// 4. Prepare stack for call.
	// Replace the interface object with its concrete value (the receiver).
	// Stack will now contain: [receiver, arg1, arg2, ..., Top of stack]
	op.vm.Stack().SetAbsolute(op.vm.Stack().StackPointer()-1-numArgs, io.Value())

	// 5. Delegate call logic to VM.
	// VM will handle creating new frame, etc.
	// Number of arguments for VM includes the receiver.
	op.vm.Call(callee, false, numArgs+1)
}
