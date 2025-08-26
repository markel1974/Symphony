package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// OpGetFree represents an operation to retrieve a free variable in a closure during execution.
type OpGetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFree creates and returns a new instance of OpGetFree, initializing its OpcodeDetails using bytecode metadata.
func NewOpGetFree(op *bytecode.Opcodes) *OpGetFree {
	return &OpGetFree{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetFree)}
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpGetFree) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 1
	freeIndex := decoder.Read(0)
	//if freeIndex != int(v.currFrame.Get8(v.ip)) {
	//	log.Println("local OpGetFree mismatch")
	//}
	val := *v.FreeVarsIndex(freeIndex).Value()
	v.Stack().Push(val)
}
