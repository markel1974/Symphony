package core

import (
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// ISequencer defines behavior for binding execution contexts and managing operation execution sequences.
// Bind initializes or links the sequencer to a virtual machine instance.
// Executors retrieves a collection of operational executors associated with the sequencer.
type ISequencer interface {
	opcodes.IOpcodes

	Bind(vm *VM) error

	Executors() []IOpExecutor
}

// IOpExecutor defines an interface for executing operations in a virtual machine, including binding, decoding, and compiling tasks.
// Bind associates a virtual machine instance with the executor.
// Opcode retrieves the opcode definition associated with the executor.
// Execute applies the decoding logic to process and execute the instruction.
// Compile compiles the associated operation into a bytecode sequence and returns it or an error if compilation fails.
type IOpExecutor interface {
	Bind(vm IVM) error

	Opcode() *opcodes.Opcode

	Execute(decoder *Decoder)

	Compile() ([]byte, error)
}
