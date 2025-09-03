package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpFalsy)
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding Opcode.
func NewOpJumpFalsy(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpFalsy{
		Opcode: op.Opcode(bytecode.OpJumpFalsy),
		vm:     vmT,
	}, nil
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	obj := op.vm.Stack().Pop()
	if obj.Falsy() {
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	}
}
