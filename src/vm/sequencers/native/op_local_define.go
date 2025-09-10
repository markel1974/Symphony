package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpLocalDefine)
}

// OpLocalDefine represents the opcode for defining a new local variable within the current frame's scope.
type OpLocalDefine struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpLocalDefine creates a new instance of OpLocalDefine with its associated opcode details.
func NewOpLocalDefine() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpLocalDefine{
		opcode: opcodes.NewOpcode(OpLocalDefineId, operands, "OpLocalDefine"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpLocalDefine) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpLocalDefine) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	localIndex := decoder.Operand(0)
	val := op.vm.Stack().Peek()
	destSlot := op.vm.Frame().BasePointer() + localIndex
	op.vm.Stack().SetAbsolute(destSlot, val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpLocalDefine) Opcode() *opcodes.Opcode {
	return op.opcode
}
