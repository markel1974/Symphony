package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/handler"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the sequencer by registering the NewOpGlobalDefine operation with the SequencerRegister function.
func init() {
	SequencerRegister(NewOpGlobalDefine)
}

// OpGlobalDefine represents an operation to define a global variable in the Core environment.
// It binds a value from the stack to a global index specified by the decoder.
// Embeds `bytecode.Opcode` for opcode-related operations and uses `handler.IVMFullAccess` for Core interactions.
type OpGlobalDefine struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpGlobalDefine creates a new OpGlobalDefine executor for the given virtual machine and opcode configuration.
// Returns an IOpExecutor instance or an error if the Core does not implement IVMFullAccess.
func NewOpGlobalDefine() handler.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
	return &OpGlobalDefine{
		opcode: opcodes.NewOpcode(OpGlobalDefineId, operands, "OpGlobalDefine"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpGlobalDefine) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided Core to IVMFullAccess and storing it.
// Returns an error if the Core does not implement the required interface.
func (op *OpGlobalDefine) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute sets a value from the stack into the global variables using the operand index from the decoder.
func (op *OpGlobalDefine) Execute(decoder *handler.Decoder) {
	index := decoder.Operand(0)
	val := op.vm.StackPeek()
	op.vm.GlobalsSet(uint(index), val)
}

// Compile generates the compiled representation of the OpGlobalDefine operation or returns an unimplemented error.
func (op *OpGlobalDefine) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
