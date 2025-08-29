package native

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpReferences)
}

// OpReferences extends Opcode to represent operations specifically related to reference handling in the bytecode.
type OpReferences struct {
	*bytecode.Opcode
	vm *core.VM
}

// NewOpReferences initializes a new OpReferences instance with corresponding Opcode from the bytecode package.
func NewOpReferences(vm *core.VM, op *bytecode.Opcodes) core.IOpExecutor {
	return &OpReferences{
		Opcode: op.Opcode(bytecode.OpReferences),
		vm:     vm,
	}
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpReferences) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	nameIndex := decoder.Read(0)
	symbol := op.vm.References().Get(uint(nameIndex))
	op.vm.Stack().Push(symbol)
}
