package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJumpAnd)
}

// OpJumpAnd represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpJumpAnd struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpAnd creates and returns a new instance of OpJumpAnd, initializing it with details for the OpJumpAnd opcode.
func NewOpJumpAnd(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpAnd{
		Opcode: op.Opcode(bytecode.OpJumpAnd),
		vm:     vmT,
	}, nil
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpJumpAnd) Execute(decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	obj := op.vm.Stack().Peek()
	if obj.Falsy() {
		pos := decoder.Read(0)
		op.vm.SetIp(pos - 1)
	} else {
		op.vm.Stack().Decrement()
	}
}
