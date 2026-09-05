package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

func init() {
	SequencerRegister(NewOpMakeChan)
}

// OpMakeChan represents a bytecode operation for creating a channel in the virtual machine.
type OpMakeChan struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpMakeChan creates and returns a new instance of OpMakeChan.
func NewOpMakeChan() handler.IOpExecutor {
	return &OpMakeChan{
		opcode: opcodes.NewOpcode(OpMakeChanId, _noOperands, "OpMakeChan"),
		vm:     nil,
	}
}

func (op *OpMakeChan) Opcode() *opcodes.Opcode {
	return op.opcode
}

func (op *OpMakeChan) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

func (op *OpMakeChan) Execute(decoder *handler.Decoder) {
	capObj := op.vm.StackPop()
	capacity := int(capObj.AsInt64())
	if capacity < 0 {
		op.vm.Shutdown(fmt.Errorf("make chan: invalid capacity %d", capacity))
		return
	}
	chObj := op.vm.Factory().NewChan(op.vm.FrameId(), capacity)
	op.vm.StackPush(chObj)
}

func (op *OpMakeChan) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
