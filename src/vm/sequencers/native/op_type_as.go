// File: vm/sequencers/native/op_type_as.go

package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init initializes the package by registering the NewOpTypeAs operation with the sequencer system.
func init() {
	SequencerRegister(NewOpTypeAs)
}

// OpTypeAs represents an executor linked to the bytecode opcode OpTypeAs for handling unchecked casts in the Core.
// It embeds a bytecode.Opcode and uses handler.IVMFullAccess for full Core functionality.
type OpTypeAs struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpTypeAs creates a new instance of OpTypeAs executor for the given virtual machine and opcode.
// Returns an error if the provided Core does not implement IVMFullAccess.
func NewOpTypeAs() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpTypeAs{
		opcode: opcodes.NewOpcode(OpTypeAsId, operands, "OpTypeAs"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpTypeAs) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the OpTypeAs instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpTypeAs) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation by popping an interface from the stack and pushing its concrete value back.
// If the popped object is not an interface, it sets an error in the virtual machine.
func (op *OpTypeAs) Execute(_ *handler.Decoder) {
	interfaceObj := op.vm.StackPop()
	concrete := op.vm.Factory().Concrete(op.vm.FrameId(), interfaceObj)
	op.vm.StackPush(concrete)
}

// Compile generates the compiled representation of the OpTypeAs operation or returns an unimplemented error.
func (op *OpTypeAs) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
