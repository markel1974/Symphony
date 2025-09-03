package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	objects2 "github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpMinus)
}

// OpMinus represents an operation for negating a numeric value.
// It embeds Opcode, providing details such as the opcode, operands, and name.
type OpMinus struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpMinus creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpMinus(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpMinus{
		Opcode: op.Opcode(bytecode.OpMinus),
		vm:     vmT,
	}, nil
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpMinus) Execute(_ *core.Decoder) {
	// Operands Offset 0
	operand := op.vm.Stack().Pop()
	switch x := operand.(type) {
	case *objects2.Int:
		res := op.vm.Factory().NewInt(op.vm.Frame().Id(), -x.Value())
		op.vm.Stack().Push(res)
	case *objects2.Float:
		res := op.vm.Factory().NewFloat(op.vm.Frame().Id(), -x.Value())
		op.vm.Stack().Push(res)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: -%s", operand.TypeName()))
	}
}
