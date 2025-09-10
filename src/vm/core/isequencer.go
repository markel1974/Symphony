package core

import (
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// IOpExecutor defines an interface for executing specific bytecode instructions within a virtual machine context.
// Opcode returns the bytecode.OpcodeId associated with the operation.
// Name returns the name of the operation represented by the executor.
// Operands retrieves the operands required for the operation's execution.
// Execute performs the operation within the provided virtual machine instance.
type IOpExecutor interface {
	Bind(vm IVM) error

	Opcode() *opcodes.Opcode

	Execute(decoder *Decoder)
}
