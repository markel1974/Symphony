package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpReturn)
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.Opcode.
type OpReturn struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpReturn creates a new instance of OpReturn with its Opcode initialized for the OpReturn operation.
func NewOpReturn(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpReturn{
		Opcode: op.Opcode(bytecode.OpReturn),
		vm:     vmT,
	}, nil
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the VM.
func (op *OpReturn) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	// collect return values from the stack using Pop(),
	// this is necessary to uncover the underlying values.
	var returnValues []objects.IObject
	if numReturnVals := decoder.Read(0); numReturnVals > 0 {
		returnValues = make([]objects.IObject, decoder.Read(0))
		for idx := numReturnVals - 1; idx >= 0; idx-- {
			returnValues[idx] = op.vm.Stack().Pop()
		}
	}
	op.vm.Return(returnValues)
}
