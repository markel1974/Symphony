package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

func init() {
	SequencerRegister(NewOpClosure)
}

// OpClosure represents a closure operation that creates a new closure in the virtual machine.
type OpClosure struct {
	*bytecode.Opcode
}

// NewOpClosure returns a new instance of OpClosure initialized with the details of the OpClosure opcode.
func NewOpClosure(op *bytecode.Opcodes) core.IOpExecutor {
	return &OpClosure{Opcode: op.Opcode(bytecode.OpClosure)}
}

// Execute performs the operation associated with the OpClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpClosure) Execute(v *core.VM, decoder *core.Decoder) {
	// Operands Offset 3 (8-bit|16-bit)
	numFree := decoder.Read(0)
	constIndex := decoder.Read(1)
	glObj := v.Constants().Get(uint(constIndex))
	fn, ok := glObj.(*objects.FuncCompiled)
	if !ok {
		v.SetError(fmt.Errorf("not a function: %s", fn.TypeName()))
		return
	}
	free := make([]*objects.ObjectPointer, numFree)
	for i := 0; i < numFree; i++ {
		o := v.Stack().PeekOffset(-numFree + i)
		switch freeVar := o.(type) {
		case *objects.ObjectPointer:
			free[i] = freeVar
		default:
			t := v.Stack().PeekOffset(-numFree + i)
			obj := v.Factory().NewObjectPointer(v.Frame().Id(), &t)
			ptr, ok := obj.(*objects.ObjectPointer)
			if !ok {
				v.SetError(fmt.Errorf("not a pointer: %s", t.TypeName()))
				return
			}
			free[i] = ptr
		}
	}
	v.Stack().DecrementCount(numFree)
	cl := v.Factory().NewFuncCompiled(v.Frame().Id(), "closure", fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	v.Stack().Push(cl)
}
