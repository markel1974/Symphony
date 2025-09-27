package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the `NewOpCallAsync` operation with the sequencer through the `SequencerRegister` function.
func init() {
	SequencerRegister(NewOpCallAsync)
}

// OpCallAsync represents an operation that performs an asynchronous function-like call within the virtual machine.
// It requires a compatible implementation of IVMFullAccess to bind and execute the operation.
type OpCallAsync struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCallAsync creates and returns a new instance of OpCallAsync as an IOpExecutor implementation.
func NewOpCallAsync() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.SzUint8}
	return &OpCallAsync{
		opcode: opcodes.NewOpcode(OpCallAsyncId, operands, "OpCallAsync"),
		vm:     nil,
	}
}

// Opcode retrieves the Opcode instance associated with the OpCallAsync.
func (op *OpCallAsync) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind assigns a virtual machine instance to the OpCallAsync object, ensuring it implements the IVMFullAccess interface.
// Returns an error if the provided VM does not meet the required interface implementation.
func (op *OpCallAsync) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the core logic of the OpCallAsync by processing operands, evaluating the stack, and invoking a call.
func (op *OpCallAsync) Execute(decoder *core.Decoder) {
	spread := decoder.Operand(0)
	numArgs := decoder.Operand(1)
	hasSpread := spread > 0
	offset := numArgs + 1
	value := op.vm.StackPeekSP(uint(offset))
	op.vm.CallAsync(value, hasSpread, numArgs)
}

// Compile generates a bytecode sequence for the operation or signals it's unimplemented by returning an error.
func (op *OpCallAsync) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
