package core

import (
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

type ISequencer interface {
	opcodes.IOpcodes

	Bind(vm *VM) error

	Executors() []IOpExecutor
}

// IOpExecutor defines the interface for executing virtual machine operations, including binding, execution, and compilation.
// Bind associates the executor with a virtual machine instance.
// Opcode retrieves the Opcode instance associated with the operation.
// Execute performs the execution logic using a decoder for operand processing.
// Compile compiles the operation into assembly and returns the result.
type IOpExecutor interface {
	Bind(vm IVM) error

	Opcode() *opcodes.Opcode

	Execute(decoder *Decoder)

	Compile() ([]byte, error)
}
