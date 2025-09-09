package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/objects"
)

func init() {
	SequencerRegister(NewOpDefer)
}

type OpDefer struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

func NewOpDefer(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpDefer{
		Opcode: op.Opcode(bytecode.OpDefer),
		vm:     vmT,
	}, nil
}

func (op *OpDefer) Execute(_ *core.Decoder) {
	// Operands Offset  0
	obj := op.vm.Stack().Pop()
	switch objT := obj.(type) {
	case *objects.FuncCompiled:
		op.vm.Frame().DeferredAdd(objT)
	default:
		op.vm.SetError(fmt.Errorf("invalid operation: defer %s", obj.TypeName()))
		return
	}
}
