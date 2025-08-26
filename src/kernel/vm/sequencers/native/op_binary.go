package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpBinary)
}

// OpBinary represents a type that performs binary operations by extending bytecode.Opcode.
type OpBinary struct {
	*bytecode.Opcode
}

// NewOpBinary creates a new instance of OpBinary with its corresponding Opcode initialized.
func NewOpBinary(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpBinary{Opcode: op.Opcode(bytecode.OpBinary)}
}

// Execute performs a binary operation using operands from the stack, updates the instruction pointer, and handles errors.
func (op *OpBinary) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset  1 (8 bits)
	opcode := decoder.Read(0)
	right := v.Stack().Pop()
	left := v.Stack().Pop()
	operator := objects.Operator(opcode)
	res, err := left.BinaryOp(v.Frame().Id(), operator, right)
	if err != nil {
		v.SetError(err)
		return
	}
	v.Stack().Push(res)
}
