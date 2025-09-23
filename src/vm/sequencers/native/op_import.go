package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpImport)
}

// OpImport extends Opcode to represent operations specifically related to reference handling in the bytecode.
type OpImport struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpImport initializes a new OpImport instance with corresponding Opcode from the bytecode package.
func NewOpImport() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpImport{
		opcode: opcodes.NewOpcode(OpImportId, operands, "OpImport"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpImport) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpImport) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpImport) Execute(decoder *core.Decoder) {
	importIdx := decoder.Operand(0)
	importObj := op.vm.ImportsGet(uint(importIdx))
	op.vm.StackPush(importObj)
}

// Compile generates the compiled representation of the OpImport operation or returns an unimplemented error.
func (op *OpImport) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
