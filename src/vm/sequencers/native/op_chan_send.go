package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the OpChanSend operation executor to the sequencer using SequencerRegister.
func init() {
	SequencerRegister(NewOpChanSend)
}

// OpChanSend represents an operation that sends a value to a channel object in the virtual machine.
type OpChanSend struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpChanSend creates a new instance of OpChanSend and returns it as an IOpExecutor implementation.
func NewOpChanSend() handler.IOpExecutor {
	return &OpChanSend{
		opcode: opcodes.NewOpcode(OpChanSendId, _noOperands, "OpChanSend"),
		vm:     nil,
	}
}

// Opcode returns the Opcode instance associated with the current OpChanSend operation.
func (op *OpChanSend) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind associates the provided IVM instance with the OpChanSend instance, requiring it to implement IVMFullAccess.
func (op *OpChanSend) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the send operation to a channel using the values from the stack. It handles synchronization and blocking.
func (op *OpChanSend) Execute(decoder *handler.Decoder) {
	val := op.vm.StackPop()
	chObj := op.vm.StackPop()

	ch, ok := chObj.(*objects.Chan)
	if !ok {
		op.vm.Shutdown(fmt.Errorf("send to non-chan type"))
		return
	}
	data := ch.Data()

	if wakeId, wakeOk := data.GetRecv(); wakeOk {
		data.AddBuffer(val)
		op.vm.WakeCore(wakeId)
		return
	}

	// If there's room in the buffer, just append and proceed.
	if data.AddBuffer(val) {
		return
	}

	// Otherwise, we must block. Restore the stack and IP to retry this instruction later.
	op.vm.StackPush(chObj)
	op.vm.StackPush(val)
	op.vm.SetIp(op.vm.GetIp() - 1)
	op.vm.BlockCurrentCore()
	data.AddSend(op.vm.CoreId())
}

// Compile generates bytecode for the operation, returning an error if the functionality is unimplemented.
func (op *OpChanSend) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
