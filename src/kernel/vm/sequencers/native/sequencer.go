package native

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// registerFunc is a function type that registers an operation with the provided bytecode.Opcodes and returns a core.IOpExecutor.
type registerFunc = func(op *bytecode.Opcodes) core.IOpExecutor

// _registerContainer holds a list of functions that register operational executors for bytecode instructions.
var _registerContainer []registerFunc

// SequencerRegister appends a registerFunc to the internal _registerContainer for further use by the sequencer system.
func SequencerRegister(fn registerFunc) {
	_registerContainer = append(_registerContainer, fn)
}

// Sequencer defines a container for managing and initializing opcode executors for a virtual machine's instruction set.
type Sequencer struct {
	op *bytecode.Opcodes
}

// NewSequencer initializes and returns a new Sequencer instance with the provided Opcodes configuration.
func NewSequencer(op *bytecode.Opcodes) *Sequencer {
	return &Sequencer{
		op: op,
	}
}

// Create initializes and returns a slice of IOpExecutor with default OpUnknown executors, supplemented by static mappings.
func (ds *Sequencer) Create() []core.IOpExecutor {
	v, _ := ds.createRegistered()
	return v
}

// createRegistered constructs and registers custom IOpExecutors using functions from _registerContainer and updates the sequence.
func (ds *Sequencer) createRegistered() ([]core.IOpExecutor, error) {
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		container[idx] = NewOpUnknown(ds.op)
	}
	for _, fn := range _registerContainer {
		seq := fn(ds.op)
		data := container[seq.Opcode()&bytecode.OpcodesMask]
		if data.Opcode() != bytecode.OpUnknown {
			return nil, fmt.Errorf("opcode %d already registered: %s", seq.Opcode(), data.Name())
		}
		container[seq.Opcode()&bytecode.OpcodesMask] = seq
	}
	return container, nil
}

// createStatic initializes and assigns specific op executors to the sequencer's container in a pre-defined sequence.
func (ds *Sequencer) createStatic() []core.IOpExecutor {
	container := make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range container {
		container[idx] = NewOpUnknown(ds.op)
	}
	ds.setSequence(container[:], NewOpConstant(ds.op))
	ds.setSequence(container[:], NewOpNull(ds.op))
	ds.setSequence(container[:], NewOpBinary(ds.op))
	ds.setSequence(container[:], NewOpReferences(ds.op))
	ds.setSequence(container[:], NewOpEqual(ds.op))
	ds.setSequence(container[:], NewOpNotEqual(ds.op))
	ds.setSequence(container[:], NewOpPop(ds.op))
	ds.setSequence(container[:], NewOpTrue(ds.op))
	ds.setSequence(container[:], NewOpFalse(ds.op))
	ds.setSequence(container[:], NewOpNotLogical(ds.op))
	ds.setSequence(container[:], NewOpBitwiseComplement(ds.op))
	ds.setSequence(container[:], NewOpMinus(ds.op))
	ds.setSequence(container[:], NewOpJumpFalsy(ds.op))
	ds.setSequence(container[:], NewOpJumpAnd(ds.op))
	ds.setSequence(container[:], NewOpJumpOr(ds.op))
	ds.setSequence(container[:], NewOpJump(ds.op))
	ds.setSequence(container[:], NewOpGlobalSet(ds.op))
	ds.setSequence(container[:], NewOpGlobalSelSet(ds.op))
	ds.setSequence(container[:], NewOpGetGlobal(ds.op))
	ds.setSequence(container[:], NewOpArray(ds.op))
	ds.setSequence(container[:], NewOpMap(ds.op))
	ds.setSequence(container[:], NewOpStruct(ds.op))
	ds.setSequence(container[:], NewOpError(ds.op))
	ds.setSequence(container[:], NewOpImmutable(ds.op))
	ds.setSequence(container[:], NewOpIndex(ds.op))
	ds.setSequence(container[:], NewOpIndexSlice(ds.op))
	ds.setSequence(container[:], NewOpCall(ds.op))
	ds.setSequence(container[:], NewOpReturn(ds.op))
	ds.setSequence(container[:], NewOpLocalDefine(ds.op))
	ds.setSequence(container[:], NewOpLocalSet(ds.op))
	ds.setSequence(container[:], NewOpLocalSelSet(ds.op))
	ds.setSequence(container[:], NewOpLocalGet(ds.op))
	ds.setSequence(container[:], NewOpClosure(ds.op))
	ds.setSequence(container[:], NewOpFreeGetPtr(ds.op))
	ds.setSequence(container[:], NewOpFreeGet(ds.op))
	ds.setSequence(container[:], NewOpFreeSet(ds.op))
	ds.setSequence(container[:], NewOpLocalPtrGet(ds.op))
	ds.setSequence(container[:], NewOpFreeSelSet(ds.op))
	ds.setSequence(container[:], NewOpIteratorInit(ds.op))
	ds.setSequence(container[:], NewOpIteratorNext(ds.op))
	ds.setSequence(container[:], NewOpIteratorKey(ds.op))
	ds.setSequence(container[:], NewOpIteratorValue(ds.op))
	ds.setSequence(container[:], NewOpIntOp(ds.op))
	ds.setSequence(container[:], NewOpDeref(ds.op))
	ds.setSequence(container[:], NewOpNoOp(ds.op))
	ds.setSequence(container[:], NewOpSuspend(ds.op))
	return container
}

// setSequence assigns the given IOpExecutor to the Sequencer's container, using the bit-masked opcode as the index.
func (ds *Sequencer) setSequence(container []core.IOpExecutor, seq core.IOpExecutor) {
	container[seq.Opcode()&bytecode.OpcodesMask] = seq
}
