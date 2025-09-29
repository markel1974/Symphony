package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeSet)
}

// OpFreeSet represents an operation to set the value of a free variable within a closure's environment.
type OpFreeSet struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpFreeSet creates and returns a new instance of OpFreeSet initialized with its corresponding Opcode.
func NewOpFreeSet() handler.IOpExecutor {
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

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpFreeSet) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpFreeSet) Execute(decoder *handler.Decoder) {
	freeIndex := decoder.Operand(0)
	o := op.vm.StackPop()
	freeObj := op.vm.FrameFreeVarsIndex(uint(freeIndex))
	if freeObj == nil {
		op.vm.Shutdown(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	op.vm.Factory().SetPointer(freeObj, o)
	//freeObj.SetValue(o)
}

// Compile generates the compiled representation of the OpFreePtrGet operation or returns an unimplemented error.
func (op *OpFreeSet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
