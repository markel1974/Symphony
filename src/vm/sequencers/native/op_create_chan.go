package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/objects"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// init registers the operation executor for creating channels with the SequencerRegister during initialization.
func init() {
	SequencerRegister(NewOpCreateChan)
}

// OpCreateChan represents an operation for creating a channel with a specified capacity in the virtual machine.
type OpCreateChan struct {
	opcode *opcodes.Opcode
	vm     handler.IVMFullAccess
}

// NewOpCreateChan initializes and returns a new instance of OpCreateChan, implementing the IOpExecutor interface.
func NewOpCreateChan() handler.IOpExecutor {
	return &OpCreateChan{
		opcode: opcodes.NewOpcode(OpCreateChanId, _noOperands, "OpCreateChan"),
		vm:     nil,
	}
}

// Opcode returns the opcode associated with the OpCreateChan instance.
func (op *OpCreateChan) Opcode() *opcodes.Opcode {
	return op.opcode
}

// Bind attempts to bind the provided IVM instance to the operation, requiring it to implement IVMFullAccess.
func (op *OpCreateChan) Bind(vm handler.IVM) error {
	vmT, ok := vm.(handler.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute processes the operation to create a channel with a specified capacity and pushes it onto the stack.
func (op *OpCreateChan) Execute(decoder *handler.Decoder) {
	capObj := op.vm.StackPop()
	capacity := int(capObj.AsInt64())
	if capacity < 0 {
		op.vm.Shutdown(fmt.Errorf("create chan: invalid capacity %d", capacity))
		return
	}
	chObj := op.vm.Factory().NewChan(op.vm.FrameId(), capacity)
	op.vm.StackPush(chObj)
}

// Compile translates the operation into a sequence of bytes or returns an error if the operation is unimplemented.
func (op *OpCreateChan) Compile() ([]byte, error) {
	return nil, objects.ErrUnimplemented
}
