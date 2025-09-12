// File: vm/sequencers/native/op_is_type.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpIsType operation with the sequencer system using the SequencerRegister function.
func init() {
	SequencerRegister(NewOpIsType)
}

// OpIsType represents an operation that checks if a given object matches a specific type.
// It integrates with bytecode and provides full VM access for execution.
// The Opcode field provides details about the operation from the bytecode perspective.
// The vm field is used to interact with the virtual machine environment, providing complete execution control.
type OpIsType struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIsType creates a new instance of OpIsType executor if the provided VM implements IVMFullAccess.
// It associates the executor with the OpIsType opcode.
// Returns an error when the vm does not support the IVMFullAccess interface.
func NewOpIsType() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpIsType{
		opcode: opcodes.NewOpcode(OpIsTypeId, operands, "OpIsType"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIsType) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the runtime logic for the `OpIsType` operation, checking if the interface value matches the target type.
func (op *OpIsType) Execute(decoder *core.Decoder) {
	// The operand is the target type name index in the constants table.
	typeNameIndex := decoder.Operand(0)

	// The interface object is at the top of the stack (we don't use Pop).
	interfaceObj := op.vm.StackPeek()

	targetTypeObj := op.vm.Constants().Get(uint(typeNameIndex))
	targetTypeName, ok := targetTypeObj.(*objects.String)
	if !ok {
		op.vm.SetError(fmt.Errorf("constant for type check is not a string"))
		return
	}

	io, isInterface := interfaceObj.(*objects.Interface)
	if !isInterface {
		op.vm.StackPush(op.vm.Factory().FalseValue())
		return
	}

	if io.Value().TypeName() == targetTypeName.Value() {
		op.vm.StackPush(op.vm.Factory().TrueValue())
	} else {
		op.vm.StackPush(op.vm.Factory().FalseValue())
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIsType) Opcode() *opcodes.Opcode {
	return op.opcode
}
