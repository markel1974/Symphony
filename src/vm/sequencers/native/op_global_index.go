package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpGlobalIndex)
}

// OpGlobalIndex represents an operation for setting a global variable's value using selectors for indexing or access.
type OpGlobalIndex struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpGlobalIndex creates a new instance of OpGlobalIndex with its corresponding Opcode initialized.
func NewOpGlobalIndex() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.Relocatable}
	return &OpGlobalIndex{
		opcode: opcodes.NewOpcode(OpGlobalIndexId, operands, "OpGlobalIndex"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalIndex) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpGlobalIndex) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation defined by OpGlobalIndex, updating the Core state and handling global index assignment.
func (op *OpGlobalIndex) Execute(decoder *handler.Decoder) {
	globalIndex := decoder.Operand(0)
	selCount := decoder.Operand(1)
	dstObj, err := op.vm.GlobalsGet(uint(globalIndex))
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	selectors := make([]objects.IObject, selCount)
	for i := 0; i < selCount; i++ {
		offset := selCount - i
		selectors[i] = op.vm.StackPeekSP(uint(offset))
	}
	offset := selCount + 1
	srcObj := op.vm.StackPeekSP(uint(offset))
	op.vm.StackDecrementCount(uint(offset))
	if err := op.vm.Factory().IndexAssign(op.vm.FrameId(), dstObj, srcObj, selectors); err != nil {
		op.vm.Shutdown(err)
		return
	}
}

// Compile generates the compiled representation of the OpGlobalIndex operation or returns an unimplemented error.
func (op *OpGlobalIndex) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
