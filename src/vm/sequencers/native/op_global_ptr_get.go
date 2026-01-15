package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the OpLocalPtrGet operation in the sequencer system when the package is initialized.
func init() {
	SequencerRegister(NewOpGlobalPtrGet)
}

// OpGlobalPtrGet represents an executable opcode for retrieving a global pointer in the virtual machine context.
type OpGlobalPtrGet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpGlobalPtrGet creates a new OpGlobalPtrGet executor, verifying the Core implements IVMFullAccess and initializing its Opcode.
// It returns the created handler.IOpExecutor or an error if the Core check fails.
func NewOpGlobalPtrGet() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpGlobalPtrGet{
		opcode: opcodes.NewOpcode(OpGlobalPtrGetId, operands, "OpGlobalPtrGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalPtrGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpGlobalPtrGet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute retrieves a global value by its index, checks its type, and pushes a pointer object onto the stack.
func (op *OpGlobalPtrGet) Execute(decoder *handler.Decoder) {
	globalIndex := decoder.Operand(0)
	val, err := op.vm.GlobalsGet(uint(globalIndex))
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	ptr := op.vm.CreateObjectPointer(val)
	op.vm.StackPush(ptr)
}

// Compile generates the compiled representation of the OpGlobalPtrGet operation or returns an unimplemented error.
func (op *OpGlobalPtrGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
