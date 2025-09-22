package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCall)
}

// OpCall represents an operation code for invoking a function call in the virtual machine.
type OpCall struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCall creates and returns a new instance of OpCall with initialized Opcode for the OpCall opcode.
func NewOpCall() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.SzUint8}
	return &OpCall{
		opcode: opcodes.NewOpcode(OpCallId, operands, "OpCall"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCall) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCall) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(decoder *core.Decoder) {
	spread := decoder.Operand(0)
	numArgs := decoder.Operand(1)
	hasSpread := spread > 0
	offset := numArgs + 1
	value := op.vm.StackPeekSP(uint(offset))
	op.vm.Call(value, hasSpread, numArgs)
}

// Compile generates the compiled representation of the OpCall operation or returns an unimplemented error.
func (op *OpCall) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
