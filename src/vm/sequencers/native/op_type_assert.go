// File: vm/sequencers/native/op_type_assert.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the system by registering the NewOpTypeAssert operation with the sequencer via SequencerRegister.
func init() {
	SequencerRegister(NewOpTypeAssert)
}

// OpTypeAssert represents a bytecode operation for performing type assertions in a virtual machine.
// It embeds bytecode.Opcode to utilize opcode-related functionalities and operates on a core.IVMFullAccess instance.
type OpTypeAssert struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpTypeAssert creates a new instance of OpTypeAssert, ensuring the provided IVM implements IVMFullAccess.
func NewOpTypeAssert() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpTypeAssert{
		opcode: opcodes.NewOpcode(OpTypeAssertId, operands, "OpTypeAssert"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpTypeAssert) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpTypeAssert) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the type assertion operation. It attempts to assert if the top stack object matches the desired type.
// It decodes the target type from the constant pool using the provided decoder and validates against the interface object.
// On success, the concrete value and a boolean 'true' are pushed onto the stack. On failure, undefined and 'false' are pushed.
func (op *OpTypeAssert) Execute(decoder *core.Decoder) {
	// The operand is the index of the target type name in the constants table.
	typeNameIndex := decoder.Operand(0)
	interfaceObj := op.vm.StackPop()
	targetTypeObj := op.vm.Constants().Get(uint(typeNameIndex))
	targetTypeName, ok := targetTypeObj.(*objects.String)
	if !ok {
		op.vm.SetError(fmt.Errorf("constant for type assertion is not a string"))
		return
	}
	concreteValue := op.vm.Factory().UndefinedValue()
	switch io := interfaceObj.(type) {
	case *objects.Interface:
		concreteValue = io.Value()
	case *objects.Struct:
		concreteValue = io
	default:
		op.vm.StackPush(op.vm.Factory().UndefinedValue())
		op.vm.StackPush(op.vm.Factory().FalseValue())
		return
	}
	if concreteValue.TypeName() == targetTypeName.Value() {
		op.vm.StackPush(concreteValue)
		op.vm.StackPush(op.vm.Factory().TrueValue())
	} else {
		op.vm.StackPush(op.vm.Factory().UndefinedValue())
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Compile generates the compiled representation of the OpTypeAssert operation or returns an unimplemented error.
func (op *OpTypeAssert) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
