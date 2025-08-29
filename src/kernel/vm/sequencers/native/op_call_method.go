package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpCallMethod)
}

// OpCallMethod represents a bytecode operation for invoking a method on an interface or object with dynamic dispatch.
type OpCallMethod struct {
	*bytecode.Opcode
}

// NewOpCallMethod creates and returns a new instance of OpCallMethod with initialized Opcode for the OpCallMethod opcode.
func NewOpCallMethod(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpCallMethod{Opcode: op.Opcode(bytecode.OpCallMethod)}
}

// Execute performs the dynamic dispatch logic.
func (op *OpCallMethod) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Definition: Method name index (16-bit), Number of arguments (8-bit) -> [2, 1]
	// Decoder Logic (Reversed):
	// - decoder.Read(0) reads the LAST operand (numArgs, 1 byte)
	// - decoder.Read(1) reads the FIRST operand (methodNameIndex, 2 bytes)

	numArgs := decoder.Read(0)
	methodNameIndex := decoder.Read(1)

	// 1. Get method name from constants table.
	methodNameObj := v.Constants().Get(uint(methodNameIndex))
	methodName, ok := methodNameObj.(*objects.String)
	if !ok {
		v.SetError(fmt.Errorf("invalid method name constant: not a string"))
		return
	}

	// 2. Get the interface object from stack. It's located below the arguments.
	interfaceObj := v.Stack().PeekOffset(-1 - numArgs)
	io, ok := interfaceObj.(*objects.Interface)
	if !ok {
		// If not an interface object, it could be a direct method call on a struct.
		// For now, we only handle the interface case. We could extend this.
		v.SetError(fmt.Errorf("method call on non-interface object type: %s", interfaceObj.TypeName()))
		return
	}

	// 3. Perform Dynamic Dispatch: lookup method in 'ITable'.
	method, found := io.ITable()[methodName.Value()]
	if !found {
		v.SetError(fmt.Errorf("undefined method '%s' for type '%s'", methodName.Value(), io.Value().TypeName()))
		return
	}

	callee, ok := method.(objects.IObject)
	if !ok {
		v.SetError(fmt.Errorf("method '%s' is not a callable function", methodName.Value()))
		return
	}

	// 4. Prepare stack for call.
	// Replace the interface object with its concrete value (the receiver).
	// Stack will now contain: [receiver, arg1, arg2, ..., Top of stack]
	v.Stack().SetAbsolute(v.Stack().StackPointer()-1-numArgs, io.Value())

	// 5. Delegate call logic to VM.
	// VM will handle creating new frame, etc.
	// Number of arguments for VM includes the receiver.
	v.Call(callee, false, numArgs+1)
}
