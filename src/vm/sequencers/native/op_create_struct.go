package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateStruct)
}

// OpCreateStruct is a wrapper around bytecode.Opcode, representing a struct creation operation in bytecode execution.
type OpCreateStruct struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCreateStruct initializes and returns a new instance of OpCreateStruct with its Opcode set to OpCreateMap details.
func NewOpCreateStruct() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpCreateStruct{
		opcode: opcodes.NewOpcode(OpCreateStructId, operands, "OpCreateStruct"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateStruct) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateStruct) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the OpCreateMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpCreateStruct) Execute(decoder *core.Decoder) {
	numElem := decoder.Operand(0)
	typeNameObj := op.vm.StackPop()
	structObj := op.vm.StackPopStruct(typeNameObj.AsString(), uint(numElem))
	op.vm.StackPush(structObj)
}

// Compile generates the compiled representation of the OpCreateStruct operation or returns an unimplemented error.
func (op *OpCreateStruct) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
