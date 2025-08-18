package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
)

// Sequencer is a type that manages a collection of IOpExecutor instances organized by their opcodes.
// It provides methods for creating and populating the sequencer with specific opcode implementations.
// The container field stores the IOpExecutor instances, indexed by their opcode with masking applied for efficiency.
type Sequencer struct {
	container []IOpExecutor
}

// NewSequencer initializes and returns a new instance of Sequencer.
func NewSequencer() *Sequencer {
	return &Sequencer{}
}

// Create initializes the internal container with default operation executors and returns it.
func (ds *Sequencer) Create() []IOpExecutor {
	ds.container = make([]IOpExecutor, bytecode.OpcodesLen)
	for idx := range ds.container {
		ds.container[idx] = NewOpUnknown()
	}
	ds.setSequence(NewOpConstant())
	ds.setSequence(NewOpNull())
	ds.setSequence(NewOpBinary())
	ds.setSequence(NewOpReferences())
	ds.setSequence(NewOpEqual())
	ds.setSequence(NewOpNotEqual())
	ds.setSequence(NewOpPop())
	ds.setSequence(NewOpTrue())
	ds.setSequence(NewOpFalse())
	ds.setSequence(NewOpLNot())
	ds.setSequence(NewOpBComplement())
	ds.setSequence(NewOpMinus())
	ds.setSequence(NewOpJumpFalsy())
	ds.setSequence(NewOpAndJump())
	ds.setSequence(NewOpOrJump())
	ds.setSequence(NewOpJump())
	ds.setSequence(NewOpSetGlobal())
	ds.setSequence(NewOpSetSelGlobal())
	ds.setSequence(NewOpGetGlobal())
	ds.setSequence(NewOpArray())
	ds.setSequence(NewOpMap())
	ds.setSequence(NewOpStruct())
	ds.setSequence(NewOpError())
	ds.setSequence(NewOpImmutable())
	ds.setSequence(NewOpIndex())
	ds.setSequence(NewOpSliceIndex())
	ds.setSequence(NewOpCall())
	ds.setSequence(NewOpReturn())
	ds.setSequence(NewOpDefineLocal())
	ds.setSequence(NewOpSetLocal())
	ds.setSequence(NewOpSetSelLocal())
	ds.setSequence(NewOpGetLocal())
	ds.setSequence(NewOpClosure())
	ds.setSequence(NewOpGetFreePtr())
	ds.setSequence(NewOpGetFree())
	ds.setSequence(NewOpSetFree())
	ds.setSequence(NewOpGetLocalPtr())
	ds.setSequence(NewOpSetSelFree())
	ds.setSequence(NewOpIteratorInit())
	ds.setSequence(NewOpIteratorNext())
	ds.setSequence(NewOpIteratorKey())
	ds.setSequence(NewOpIteratorValue())
	ds.setSequence(NewOpSuspend())
	return ds.container
}

// setSequence assigns a specific IOpExecutor implementation to its corresponding opcode index in the container.
func (ds *Sequencer) setSequence(seq IOpExecutor) {
	ds.container[seq.Opcode()&bytecode.OpcodesMask] = seq
}
