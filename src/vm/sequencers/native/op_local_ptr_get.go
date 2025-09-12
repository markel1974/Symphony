package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalPtrGet)
}

// OpLocalPtrGet retrieves a local variable as a pointer using its index within the current frame.
type OpLocalPtrGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpLocalPtrGet creates and returns a new instance of OpLocalPtrGet, initializing its Opcode.
func NewOpLocalPtrGet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalPtrGet{
		opcode: opcodes.NewOpcode(OpLocalPtrGetId, operands, "OpLocalPtrGet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLocalPtrGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpLocalPtrGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	localIndex := decoder.Operand(0)
	val := op.vm.StackPeekOffsetBP(uint(localIndex))
	if obj, ok := val.(*objects.ObjectPointer); ok {
		op.vm.StackPush(obj)
		return
	}
	freeVar := op.vm.Factory().NewObjectPointer(op.vm.FrameId(), &val)
	op.vm.StackPush(freeVar)
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalPtrGet) Opcode() *opcodes.Opcode {
	return op.opcode
}
