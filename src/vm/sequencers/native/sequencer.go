package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/core"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// _noOperands is a predefined empty slice of OperandFeature used for opcodes that do not require any operands.
var _noOperands []opcodes.OperandFeature

// registerFunc defines a function type that returns a core.IOpExecutor instance used to register operation executors.
type registerFunc = func() core.IOpExecutor

// _registerContainer is a slice of registerFunc used to store functions that register IOpExecutor instances.
var _registerContainer []registerFunc

// SequencerRegister appends the provided registerFunc to the internal registration container for later execution.
func SequencerRegister(fn registerFunc) {
	_registerContainer = append(_registerContainer, fn)
}

// Sequencer is responsible for managing and resolving opcode executors within a virtual machine execution context.
type Sequencer struct {
	mask      int
	executors []core.IOpExecutor
	unknownId opcodes.OpcodeId
}

// NewSequencer initializes and returns a new instance of the Sequencer struct.
// It sets up opcode executors, calculates mask bits, and registers executors based on available opcodes.
func NewSequencer() *Sequencer {
	seq := &Sequencer{
		unknownId: -1,
		mask:      0,
	}
	return seq
}

// Setup initializes the Sequencer by registering and configuring opcode executors, ensuring no duplicate opcodes exist.
func (ds *Sequencer) Setup() error {
	var container []core.IOpExecutor
	maxId := -1
	for _, fn := range _registerContainer {
		executor := fn()
		if opcodeId := executor.Opcode().OpcodeId(); opcodeId > maxId {
			maxId = opcodeId
		}
		container = append(container, executor)
	}
	maskBits := 0
	for (1 << maskBits) <= maxId {
		maskBits++
	}
	unknown := NewOpUnknown()
	ds.unknownId = unknown.Opcode().OpcodeId()
	ds.mask = (1 << maskBits) - 1
	ds.executors = make([]core.IOpExecutor, ds.mask+1)
	for idx := range ds.executors {
		ds.executors[idx] = unknown
	}
	for _, executor := range container {
		target := ds.executors[executor.Opcode().OpcodeId()]
		if target.Opcode().OpcodeId() != OpUnknownId {
			return fmt.Errorf("duplicate opcode registration: %s", target.Opcode().Name())
		}
		ds.executors[executor.Opcode().OpcodeId()] = executor
	}
	return nil
}

// Bind links all opcode executors within the Sequencer to the provided virtual machine instance, returning an error if any fail.
func (ds *Sequencer) Bind(vm *core.VM) error {
	for _, executor := range ds.executors {
		if err := executor.Bind(vm); err != nil {
			return err
		}
	}
	return nil
}

// Opcode retrieves the Opcode instance corresponding to the given opcodeId from the Sequencer's executors list.
func (ds *Sequencer) Opcode(opcodeId opcodes.OpcodeId) *opcodes.Opcode {
	return ds.executors[opcodeId&ds.mask].Opcode()
}

// Compile generates a bytecode representation for a given opcode and its operands, or returns an error if compilation fails.
func (ds *Sequencer) Compile(opcodeId opcodes.OpcodeId, operands ...int) ([]byte, error) {
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
	return len(ds.executors)
}

// Executors returns the slice of IOpExecutor instances managed by the Sequencer.
func (ds *Sequencer) Executors() []core.IOpExecutor {
	return ds.executors
}

// Id returns the identifier string for the Sequencer instance, typically used to describe its type or implementation.
func (ds *Sequencer) Id() string {
	return "native"
}
