package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpLocalGet)
}

// OpLocalGet represents an operation to retrieve a local variable from the stack using its index.
type OpLocalGet struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpLocalGet creates a new OpLocalGet instance and initializes it with details for the OpLocalGet opcode.
func NewOpLocalGet(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpLocalGet{
		Opcode: op.Opcode(bytecode.OpLocalGet),
		vm:     vmT,
	}, nil
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpLocalGet) Execute(decoder *core.Decoder) {
	// Operands Offset 1 (8-bit)
	localIndex := decoder.Read(0)
	val := op.vm.Stack().PeekAbsolute(op.vm.Frame().BasePointer() + localIndex)
	op.vm.Stack().Push(val)
}
