package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpSetSelLocal represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds OpcodeDetails to utilize its properties like opcode, name, and operands.
type OpSetSelLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelLocal creates and returns a new instance of the OpSetSelLocal operation executor.
func NewOpSetSelLocal(op *bytecode.Opcodes) *OpSetSelLocal {
	return &OpSetSelLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetSelLocal)}
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpSetSelLocal) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	numSelectors := decoder.Read(0)
	localIndex := decoder.Read(1)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.Stack().PeekOffset(-numSelectors + i)
	}
	val := v.Stack().PeekOffset(-numSelectors - 1)
	v.Stack().DecrementCount(numSelectors + 1)
	dst := v.Stack().PeekAbsolute(v.BasePointer() + localIndex)
	if obj, ok := dst.(*objects.ObjectPointer); ok {
		dst = *obj.Value()
	}
	if err := op.Factory().IndexAssign(v.FrameID(), dst, val, selectors); err != nil {
		v.SetError(err)
		return
	}
}
