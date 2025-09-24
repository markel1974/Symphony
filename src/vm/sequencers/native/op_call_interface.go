package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCallInterface)
}

// OpCallInterface represents a bytecode operation for invoking a method on an interface or object with dynamic dispatch.
type OpCallInterface struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCallInterface creates and returns a new instance of OpCallInterface with initialized Opcode for the OpCallInterface opcode.
func NewOpCallInterface() core.IOpExecutor {
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

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCallInterface) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the dynamic dispatch logic.
func (op *OpCallInterface) Execute(decoder *core.Decoder) {
	methodNameIndex := decoder.Operand(0)
	spread := decoder.Operand(1)
	numArgs := decoder.Operand(2)
	methodNameObj := op.vm.ConstantsGet(uint(methodNameIndex))
	offset := numArgs + 1
	interfaceObj := op.vm.StackPeekSP(uint(offset))
	io, ok := interfaceObj.(*objects.Interface)
	if !ok {
		op.vm.SetError(fmt.Errorf("method call on non-interface object type: %s", interfaceObj.TypeName()))
		return
	}
	callee, ok := io.Method(methodNameObj.AsString())
	if !ok {
		op.vm.SetError(fmt.Errorf("undefined method '%s' for type '%s'", methodNameObj.AsString(), io.Value().TypeName()))
		return
	}
	op.vm.StackSetSP(uint(offset), io.Value())
	hasSpread := spread > 0
	totalArgs := numArgs + 1 //receiver
	op.vm.Call(callee, hasSpread, totalArgs)
}

// Compile generates the compiled representation of the OpCallInterface operation or returns an unimplemented error.
func (op *OpCallInterface) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
