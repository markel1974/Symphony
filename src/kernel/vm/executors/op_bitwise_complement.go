package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpBitwiseComplement)
}

// OpBitwiseComplement represents an operation for performing a bitwise complement on an operand.
// It extends OpcodeDetails, inheriting its metadata and behaviors.
type OpBitwiseComplement struct {
	*bytecode.OpcodeDetails
}

// NewOpBitwiseComplement initializes and returns an OpBitwiseComplement instance with the corresponding OpcodeDetails configuration.
func NewOpBitwiseComplement(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpBitwiseComplement{OpcodeDetails: op.OpcodeToDetails(bytecode.OpBitwiseComplement)}
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpBitwiseComplement) Execute(v *core.VM, _ *core.Decoder) {
	// Operands Offset 0
	operand := v.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.Factory().NewInt(v.FrameID(), ^x.Value())
		v.Stack().Push(res)
	default:
		v.SetError(fmt.Errorf("invalid operation: ^%s", operand.TypeName()))
		return
	}
}
