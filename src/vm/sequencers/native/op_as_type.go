// File: vm/sequencers/native/op_as_type.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the package by registering the NewOpAsType operation with the sequencer system.
func init() {
	SequencerRegister(NewOpAsType)
}

// OpAsType represents an executor linked to the bytecode opcode OpAsType for handling unchecked casts in the VM.
// It embeds a bytecode.Opcode and uses core.IVMFullAccess for full VM functionality.
type OpAsType struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpAsType creates a new instance of OpAsType executor for the given virtual machine and opcode.
// Returns an error if the provided VM does not implement IVMFullAccess.
func NewOpAsType() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpAsType{
		opcode: opcodes.NewOpcode(OpAsTypeId, operands, "OpAsType"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpAsType) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the OpAsType instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpAsType) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation by popping an interface from the stack and pushing its concrete value back.
// If the popped object is not an interface, it sets an error in the virtual machine.
func (op *OpAsType) Execute(_ *core.Decoder) {
	interfaceObj := op.vm.StackPop()
	concrete := op.vm.Factory().Concrete(interfaceObj)
	op.vm.StackPush(concrete)
}

// Compile generates the compiled representation of the OpAsType operation or returns an unimplemented error.
func (op *OpAsType) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
