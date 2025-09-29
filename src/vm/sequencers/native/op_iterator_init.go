package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIteratorInit)
}

// OpIteratorInit represents an operation that initializes an iterator over an iterable object.
// It embeds Opcode for additional opcode-specific metadata.
type OpIteratorInit struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpIteratorInit{
		opcode: opcodes.NewOpcode(OpIteratorInitId, operands, "OpIteratorInit"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpIteratorInit) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpIteratorInit) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(decoder *handler.Decoder) {
	stableIdx := decoder.Operand(0)
	obj := op.vm.StackPop()
	if !obj.Iterable() {
		op.vm.SetError(objects.ComputeIteratorError(objects.ErrNotIterable, obj.TypeName()))
		return
	}
	itObj := obj.Iterate(op.vm.FrameId())
	op.vm.StackSetBP(uint(stableIdx), itObj)
}

// Compile generates the compiled representation of the OpIteratorInit operation or returns an unimplemented error.
func (op *OpIteratorInit) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
