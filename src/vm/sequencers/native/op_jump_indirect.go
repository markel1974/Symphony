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
	opcode *opcodes.Opcode
	vm     core.IVMFullAccess
}

// NewOpJumpIndirect creates an `OpJumpIndirect` executor for handling indirect jump instructions in the virtual machine.
// It requires a core.IVM instance and a bytecode.Opcodes reference. Returns an error if the vm does not implement IVMFullAccess.
func NewOpJumpIndirect() core.IOpExecutor {
	operands := _noOperands
	return &OpJumpIndirect{
		opcode: opcodes.NewOpcode(OpJumpIndirectId, operands, "OpJumpIndirect"),
		vm:     nil,
	}
}

// Bind initializes the instance by casting the provided VM to IVMFullAccess and storing it.
// Returns an error if the VM does not implement the required interface.
func (op *OpJumpIndirect) Bind(vm core.IVM) error {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return fmt.Errorf("vm does not implement IVMFullAccess")
	}
	op.vm = vmT
	return nil
}

// Execute performs an indirect jump by popping an address from the stack and setting the instruction pointer to it.
func (op *OpJumpIndirect) Execute(decoder *core.Decoder) {
	addrObj := op.vm.Stack().Pop()
	addr := addrObj.AsInt64()
	op.vm.SetIp(int(addr) - 1)
}

// Opcode returns the opcode associated with the instance.
func (op *OpJumpIndirect) Opcode() *opcodes.Opcode {
	return op.opcode
}
