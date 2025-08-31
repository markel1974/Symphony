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
func (op *OpClosure) Execute(
	decoder *core.Decoder) {
	// Operands Offset 3 (8-bit|16-bit)
	numFree := decoder.Read(0)
	constIndex := decoder.Read(1)
	glObj := op.vm.Constants().Get(uint(constIndex))
	fn, ok := glObj.(*objects.FuncCompiled)
	if !ok {
		op.vm.SetError(fmt.Errorf("not a function: %s", fn.TypeName()))
		return
	}

	free := make([]*objects.ObjectPointer, numFree)
	for i := 0; i < numFree; i++ {
		offset := numFree - i
		objOffset := op.vm.Stack().PeekOffset(-offset)
		switch objType := objOffset.(type) {
		case *objects.ObjectPointer:
			free[i] = objType
		default:
			obj := op.vm.Factory().NewObjectPointer(op.vm.Frame().Id(), &objOffset)
			freeObjPtr, ok := obj.(*objects.ObjectPointer)
			if !ok {
				op.vm.SetError(fmt.Errorf("not a pointer: %s", obj.TypeName()))
				return
			}
			free[i] = freeObjPtr
		}
	}
	op.vm.Stack().DecrementCount(numFree)
	cl := op.vm.Factory().NewFuncCompiled(op.vm.Frame().Id(), "closure", fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	op.vm.Stack().Push(cl)
}
