package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeGetPtr)
}

// OpFreePtrGet represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds Opcode, which provides opcode metadata such as identifier, operands, and name.
type OpFreePtrGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpFreeGetPtr creates a new instance of OpFreePtrGet initialized with the corresponding Opcode.
func NewOpFreeGetPtr() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpFreePtrGet{
		opcode: opcodes.NewOpcode(OpFreePtrGetId, operands, "OpFreePtrGet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpFreePtrGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute executes the OpFreePtrGet operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpFreePtrGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	freeIndex := decoder.Operand(0)
	val := op.vm.Frame().FreeVarsIndex(uint(freeIndex))
	if val == nil {
		op.vm.SetError(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	op.vm.Stack().Push(val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpFreePtrGet) Opcode() *opcodes.Opcode {
	return op.opcode
}
