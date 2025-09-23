package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalGet)
}

// OpGlobalGet represents an operation to retrieve a global variable in the virtual machine.
// It embeds Opcode for detailed opcode information.
type OpGlobalGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpGlobalGet creates a new instance of OpGlobalGet with its associated opcode details.
func NewOpGlobalGet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpGlobalGet{
		opcode: opcodes.NewOpcode(OpGlobalGetId, operands, "OpGlobalGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpGlobalGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGlobalGet) Execute(decoder *core.Decoder) {
	index := decoder.Operand(0)
	obj := op.vm.GlobalsGet(uint(index))
	op.vm.StackPush(obj)
}

// Compile generates the compiled representation of the OpGlobalGet operation or returns an unimplemented error.
func (op *OpGlobalGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
