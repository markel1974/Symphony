package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpIntOp)
}

// OpIntOp extends OpcodeDetails and represents integer operations performed on a virtual machine.
type OpIntOp struct {
	*bytecode.OpcodeDetails
}

// NewOpIntOp initializes and returns a new instance of OpIntOp with relevant opcode details provided by bytecode.Opcodes.
func NewOpIntOp(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpIntOp{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIntOp)}
}

// Execute performs a specified binary operation between two integers from the stack and stores the result in a destination slot.
// It retrieves operands, validates types, and executes the operation, setting an error on unsupported cases or type mismatches.
func (op *OpIntOp) Execute(v *core.VM, decoder *core.Decoder) {
	dstObj := v.Stack().PeekAbsolute(decoder.Read(0))
	binaryOp := objects.Operator(decoder.Read(1))
	dst, ok := dstObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	rhsObj := v.Stack().Pop()
	rhs, ok := rhsObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("rhs expected int, got %s", rhsObj.TypeName()))
		return
	}
	lhsObj := v.Stack().Pop()
	lhs, ok := lhsObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("lhs expected int, got %s", lhsObj.TypeName()))
		return
	}
	result, err := op.Factory().BinaryOpInt64(binaryOp, lhs.Value(), rhs.Value())
	if err != nil {
		v.SetError(err)
	}
	dst.SetValue(result)
}
