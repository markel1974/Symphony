package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeGetPtr)
}

// OpFreePtrGet represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds Opcode, which provides opcode metadata such as identifier, operands, and name.
type OpFreePtrGet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpFreeGetPtr creates a new instance of OpFreePtrGet initialized with the corresponding Opcode.
func NewOpFreeGetPtr() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpFreePtrGet{
		opcode: opcodes.NewOpcode(OpFreePtrGetId, operands, "OpFreePtrGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpFreePtrGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpFreePtrGet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute executes the OpFreePtrGet operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpFreePtrGet) Execute(decoder *handler.Decoder) {
	freeIndex := decoder.Operand(0)
	val := op.vm.FrameFreeVarsIndex(uint(freeIndex))
	if val == nil {
		op.vm.Shutdown(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	op.vm.StackPush(val)
}

// Compile generates the compiled representation of the OpFreePtrGet operation or returns an unimplemented error.
func (op *OpFreePtrGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
