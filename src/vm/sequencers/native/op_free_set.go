package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeSet)
}

// OpFreeSet represents an operation to set the value of a free variable within a closure's environment.
type OpFreeSet struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpFreeSet creates and returns a new instance of OpFreeSet initialized with its corresponding Opcode.
func NewOpFreeSet(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpFreeSet{
		Opcode: op.Opcode(opcodes.OpFreeSet),
		vm:     vmT,
	}, nil
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpFreeSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	freeIndex := decoder.Read(0)
	o := op.vm.Stack().Pop()
	freeObj := op.vm.Frame().FreeVarsIndex(uint(freeIndex))
	if freeObj == nil {
		op.vm.SetError(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	op.vm.Factory().SetPointer(freeObj, o)
	//freeObj.SetValue(o)
}
