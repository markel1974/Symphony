package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpJump)
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJump{
		Opcode: op.Opcode(bytecode.OpJump),
		vm:     vmT,
	}, nil
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`VM`) to a calculated position in the frame.
func (op *OpJump) Execute(decoder *core.Decoder) {
	// Operands Offset  2 (16-bit)
	pos := decoder.Read(0)
	op.vm.SetIp(pos - 1)
}
