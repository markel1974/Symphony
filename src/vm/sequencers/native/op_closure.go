package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpClosure)
}

// OpClosure represents a closure operation that creates a new closure in the virtual machine.
type OpClosure struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpClosure returns a new instance of OpClosure initialized with the details of the OpClosure opcode.
func NewOpClosure(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpClosure{
		Opcode: op.Opcode(bytecode.OpClosure),
		vm:     vmT,
	}, nil
}

// Execute performs the operation associated with the OpClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpClosure) Execute(decoder *core.Decoder) {
	// Operands Offset 3 (8-bit|16-bit)
	numTotal := decoder.Read(0)
	closureIndex := decoder.Read(1)
	closureObj := op.vm.Constants().Get(uint(closureIndex))
	fn, ok := closureObj.(*objects.FuncCompiled)
	if !ok {
		op.vm.SetError(fmt.Errorf("not a function: %s", fn.TypeName()))
		return
	}
	freeArgs := op.vm.Stack().Pop()
	freeIndices, ok := freeArgs.(*objects.Array)
	if !ok {
		op.vm.SetError(fmt.Errorf("invalid operation: cannot create closure without arguments"))
		return
	}
	free := make([]*objects.ObjectPointer, freeIndices.Length())
	for idx, freeObjIndex := range freeIndices.Values() {
		freeIndex, ok := freeObjIndex.(*objects.Int)
		if !ok {
			op.vm.SetError(fmt.Errorf("invalid operation: cannot create closure without arguments"))
			return
		}
		offset := numTotal - int(freeIndex.Value())
		objOffset := op.vm.Stack().PeekOffset(-offset - 1)
		switch objType := objOffset.(type) {
		case *objects.ObjectPointer:
			free[idx] = objType
		default:
			obj := op.vm.Factory().NewObjectPointer(op.vm.Frame().Id(), &objOffset)
			freeObjPtr, ok := obj.(*objects.ObjectPointer)
			if !ok {
				op.vm.SetError(fmt.Errorf("not a pointer: %s", obj.TypeName()))
				return
			}
			free[idx] = freeObjPtr
		}
	}
	op.vm.Stack().DecrementCount(freeIndices.Length())
	cl := op.vm.Factory().NewFuncCompiled(op.vm.Frame().Id(), fn.Name(), fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	op.vm.Stack().Push(cl)
}
