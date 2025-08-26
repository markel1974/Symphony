package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

func init() {
	SequencerRegister(NewOpGetGlobal)
}

// OpGlobalGet represents an operation to retrieve a global variable in the virtual machine.
// It embeds Opcode for detailed opcode information.
type OpGlobalGet struct {
	*bytecode.Opcode
}

// NewOpGetGlobal creates a new instance of OpGlobalGet with its associated opcode details.
func NewOpGetGlobal(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpGlobalGet{Opcode: op.Opcode(bytecode.OpGlobalGet)}
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGlobalGet) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	glIndex := decoder.Read(0)
	glObj := v.Globals().Get(uint(glIndex))
	if glObj == nil {
		v.SetError(fmt.Errorf("undefined global: %d", glIndex))
		return
	}
	v.Stack().Push(glObj)
}
