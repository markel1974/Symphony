package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalSet)
}

// OpLocalSet represents an operation to set the value of a local variable within the current frame.
// It embeds Opcode for opcode-specific information such as name, operands, and code.
type OpLocalSet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpLocalSet initializes and returns a new instance of OpLocalSet with associated opcode details.
func NewOpLocalSet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalSet{
		opcode: opcodes.NewOpcode(OpLocalSetId, operands, "OpLocalSet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLocalSet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpLocalSet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	localIndex := decoder.Operand(0)
	val := op.vm.Stack().Peek()
	obj := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	if freeObj, ok := obj.(*objects.ObjectPointer); ok {
		op.vm.Factory().SetPointer(freeObj, val)
	} else {
		op.vm.Stack().SetAbsolute(op.vm.Frame().BasePointer()+localIndex, val)
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalSet) Opcode() *opcodes.Opcode {
	return op.opcode
}
