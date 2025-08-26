package executors

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/core"
)

// Sequencer is a type that manages a collection of IOpExecutor instances organized by their opcodes.
// It provides methods for creating and populating the sequencer with specific opcode implementations.
// The container field stores the IOpExecutor instances, indexed by their opcode with masking applied for efficiency.
type Sequencer struct {
	op        *bytecode.Opcodes
	container []core.IOpExecutor
}

// NewSequencer initializes and returns a new instance of Sequencer.
func NewSequencer(op *bytecode.Opcodes) *Sequencer {
	return &Sequencer{
		op: op,
	}
}

// Create initializes the internal container with default operation executors and returns it.
func (ds *Sequencer) Create() []core.IOpExecutor {
	ds.container = make([]core.IOpExecutor, bytecode.OpcodesLen)
	for idx := range ds.container {
		ds.container[idx] = NewOpUnknown(ds.op)
	}
	ds.setSequence(NewOpConstant(ds.op))
	ds.setSequence(NewOpNull(ds.op))
	ds.setSequence(NewOpBinary(ds.op))
	ds.setSequence(NewOpReferences(ds.op))
	ds.setSequence(NewOpEqual(ds.op))
	ds.setSequence(NewOpNotEqual(ds.op))
	ds.setSequence(NewOpPop(ds.op))
	ds.setSequence(NewOpTrue(ds.op))
	ds.setSequence(NewOpFalse(ds.op))
	ds.setSequence(NewOpNotLogical(ds.op))
	ds.setSequence(NewOpBitwiseComplement(ds.op))
	ds.setSequence(NewOpMinus(ds.op))
	ds.setSequence(NewOpJumpFalsy(ds.op))
	ds.setSequence(NewOpJumpAnd(ds.op))
	ds.setSequence(NewOpJumpOr(ds.op))
	ds.setSequence(NewOpJump(ds.op))
	ds.setSequence(NewOpGlobalSet(ds.op))
	ds.setSequence(NewOpGlobalSelSet(ds.op))
	ds.setSequence(NewOpGetGlobal(ds.op))
	ds.setSequence(NewOpArray(ds.op))
	ds.setSequence(NewOpMap(ds.op))
	ds.setSequence(NewOpStruct(ds.op))
	ds.setSequence(NewOpError(ds.op))
	ds.setSequence(NewOpImmutable(ds.op))
	ds.setSequence(NewOpIndex(ds.op))
	ds.setSequence(NewOpIndexSlice(ds.op))
	ds.setSequence(NewOpCall(ds.op))
	ds.setSequence(NewOpReturn(ds.op))
	ds.setSequence(NewOpLocalDefine(ds.op))
	ds.setSequence(NewOpLocalSet(ds.op))
	ds.setSequence(NewOpLocalSelSet(ds.op))
	ds.setSequence(NewOpLocalGet(ds.op))
	ds.setSequence(NewOpClosure(ds.op))
	ds.setSequence(NewOpFreeGetPtr(ds.op))
	ds.setSequence(NewOpFreeGet(ds.op))
	ds.setSequence(NewOpFreeSet(ds.op))
	ds.setSequence(NewOpLocalPtrGet(ds.op))
	ds.setSequence(NewOpFreeSelSet(ds.op))
	ds.setSequence(NewOpIteratorInit(ds.op))
	ds.setSequence(NewOpIteratorNext(ds.op))
	ds.setSequence(NewOpIteratorKey(ds.op))
	ds.setSequence(NewOpIteratorValue(ds.op))
	ds.setSequence(NewOpIntOp(ds.op))
	ds.setSequence(NewOpDeref(ds.op))
	ds.setSequence(NewOpSuspend(ds.op))
	return ds.container
}

// setSequence assigns a specific IOpExecutor implementation to its corresponding opcode index in the container.
func (ds *Sequencer) setSequence(seq core.IOpExecutor) {
	ds.container[seq.Opcode()&bytecode.OpcodesMask] = seq
}
