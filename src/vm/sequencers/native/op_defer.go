package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpDefer)
}

type OpDefer struct {
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

func NewOpDefer() core.IOpExecutor {
	operands := _noOperands
	return &OpDefer{
		opcode: opcodes.NewOpcode(OpDeferId, operands, "OpDefer"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpDefer) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

func (op *OpDefer) Execute(_ *core.Decoder) {
	// Operands Offset  0
	obj := op.vm.StackPop()
	switch objT := obj.(type) {
	case *objects.FuncCompiled:
		op.vm.FrameDeferredAdd(objT)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: defer %s", obj.TypeName()))
		return
	}
}

// Opcode returns the opcode associated with the instance.
func (op *OpDefer) Opcode() *opcodes.Opcode {
	return op.opcode
}
