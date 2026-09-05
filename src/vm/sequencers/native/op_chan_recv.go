package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init initializes the package by registering the NewOpChanRecv operation using SequencerRegister.
func init() {
	SequencerRegister(NewOpChanRecv)
}

// OpChanRecv represents an operation for receiving values from a channel in a virtual machine context.
type OpChanRecv struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpChanRecv initializes and returns an OpChanRecv operation executor for receiving values from a channel.
func NewOpChanRecv() handler.IOpExecutor {
	return &OpChanRecv{
		opcode: opcodes.NewOpcode(OpChanRecvId, _noOperands, "OpChanRecv"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the OpChanRecv instance.
func (op *OpChanRecv) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind links the given virtual machine instance to the OpChanRecv, ensuring it implements IVMFullAccess.
func (op *OpChanRecv) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs the receive operation from a channel, handling blocking, unblocking, and stack state changes.
func (op *OpChanRecv) Execute(decoder *handler.Decoder) {
	chObj := op.vm.StackPop()

	ch, ok := chObj.(*objects.Chan)
	if !ok {
		op.vm.Shutdown(fmt.Errorf("receive from non-chan type"))
		return
	}
	data := ch.Data()

	// Take the value from the buffer
	val := data.GetBuffer()

	if val != nil {
		// If a sender was blocked because the buffer was full (or unbuffered), wake them up.
		if wakeId, wakeOk := data.GetSend(); wakeOk {
			op.vm.WakeCore(wakeId)
		}
		op.vm.StackPush(val)
		return
	}

	// Buffer is empty.
	if wakeId, wakeOk := data.GetSend(); wakeOk {
		op.vm.WakeCore(wakeId)
	}

	// We must block until a value is available.
	op.vm.StackPush(chObj)
	op.vm.SetIp(op.vm.GetIp() - 1)
	op.vm.BlockCurrentCore()
	data.AddRecv(op.vm.CoreId())
}

// Compile translates the operation into bytecode and returns the resulting byte slice or an unimplemented error.
func (op *OpChanRecv) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
