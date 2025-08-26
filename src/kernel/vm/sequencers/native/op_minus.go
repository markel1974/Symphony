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
// It embeds Opcode, providing details such as the opcode, operands, and name.
type OpMinus struct {
	*bytecode.Opcode
}

// NewOpMinus creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpMinus(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpMinus{Opcode: op.Opcode(bytecode.OpMinus)}
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpMinus) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	operand := v.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := v.Factory().NewInt(v.Frame().Id(), -x.Value())
		v.Stack().Push(res)
	case *objects.Float:
		res := v.Factory().NewFloat(v.Frame().Id(), -x.Value())
		v.Stack().Push(res)
	default:
		v.SetError(fmt.Errorf("invalid operation: -%s", operand.TypeName()))
	}
}
