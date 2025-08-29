package native

import (
	"fmt"

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
	vm core.IVMFullAccess
}

// NewOpBinary creates a new instance of OpBinary with its corresponding Opcode initialized.
func NewOpBinary(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpBinary{
		Opcode: op.Opcode(bytecode.OpBinary),
		vm:     vmT,
	}, nil
}

// Execute performs a binary operation using operands from the stack, updates the instruction pointer, and handles errors.
func (op *OpBinary) Execute(decoder *core.Decoder) {
	// Operands Offset  1 (8 bits)
	opcode := decoder.Read(0)
	right := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	operator := objects.Operator(opcode)
	res, err := left.BinaryOp(op.vm.Frame().Id(), operator, right)
	if err != nil {
		op.vm.SetError(err)
		return
	}
	op.vm.Stack().Push(res)
}
