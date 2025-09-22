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
	operands := []opcodes.OperandFeature{opcodes.Relocatable}
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
	closureFnIndex := decoder.Operand(0)
	closureFnObj := op.vm.Constants().Get(uint(closureFnIndex))
	fn, ok := closureFnObj.(*objects.Func)
	if !ok {
		op.vm.SetError(fmt.Errorf("not a function: %s", closureFnObj.TypeName()))
		return
	}
	var required []objects.IObject
	if reqIndices := fn.FreeIndices(); len(reqIndices) > 0 {
		required = make([]objects.IObject, len(reqIndices))
		for idx, freeObjIndex := range reqIndices {
			obj := op.vm.StackPeekBP(uint(freeObjIndex))
			required[idx] = obj
		}
	}
	cl := op.vm.CreateClosure(fn, required)
	op.vm.StackPush(cl)
}

// Compile generates the compiled representation of the OpCreateClosure operation or returns an unimplemented error.
func (op *OpCreateClosure) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
