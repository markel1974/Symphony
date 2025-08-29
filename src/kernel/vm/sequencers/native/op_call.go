package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpCall)
}

// OpCall represents an operation code for invoking a function call in the virtual machine.
type OpCall struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpCall creates and returns a new instance of OpCall with initialized Opcode for the OpCall opcode.
func NewOpCall(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCall{
		Opcode: op.Opcode(bytecode.OpCall),
		vm:     vmT,
	}, nil
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	spread := decoder.Read(0)
	numArgs := decoder.Read(1)
	value := op.vm.Stack().PeekOffset(-1 - numArgs)
	op.vm.Call(value, spread == 1, numArgs)
}
