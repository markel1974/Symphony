package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpCreateClosure)
}

// OpCreateClosure represents a closure operation that creates a new closure in the virtual machine.
type OpCreateClosure struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpCreateClosure returns a new instance of OpCreateClosure initialized with the details of the OpCreateClosure opcode.
func NewOpCreateClosure(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpCreateClosure{
		Opcode: op.Opcode(bytecode.OpCreateClosure),
		vm:     vmT,
	}, nil
}

// Execute performs the operation associated with the OpCreateClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpCreateClosure) Execute(decoder *core.Decoder) {
	// Operands Offset 3 (16-bit|8-bit)
	closureIndex := decoder.Read(0)
	numTotal := decoder.Read(1)
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
