package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpCall)
}

// OpCall represents an operation code for invoking a function call in the virtual machine.
type OpCall struct {
	*bytecode.Opcode
}

// NewOpCall creates and returns a new instance of OpCall with initialized Opcode for the OpCall opcode.
func NewOpCall(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpCall{Opcode: op.Opcode(bytecode.OpCall)}
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	spread := decoder.Read(0)
	numArgs := decoder.Read(1)
	value := v.Stack().PeekOffset(-1 - numArgs)
	v.Call(value, spread == 1, numArgs)
}
