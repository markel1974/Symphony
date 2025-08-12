package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecodes"
)

// Sequencer is a type that provides functionality to create a sequence of operations for a virtual machine.
type Sequencer struct {
}

func NewSequencer() *Sequencer {
	return &Sequencer{}
}

// Create initializes a sequencer array with operation functions mapped to their corresponding opcodes.
func (ds *Sequencer) Create(v *VM) []func() {
	sequencer := make([]func(), sequenceLen)
	for idx := range sequencer {
		sequencer[idx] = v.doOpUnknown
	}
	sequencer[bytecodes.OpConstant] = v.doOpConstant
	sequencer[bytecodes.OpNull] = v.doOpNull
	sequencer[bytecodes.OpBinaryOp] = v.doOpBinary
	sequencer[bytecodes.OpEqual] = v.doOpEqual
	sequencer[bytecodes.OpNotEqual] = v.doOpNotEqual
	sequencer[bytecodes.OpPop] = v.doOpPop
	sequencer[bytecodes.OpTrue] = v.doOpTrue
	sequencer[bytecodes.OpFalse] = v.doOpFalse
	sequencer[bytecodes.OpLNot] = v.doOpLNot
	sequencer[bytecodes.OpBComplement] = v.doOpBComplement
	sequencer[bytecodes.OpMinus] = v.doOpMinus
	sequencer[bytecodes.OpJumpFalsy] = v.doOpJumpFalsy
	sequencer[bytecodes.OpAndJump] = v.doOpAndJump
	sequencer[bytecodes.OpOrJump] = v.doOpOrJump
	sequencer[bytecodes.OpJump] = v.doOpJump
	sequencer[bytecodes.OpSetGlobal] = v.doOpSetGlobal
	sequencer[bytecodes.OpSetSelGlobal] = v.doOpSetSelGlobal
	sequencer[bytecodes.OpGetGlobal] = v.doOpGetGlobal
	sequencer[bytecodes.OpArray] = v.doOpArray
	sequencer[bytecodes.OpMap] = v.doOpMap
	sequencer[bytecodes.OpError] = v.doOpError
	sequencer[bytecodes.OpImmutable] = v.doOpImmutable
	sequencer[bytecodes.OpIndex] = v.doOpIndex
	sequencer[bytecodes.OpSliceIndex] = v.doOpSliceIndex
	sequencer[bytecodes.OpCall] = v.doOpCall
	sequencer[bytecodes.OpReturn] = v.doOpReturn
	sequencer[bytecodes.OpDefineLocal] = v.doOpDefineLocal
	sequencer[bytecodes.OpSetLocal] = v.doOpSetLocal
	sequencer[bytecodes.OpSetSelLocal] = v.doOpSetSelLocal
	sequencer[bytecodes.OpGetLocal] = v.doOpGetLocal
	sequencer[bytecodes.OpGetBuiltin] = v.doOpGetBuiltin
	sequencer[bytecodes.OpClosure] = v.doOpClosure
	sequencer[bytecodes.OpGetFreePtr] = v.doOpGetFreePtr
	sequencer[bytecodes.OpGetFree] = v.doOpGetFree
	sequencer[bytecodes.OpSetFree] = v.doOpSetFree
	sequencer[bytecodes.OpGetLocalPtr] = v.doOpGetLocalPtr
	sequencer[bytecodes.OpSetSelFree] = v.doOpSetSelFree
	sequencer[bytecodes.OpIteratorInit] = v.doOpIteratorInit
	sequencer[bytecodes.OpIteratorNext] = v.doOpIteratorNext
	sequencer[bytecodes.OpIteratorKey] = v.doOpIteratorKey
	sequencer[bytecodes.OpIteratorValue] = v.doOpIteratorValue
	sequencer[bytecodes.OpSuspend] = v.doOpSuspend
	return sequencer
}
