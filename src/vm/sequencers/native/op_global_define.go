package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
)

func init() {
	SequencerRegister(NewOpGlobalDefine)
}

// OpGlobalDefine represents the opcode for defining a new local variable within the current frame's scope.
type OpGlobalDefine struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpGlobalDefine creates a new instance of OpGlobalDefine with its associated opcode details.
func NewOpGlobalDefine(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpGlobalDefine{
		Opcode: op.Opcode(bytecode.OpGlobalDefine),
		vm:     vmT,
	}, nil
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpGlobalDefine) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	glIndex := decoder.Read(0)
	val := op.vm.Stack().Peek()
	op.vm.Globals().Set(uint(glIndex), val)
}
