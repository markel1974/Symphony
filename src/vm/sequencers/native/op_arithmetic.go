package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// init registers the NewOpLogical operation with the sequencer system by adding it to the _registerContainer.
func init() {
	SequencerRegister(NewOpArithmetic)
}

// OpArithmetic represents an operation responsible for performing arithmetic computations within the virtual machine.
type OpArithmetic struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpArithmetic creates a new arithmetic operation executor with the given virtual machine and opcode.
// It ensures the virtual machine implements the IVMFullAccess interface before creating the executor.
// Returns an IOpExecutor for arithmetic operations or an error if the provided VM does not support full access.
func NewOpArithmetic(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpArithmetic{
		Opcode: op.Opcode(bytecode.OpArithmetic),
		vm:     vmT,
	}, nil
}

// Execute processes the arithmetic operation by decoding the operator and applying it to the top two stack elements.
// It retrieves the operator from the bytecode, performs the calculation, and pushes the result back onto the stack.
// In case of an operation error, it sets the virtual machine's error state.
func (op *OpArithmetic) Execute(decoder *core.Decoder) {
	// Operands Offset  1 (8 bits)
	opcode := decoder.Read(0)
	right := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	operator := objects.ArithmeticOperator(opcode)
	res, err := left.ArithmeticOp(op.vm.Frame().Id(), operator, right)
	if err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.Stack().Push(res)
}
