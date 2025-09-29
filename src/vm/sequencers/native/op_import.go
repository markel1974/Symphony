package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpImport)
}

// OpImport extends Opcode to represent operations specifically related to reference handling in the bytecode.
type OpImport struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpImport initializes a new OpImport instance with corresponding Opcode from the bytecode package.
func NewOpImport() handler.IOpExecutor {
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

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpImport) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the specified Core instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpImport) Execute(decoder *handler.Decoder) {
	importIdx := decoder.Operand(0)
	importObj, err := op.vm.ImportsGet(uint(importIdx))
	if err != nil {
		op.vm.Shutdown(err)
		return
	}
	op.vm.StackPush(importObj)
}

// Compile generates the compiled representation of the OpImport operation or returns an unimplemented error.
func (op *OpImport) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
