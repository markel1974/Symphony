package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpEqual)
}

// OpEqual represents an operation that checks if two values are equal and updates the stack accordingly.
type OpEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpEqual creates and returns an instance of OpEqual, initialized with its corresponding opcode details.
func NewOpEqual(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpEqual{OpcodeDetails: op.OpcodeToDetails(bytecode.OpEqual)}
}

// Execute performs the equality comparison between the top two stack values and pushes the result (true or false) back onto the stack.
func (op *OpEqual) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	right := v.Stack().Pop()
	left := v.Stack().Pop()
	var val objects.IObject
	if left.Equals(right) {
		val = op.Factory().FalseValue()
	} else {
		val = op.Factory().TrueValue()
	}
	v.Stack().Push(val)
}
