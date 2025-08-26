package executors

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpCall represents an operation code for invoking a function call in the virtual machine.
type OpCall struct {
	*bytecode.OpcodeDetails
}

// NewOpCall creates and returns a new instance of OpCall with initialized OpcodeDetails for the OpCall opcode.
func NewOpCall(op *bytecode.Opcodes) *OpCall {
	return &OpCall{OpcodeDetails: op.OpcodeToDetails(bytecode.OpCall)}
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 2 (8-bit|8-bit)
	spread := decoder.Read(0)
	numArgs := decoder.Read(1)
	value := v.Stack().PeekOffset(-1 - numArgs)
	if !value.CanCall() {
		v.SetError(fmt.Errorf("%s is not callable: %s", value.String(), value.TypeName()))
		return
	}
	if spread == 1 {
		arrObj := v.Stack().Pop()
		switch z := arrObj.(type) {
		case *objects.Array:
			for _, item := range z.Values() {
				v.Stack().Push(item)
			}
			numArgs += z.Length() - 1
		case *objects.ArrayImmutable:
			for _, item := range z.Values() {
				v.Stack().Push(item)
			}
			numArgs += z.Length() - 1
		default:
			v.SetError(fmt.Errorf("not an array: %s", arrObj.TypeName()))
			return
		}
	}

	if callee, ok := value.(*objects.FuncCompiled); ok {
		if callee.VarArgs() {
			v.Stack().PushVarArgs(v.FrameID(), numArgs, callee.NumParameters()-1)
			numArgs = callee.NumParameters()
		}
		if numArgs != callee.NumParameters() {
			numParams := callee.NumParameters()
			if callee.VarArgs() {
				numParams--
			}
			v.SetError(fmt.Errorf("%s wrong number of arguments: want>=%d, got=%d", callee.Name(), numParams, numArgs))
			return
		}
		v.FunctionCompiledCall(callee, numArgs)
	} else {
		var args []objects.IObject
		args = append(args, v.Stack().PeekArrayObject(numArgs)...)
		v.FunctionLibraryCall(value, args, numArgs)
	}
}
