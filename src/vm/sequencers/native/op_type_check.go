// File: vm/sequencers/native/op_type_check.go

package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the NewOpTypeCheck operation with the sequencer system using the SequencerRegister function.
func init() {
	SequencerRegister(NewOpTypeCheck)
}

// OpTypeCheck represents an operation that checks if a given object matches a specific type.
// It integrates with bytecode and provides full Core access for execution.
// The Opcode field provides details about the operation from the bytecode perspective.
// The vm field is used to interact with the virtual machine environment, providing complete execution control.
type OpTypeCheck struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpTypeCheck creates a new instance of OpTypeCheck executor if the provided Core implements IVMFullAccess.
// It associates the executor with the OpTypeCheck opcode.
// Returns an error when the vm does not support the IVMFullAccess interface.
func NewOpTypeCheck() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpTypeCheck{
		opcode: opcodes.NewOpcode(OpTypeCheckId, operands, "OpTypeCheck"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpTypeCheck) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpTypeCheck) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the runtime logic for the `OpTypeCheck` operation, checking if the interface value matches the target type.
func (op *OpTypeCheck) Execute(decoder *handler.Decoder) {
	interfaceObj := op.vm.StackPeek()
	switch io := interfaceObj.(type) {
	case *objects.Interface:
		typeNameIndex := decoder.Operand(0)
		targetTypeObj, err := op.vm.ConstantsGet(uint(typeNameIndex))
		if err != nil {
			op.vm.Shutdown(err)
			return
		}
		if io.Value().TypeName() == targetTypeObj.AsString() {
			op.vm.StackPush(op.vm.Factory().TrueValue())
		} else {
			op.vm.StackPush(op.vm.Factory().FalseValue())
		}
	default:
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Compile generates the compiled representation of the OpTypeCheck operation or returns an unimplemented error.
func (op *OpTypeCheck) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
