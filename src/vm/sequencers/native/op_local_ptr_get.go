package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalPtrGet)
}

// OpLocalPtrGet retrieves a local variable as a pointer using its index within the current frame.
type OpLocalPtrGet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its Opcode.
func NewOpLocalPtrGet() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalPtrGet{
		opcode: opcodes.NewOpcode(OpLocalPtrGetId, operands, "OpLocalPtrGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalPtrGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpLocalPtrGet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpLocalPtrGet) Execute(decoder *handler.Decoder) {
	localIndex := decoder.Operand(0)
	val := op.vm.StackPeekBP(uint(localIndex))
	ptr := op.vm.CreateObjectPointer(val)
	op.vm.StackPush(ptr)
}

// Compile generates the compiled representation of the OpLocalPtrGet operation or returns an unimplemented error.
func (op *OpLocalPtrGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
