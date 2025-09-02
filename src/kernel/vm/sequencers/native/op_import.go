package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpImport)
}

// OpImport extends Opcode to represent operations specifically related to reference handling in the bytecode.
type OpImport struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpImport initializes a new OpImport instance with corresponding Opcode from the bytecode package.
func NewOpImport(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpImport{
		Opcode: op.Opcode(bytecode.OpImport),
		vm:     vmT,
	}, nil
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpImport) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	nameIndex := decoder.Read(0)
	symbol := op.vm.Imports().Get(uint(nameIndex))
	op.vm.Stack().Push(symbol)
}
