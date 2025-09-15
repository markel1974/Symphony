package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpFreeGet)
}

// OpFreeGet represents an operation to retrieve a free variable in a closure during execution.
type OpFreeGet struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpFreeGet creates and returns a new instance of OpFreeGet, initializing its Opcode using bytecode metadata.
func NewOpFreeGet() core.IOpExecutor {
	operands := []opcodes.OperandFeature{opcodes.SzUint16}
	return &OpFreeGet{
		opcode: opcodes.NewOpcode(OpFreeGetId, operands, "OpFreeGet"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpFreeGet) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpFreeGet) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpFreeGet) Execute(decoder *core.Decoder) {
	// Operands Offset 2 (16-bit)
	freeIndex := decoder.Operand(0)
	freeVar := op.vm.FrameFreeVarsIndex(uint(freeIndex))
	if freeVar == nil {
		op.vm.SetError(fmt.Errorf("free variable %d not found", freeIndex))
		return
	}
	z := *freeVar.Value()
	op.vm.StackPush(z)
}

// Compile generates the compiled representation of the OpFreeGet operation or returns an unimplemented error.
func (op *OpFreeGet) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
