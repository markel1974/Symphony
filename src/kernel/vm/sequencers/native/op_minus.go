package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpMinus)
}

// OpMinus represents an operation for negating a numeric value.
// It embeds OpcodeDetails, providing details such as the opcode, operands, and name.
type OpMinus struct {
	*bytecode.OpcodeDetails
}

// NewOpMinus creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpMinus(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpMinus{OpcodeDetails: op.OpcodeToDetails(bytecode.OpMinus)}
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpMinus) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	operand := v.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.Factory().NewInt(v.FrameID(), -x.Value())
		v.Stack().Push(res)
	case *objects.Float:
		res := op.Factory().NewFloat(v.FrameID(), -x.Value())
		v.Stack().Push(res)
	default:
		v.SetError(fmt.Errorf("invalid operation: -%s", operand.TypeName()))
	}
}
