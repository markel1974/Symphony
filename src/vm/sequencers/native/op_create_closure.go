package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpCreateClosure)
}

// OpCreateClosure represents a closure operation that creates a new closure in the virtual machine.
type OpCreateClosure struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCreateClosure returns a new instance of OpCreateClosure initialized with the details of the OpCreateClosure opcode.
func NewOpCreateClosure() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint8, opcodes.Relocatable}
	return &OpCreateClosure{
		opcode: opcodes.NewOpcode(OpCreateClosureId, operands, "OpCreateClosure"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCreateClosure) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCreateClosure) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the operation associated with the OpCreateClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpCreateClosure) Execute(decoder *core.Decoder) {
	closureIndex := decoder.Operand(0)
	//numTotal := decoder.Operand(1)
	closureObj := op.vm.Constants().Get(uint(closureIndex))
	fn, ok := closureObj.(*objects.Func)
	if !ok {
		op.vm.SetError(fmt.Errorf("not a function: %s", fn.TypeName()))
		return
	}
	freeArgs := op.vm.StackPop()
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
		index := int(freeIndex.Value())
		objOffset := op.vm.StackPeekBP(uint(index))
		switch objType := objOffset.(type) {
		case *objects.ObjectPointer:
			free[idx] = objType
		default:
			obj := op.vm.Factory().NewObjectPointer(op.vm.FrameId(), &objOffset)
			freeObjPtr, ok := obj.(*objects.ObjectPointer)
			if !ok {
				op.vm.SetError(fmt.Errorf("not a pointer: %s", obj.TypeName()))
				return
			}
			free[idx] = freeObjPtr
		}
	}
	op.vm.StackDecrementCount(uint(freeIndices.Length()))
	cl := op.vm.Factory().NewFunc(op.vm.FrameId(), fn.Name(), fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	op.vm.StackPush(cl)
}

// Compile generates the compiled representation of the OpCreateClosure operation or returns an unimplemented error.
func (op *OpCreateClosure) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
