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
	OpcodeId() opcodes.OpcodeId

	Name() string

	Operands() []opcodes.OperandFeature

	Execute(decoder *Decoder)
}

// ISequencer defines an interface to generate a sequence of functions for a given Virtual Machine instance.
type ISequencer interface {
	Create(vm *VM) ([]IOpExecutor, error)
}
