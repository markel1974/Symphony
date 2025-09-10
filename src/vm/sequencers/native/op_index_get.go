package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpIndexGet)
}

// OpIndexGet represents the operation for performing an indexing operation on a value.
type OpIndexGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpIndexGet creates and returns a new instance of OpIndexGet initialized with its associated Opcode.
func NewOpIndexGet() core.IOpExecutor {
	operands := _noOperands
	return &OpIndexGet{
		opcode: opcodes.NewOpcode(OpIndexGetId, operands, "OpIndexGet"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpIndexGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndexGet) Execute(_ *core.Decoder) {
	// Operands Offset  0
	index := op.vm.Stack().Pop()
	left := op.vm.Stack().Pop()
	val, err := left.IndexGet(op.vm.Frame().Id(), index)
	if err != nil {
		op.vm.SetError(objects.ComputeIndexGetError(err, index.TypeName(), index.TypeName()))
		return
	}
	if val == nil {
		val = op.vm.Factory().UndefinedValue()
	}
	op.vm.Stack().Push(val)
}

// Opcode returns the opcode associated with the instance.
func (op *OpIndexGet) Opcode() *opcodes.Opcode {
	return op.opcode
}
