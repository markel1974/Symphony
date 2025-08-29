package native

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
// It extends Opcode, inheriting its metadata and behaviors.
type OpBitwiseComplement struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpBitwiseComplement initializes and returns an OpBitwiseComplement instance with the corresponding Opcode configuration.
func NewOpBitwiseComplement(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpBitwiseComplement{
		Opcode: op.Opcode(bytecode.OpBitwiseComplement),
		vm:     vm,
	}
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpBitwiseComplement) Execute(_ *core.Decoder) {
	// Operands Offset 0
	operand := op.vm.Stack().Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.vm.Factory().NewInt(op.vm.Frame().Id(), ^x.Value())
		op.vm.Stack().Push(res)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: ^%s", operand.TypeName()))
		return
	}
}
