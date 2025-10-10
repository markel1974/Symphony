// File: vm/sequencers/native/op_type_assert.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the system by registering the NewOpTypeAssert operation with the sequencer via SequencerRegister.
func init() {
	SequencerRegister(NewOpTypeAssert)
}

// OpTypeAssert represents a bytecode operation for performing type assertions in a virtual machine.
// It embeds bytecode.Opcode to utilize opcode-related functionalities and operates on a handler.IVMFullAccess instance.
type OpTypeAssert struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpTypeAssert creates a new instance of OpTypeAssert, ensuring the provided IVM implements IVMFullAccess.
func NewOpTypeAssert() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.Relocatable}
	return &OpTypeAssert{
		opcode: opcodes.NewOpcode(OpTypeAssertId, operands, "OpTypeAssert"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpTypeAssert) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpTypeAssert) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the type assertion operation. It attempts to assert if the top stack object matches the desired type.
// It decodes the target type from the constant pool using the provided decoder and validates against the interface object.
// On success, the concrete value and a boolean 'true' are pushed onto the stack. On failure, undefined and 'false' are pushed.
func (op *OpTypeAssert) Execute(decoder *handler.Decoder) {
	typeNameIndex := decoder.Operand(0)
	hasOk := decoder.Operand(1)
	interfaceObj := op.vm.StackPop()
	valid := false
	concreteValue := op.vm.Factory().Concrete(op.vm.FrameId(), interfaceObj)
	switch concreteValue.(type) {
	case *objects.Any:
		valid = true
	default:
		targetTypeObj, err := op.vm.ConstantsGet(uint(typeNameIndex))
		if err != nil {
			op.vm.Shutdown(err)
			return
		}
		if concreteValue.TypeName() == targetTypeObj.AsString() {
			valid = true
		}
	}

	if valid {
		op.vm.StackPush(concreteValue)
		if hasOk != 0 {
			op.vm.StackPush(op.vm.Factory().TrueValue())
		}
	} else {
		op.vm.StackPush(op.vm.Factory().UndefinedValue())
		if hasOk != 0 {
			op.vm.StackPush(op.vm.Factory().FalseValue())
		}
	}
}

// Compile generates the compiled representation of the OpTypeAssert operation or returns an unimplemented error.
func (op *OpTypeAssert) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
