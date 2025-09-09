package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the OpLocalPtrGet operation in the sequencer system when the package is initialized.
func init() {
	SequencerRegister(NewOpGlobalPtrGet)
}

// OpGlobalPtrGet represents an executable opcode for retrieving a global pointer in the virtual machine context.
type OpGlobalPtrGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalPtrGet creates a new OpGlobalPtrGet executor, verifying the VM implements IVMFullAccess and initializing its Opcode.
// It returns the created core.IOpExecutor or an error if the VM check fails.
func NewOpGlobalPtrGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalPtrGet{
		Opcode: op.Opcode(bytecode.OpGlobalPtrGet),
		vm:     vmT,
	}, nil
}

// Execute retrieves a global value by its index, checks its type, and pushes a pointer object onto the stack.
func (op *OpGlobalPtrGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	globalIndex := decoder.Read(0)
	val := op.vm.Globals().Get(uint(globalIndex))
	if obj, ok := val.(*objects.ObjectPointer); ok {
		op.vm.Stack().Push(obj)
		return
	}
	freeVar := op.vm.Factory().NewObjectPointer(op.vm.Frame().Id(), &val)
	op.vm.Stack().Push(freeVar)
}
