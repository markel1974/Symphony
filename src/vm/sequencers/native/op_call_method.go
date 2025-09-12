package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	objects "github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCallMethod)
}

// OpCallMethod represents a bytecode operation for invoking a method on an interface or object with dynamic dispatch.
type OpCallMethod struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCallMethod creates and returns a new instance of OpCallMethod with initialized Opcode for the OpCallMethod opcode.
func NewOpCallMethod() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16, opcodes.SzUint8}
	return &OpCallMethod{
		opcode: opcodes.NewOpcode(OpCallMethodId, operands, "OpCallMethod"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCallMethod) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the dynamic dispatch logic.
func (op *OpCallMethod) Execute(decoder *core.Decoder) {
	numArgs := decoder.Operand(0)
	methodNameIndex := decoder.Operand(1)
	methodNameObj := op.vm.Constants().Get(uint(methodNameIndex))
	methodName, ok := methodNameObj.(*objects.String)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid method name constant: not a string"))
		return
	}
	offset := numArgs + 1
	interfaceObj := op.vm.StackPeekOffsetSP(uint(offset))
	io, ok := interfaceObj.(*objects.Interface)
	if !ok {
		op.vm.SetError(fmt.Errorf("method call on non-interface object type: %s", interfaceObj.TypeName()))
		return
	}
	method, found := io.ITable()[methodName.Value()]
	if !found {
		op.vm.SetError(fmt.Errorf("undefined method '%s' for type '%s'", methodName.Value(), io.Value().TypeName()))
		return
	}
	callee, ok := method.(objects.IObject)
	if !ok {
		op.vm.SetError(fmt.Errorf("method '%s' is not a callable function", methodName.Value()))
		return
	}
	target := numArgs + 1
	op.vm.StackSetOffsetSP(uint(target), io.Value())
	//op.vm.Stack().SetAbsolute(op.vm.Stack().StackPointer()-1-numArgs, io.Value())
	op.vm.Call(callee, false, numArgs+1)
}

// Opcode returns the opcode associated with the instance.
func (op *OpCallMethod) Opcode() *opcodes.Opcode {
	return op.opcode
}
