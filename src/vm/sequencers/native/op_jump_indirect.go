package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// init initializes the package by registering the OpJumpIndirect operation with the sequencer system.
func init() {
	SequencerRegister(NewOpJumpIndirect)
}

// OpJumpIndirect is a bytecode operation that performs an indirect jump to an address popped from the VM stack.
// It adjusts the instruction pointer (IP) to the target address minus one after popping it from the stack.
// The instruction supports VM implementations providing full access through the IVMFullAccess interface.
type OpJumpIndirect struct {
	*opcodes.Opcode
	vm core.IVMFullAccess
}

// NewOpJumpIndirect creates an `OpJumpIndirect` executor for handling indirect jump instructions in the virtual machine.
// It requires a core.IVM instance and a bytecode.Opcodes reference. Returns an error if the vm does not implement IVMFullAccess.
func NewOpJumpIndirect(vm core.IVM, op *opcodes.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpJumpIndirect{
		Opcode: op.Opcode(opcodes.OpJumpIndirect),
		vm:     vmT,
	}, nil
}

// Execute performs an indirect jump by popping an address from the stack and setting the instruction pointer to it.
func (op *OpJumpIndirect) Execute(decoder *core.Decoder) {
	addrObj := op.vm.Stack().Pop()
	addr := addrObj.AsInt64()
	op.vm.SetIp(int(addr) - 1)
}
