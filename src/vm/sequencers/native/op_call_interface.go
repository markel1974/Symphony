package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCallInterface)
}

// OpCallInterface represents a bytecode operation for invoking a method on an interface or object with dynamic dispatch.
type OpCallInterface struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCallInterface creates and returns a new instance of OpCallInterface with initialized Opcode for the OpCallInterface opcode.
func NewOpCallInterface() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.SzUint8, opcodes.Relocatable}
	return &OpCallInterface{
		opcode: opcodes.NewOpcode(OpCallInterfaceId, operands, "OpCallInterface"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCallInterface) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpCallInterface) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the dynamic dispatch logic.
func (op *OpCallInterface) Execute(decoder *handler.Decoder) {
	methodNameIndex := decoder.Operand(0)
	spread := decoder.Operand(1)
	numArgs := decoder.Operand(2)
	methodNameObj, err := op.vm.ConstantsGet(uint(methodNameIndex))
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	offset := numArgs + 1
	interfaceObj := op.vm.StackPeekSP(uint(offset))
	var callee objects.IObject
	var value objects.IObject
	var ok bool

	switch io := interfaceObj.(type) {
	case *objects.Interface:
		if callee, ok = io.Method(methodNameObj.AsString()); !ok {
			op.vm.Shutdown(fmt.Errorf("undefined method '%s' for type '%s'", methodNameObj.AsString(), io.Value().TypeName()))
			return
		}
		value = io.Value()
	case *objects.Any:
		if callee, ok = io.Method(methodNameObj.AsString()); !ok {
			op.vm.Shutdown(fmt.Errorf("undefined method '%s' for any object", methodNameObj.AsString()))
			return
		}
		value = io
	}
	op.vm.StackSetSP(uint(offset), value)
	hasSpread := spread > 0
	totalArgs := numArgs + 1 //receiver
	op.vm.Call(callee, false, hasSpread, totalArgs)
}

// Compile generates the compiled representation of the OpCallInterface operation or returns an unimplemented error.
func (op *OpCallInterface) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
