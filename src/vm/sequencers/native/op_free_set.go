package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeSet)
}

// OpFreeSet represents an operation to set the value of a free variable within a closure's environment.
type OpFreeSet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpFreeSet creates and returns a new instance of OpFreeSet initialized with its corresponding Opcode.
func NewOpFreeSet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpFreeSet{
		opcode: opcodes.NewOpcode(OpFreeSetId, operands, "OpFreeSet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpFreeSet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpFreeSet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpFreeSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	freeIndex := decoder.Operand(0)
	o := op.vm.StackPop()
	freeObj := op.vm.FrameFreeVarsIndex(uint(freeIndex))
	if freeObj == nil {
		op.vm.SetError(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	op.vm.Factory().SetPointer(freeObj, o)
	//freeObj.SetValue(o)
}

// Compile generates the compiled representation of the OpFreePtrGet operation or returns an unimplemented error.
func (op *OpFreeSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
