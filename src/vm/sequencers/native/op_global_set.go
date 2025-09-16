package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalSet)
}

// OpGlobalSet represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpGlobalSet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpGlobalSet creates and returns a new instance of OpGlobalSet with initialized Opcode.
func NewOpGlobalSet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpGlobalSet{
		opcode: opcodes.NewOpcode(OpGlobalSetId, operands, "OpGlobalSet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalSet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpGlobalSet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpGlobalSet) Execute(decoder *core.Decoder) {
	index := decoder.Operand(0)
	val := op.vm.StackPeek()
	op.vm.Globals().Set(uint(index), val)
}

// Compile generates the compiled representation of the OpGlobalSet operation or returns an unimplemented error.
func (op *OpGlobalSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
