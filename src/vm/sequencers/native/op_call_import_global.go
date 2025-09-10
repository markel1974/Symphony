package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init registers the NewOpCallImportGlobal function into the sequencing system using SequencerRegister.
func init() {
	SequencerRegister(NewOpCallImportGlobal)
}

// OpCallImportGlobal represents an operation to call a function from the global import table within the virtual machine.
// It uses a VM with full access permissions and an associated opcode for execution configuration.
type OpCallImportGlobal struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpCallImportGlobal creates a new instance of OpCallImportGlobal, ensuring the VM implements IVMFullAccess.
// Returns the new OpCallImportGlobal instance or an error if VM type assertion fails.
func NewOpCallImportGlobal() core.IOpExecutor {
	operands := []opcodes.OperandFeature{
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.SzUint16,
		opcodes.Relocatable,
	}
	return &OpCallImportGlobal{
		opcode: opcodes.NewOpcode(OpCallImportGlobalId, operands, "OpCallImportGlobal"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpCallImportGlobal) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute decodes and executes the current import-global call instruction using the provided decoder.
func (op *OpCallImportGlobal) Execute(decoder *core.Decoder) {
	funcImportIndex := decoder.Operand(0)
	numArgs := decoder.Operand(1)
	callee := op.vm.Imports().Get(uint(funcImportIndex))
	switch numArgs {
	case 0:
		op.vm.CallObject(callee, 0)
	case 1:
		i2 := op.vm.Globals().Get(uint(decoder.Operand(2)))
		op.vm.CallObject(callee, numArgs, i2)
	case 2:
		i2 := op.vm.Globals().Get(uint(decoder.Operand(2)))
		i3 := op.vm.Globals().Get(uint(decoder.Operand(3)))
		op.vm.CallObject(callee, numArgs, i2, i3)
	case 3:
		i2 := op.vm.Globals().Get(uint(decoder.Operand(2)))
		i3 := op.vm.Globals().Get(uint(decoder.Operand(3)))
		i4 := op.vm.Globals().Get(uint(decoder.Operand(4)))
		op.vm.CallObject(callee, numArgs, i2, i3, i4)
	default:
		args := make([]objects.IObject, numArgs)
		for i := 0; i < numArgs; i++ {
			globalIndex := uint(decoder.Operand(2 + i))
			args[i] = op.vm.Globals().Get(globalIndex)
		}
		op.vm.CallObject(callee, len(args), args...)
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpCallImportGlobal) Opcode() *opcodes.Opcode {
	return op.opcode
}
