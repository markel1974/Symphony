package native

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/handler"
	"github.com/markel1974/symphony/src/vm/opcodes"
)

// _noOperands is a predefined empty slice of OperandFeature used for opcodes that do not require any operands.
var _noOperands []opcodes.OperandFeature

// registerFunc defines a function type that returns a handler.IOpExecutor instance used to register operation executors.
type registerFunc = func() handler.IOpExecutor

// _registerContainer is a slice of registerFunc used to store functions that register IOpExecutor instances.
var _registerContainer []registerFunc

// SequencerRegister appends the provided registerFunc to the internal registration container for later execution.
func SequencerRegister(fn registerFunc) {
	_registerContainer = append(_registerContainer, fn)
}

// Sequencer is responsible for managing and resolving opcode executors within a virtual machine execution context.
type Sequencer struct {
	mask            int
	defaultExecutor []handler.IOpExecutor
	unknown         handler.IOpExecutor
	unknownId       opcodes.OpcodeId
}

// NewSequencer initializes and returns a new instance of the Sequencer struct.
// It sets up opcode executors, calculates mask bits, and registers executors based on available opcodes.
func NewSequencer() *Sequencer {
	unknown := NewOpUnknown()
	seq := &Sequencer{
		unknown:   unknown,
		unknownId: unknown.Opcode().OpcodeId(),
		mask:      0,
	}
	return seq
}

// Setup initializes the Sequencer by registering and configuring opcode executors, ensuring no duplicate opcodes exist.
func (ds *Sequencer) Setup() error {
	maxId := -1
	container := make([]handler.IOpExecutor, len(_registerContainer))
	for idx, fn := range _registerContainer {
		executor := fn()
		if opcodeId := executor.Opcode().OpcodeId(); opcodeId > maxId {
			maxId = opcodeId
		}
		container[idx] = executor
	}
	maskBits := 0
	for (1 << maskBits) <= maxId {
		maskBits++
	}
	ds.mask = (1 << maskBits) - 1
	ds.defaultExecutor = make([]handler.IOpExecutor, ds.mask+1)
	for idx := range ds.defaultExecutor {
		ds.defaultExecutor[idx] = ds.unknown
	}
	for _, executor := range container {
		target := ds.defaultExecutor[executor.Opcode().OpcodeId()]
		if target.Opcode().OpcodeId() != OpUnknownId {
			return fmt.Errorf("duplicate opcode registration: %s", target.Opcode().Name())
		}
		ds.defaultExecutor[executor.Opcode().OpcodeId()] = executor
	}
	return nil
}

// Opcode retrieves the Opcode instance corresponding to the given opcodeId from the Sequencer's executors list.
func (ds *Sequencer) Opcode(opcodeId opcodes.OpcodeId) *opcodes.Opcode {
	return ds.defaultExecutor[opcodeId&ds.mask].Opcode()
}

// Bytecode generates a bytecode representation for a given opcode and its operands, or returns an error if compilation fails.
func (ds *Sequencer) Bytecode(opcodeId opcodes.OpcodeId, operands ...int) ([]byte, error) {
	opcode := ds.Opcode(opcodeId)
	if opcode.OpcodeId() == ds.unknownId {
		return nil, fmt.Errorf("compile: Unknown opcode: %d", opcodeId)
	}
	inst, err := opcode.Compile(operands)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// Mask returns the bitmask used to determine the executor index for opcode handling in the Sequencer.
func (ds *Sequencer) Mask() int {
	return ds.mask
}

// Len returns the total number of executors available in the Sequencer.
func (ds *Sequencer) Len() int {
	return len(ds.defaultExecutor)
}

// Executors retrieves a list of opcode executors and the bitmask used to index them in the Sequencer.
func (ds *Sequencer) Executors() ([]handler.IOpExecutor, int) {
	executors := make([]handler.IOpExecutor, len(ds.defaultExecutor))
	for idx := range executors {
		executors[idx] = ds.unknown
	}
	for _, fn := range _registerContainer {
		executor := fn()
		executors[executor.Opcode().OpcodeId()] = executor
	}
	return executors, ds.mask
}

// Id returns the identifier string for the Sequencer instance, typically used to describe its type or implementation.
func (ds *Sequencer) Id() string {
	return "native"
}
