package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*bytecode.OpcodeDetails
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding OpcodeDetails configuration set.
func NewOpUnknown(op *bytecode.Opcodes) *OpUnknown {
	return &OpUnknown{OpcodeDetails: op.OpcodeToDetails(bytecode.OpUnknown)}
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 0
	//currentIp := int(v.currFrame.Get8(v.ip))
	//if nameIndex != int(v.currFrame.Get8(v.ip)) {
	//	log.Println("name index mismatch: %d != %d", nameIndex, int(v.currFrame.Get16(v.ip)))
	//}
	v.SetError(fmt.Errorf("unknown opcode at: %d", v.GetIp()))
}
