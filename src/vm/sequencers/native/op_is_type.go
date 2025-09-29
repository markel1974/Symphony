// File: vm/sequencers/native/op_is_type.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpIsType operation with the sequencer system using the SequencerRegister function.
func init() {
	SequencerRegister(NewOpIsType)
}

// OpIsType represents an operation that checks if a given object matches a specific type.
// It integrates with bytecode and provides full Core access for execution.
// The Opcode field provides details about the operation from the bytecode perspective.
// The vm field is used to interact with the virtual machine environment, providing complete execution control.
type OpIsType struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIsType creates a new instance of OpIsType executor if the provided Core implements IVMFullAccess.
// It associates the executor with the OpIsType opcode.
// Returns an error when the vm does not support the IVMFullAccess interface.
func NewOpIsType() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpIsType{
		opcode: opcodes.NewOpcode(OpIsTypeId, operands, "OpIsType"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIsType) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIsType) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the runtime logic for the `OpIsType` operation, checking if the interface value matches the target type.
func (op *OpIsType) Execute(decoder *handler.Decoder) {
	interfaceObj := op.vm.StackPeek()
	switch io := interfaceObj.(type) {
	case *objects.Interface:
		typeNameIndex := decoder.Operand(0)
		targetTypeObj := op.vm.ConstantsGet(uint(typeNameIndex))
		if io.Value().TypeName() == targetTypeObj.AsString() {
			op.vm.StackPush(op.vm.Factory().TrueValue())
		} else {
			op.vm.StackPush(op.vm.Factory().FalseValue())
		}
	default:
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Compile generates the compiled representation of the OpIsType operation or returns an unimplemented error.
func (op *OpIsType) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
