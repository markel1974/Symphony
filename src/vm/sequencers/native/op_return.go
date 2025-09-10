package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpReturn)
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.Opcode.
type OpReturn struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpReturn creates a new instance of OpReturn with its Opcode initialized for the OpReturn operation.
func NewOpReturn() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8}
	return &OpReturn{
		opcode: opcodes.NewOpcode(OpReturnId, operands, "OpReturn"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpReturn) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the VM.
func (op *OpReturn) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	// collect return values from the stack using Pop(),
	// this is necessary to uncover the underlying values.
	var returnValues []objects.IObject
	if numReturnVals := decoder.Operand(0); numReturnVals > 0 {
		returnValues = make([]objects.IObject, numReturnVals)
		for i := 0; i < numReturnVals; i++ {
			returnValues[i] = op.vm.Stack().Pop()
		}
	}
	op.vm.Return(returnValues)
}

// Opcode returns the opcode associated with the instance.
func (op *OpReturn) Opcode() *opcodes.Opcode {
	return op.opcode
}
