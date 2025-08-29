// File: vm/sequencers/native/op_as_type.go

package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// init initializes the package by registering the NewOpAsType operation with the sequencer system.
func init() {
	SequencerRegister(NewOpAsType)
}

// OpAsType represents an executor linked to the bytecode opcode OpAsType for handling unchecked casts in the VM.
// It embeds a bytecode.Opcode and uses core.IVMFullAccess for full VM functionality.
type OpAsType struct {
	*bytecode.Opcode
	vm core.IVMFullAccess
}

// NewOpAsType creates a new instance of OpAsType executor for the given virtual machine and opcode.
// Returns an error if the provided VM does not implement IVMFullAccess.
func NewOpAsType(vm core.IVM, op *bytecode.Opcodes) (core.IOpExecutor, error) {
	vmT, ok := vm.(core.IVMFullAccess)
	if !ok {
		return nil, fmt.Errorf("vm does not implement IVMFullAccess")
	}
	return &OpAsType{
		Opcode: op.Opcode(bytecode.OpAsType),
		vm:     vmT,
	}, nil
}

// Execute performs the operation by popping an interface from the stack and pushing its concrete value back.
// If the popped object is not an interface, it sets an error in the virtual machine.
func (op *OpAsType) Execute(decoder *core.Decoder) {
	// Questo opcode non usa i suoi operandi, ma li manteniamo per coerenza
	// se in futuro servissero per qualche controllo di sicurezza.

	interfaceObj := op.vm.Stack().Pop()

	io, isInterface := interfaceObj.(*objects.Interface)
	if !isInterface {
		// Questo non dovrebbe mai accadere in un type switch valido.
		op.vm.SetError(fmt.Errorf("cannot perform unchecked cast on a non-interface type: %s", interfaceObj.TypeName()))
		return
	}

	// Sostituisce l'interfaccia con il suo valore concreto sullo stack.
	op.vm.Stack().Push(io.Value())
}
