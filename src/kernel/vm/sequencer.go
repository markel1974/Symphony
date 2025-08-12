package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
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
	sequencer[bytecode.OpConstant] = v.doOpConstant
	sequencer[bytecode.OpNull] = v.doOpNull
	sequencer[bytecode.OpBinaryOp] = v.doOpBinary
	sequencer[bytecode.OpEqual] = v.doOpEqual
	sequencer[bytecode.OpNotEqual] = v.doOpNotEqual
	sequencer[bytecode.OpPop] = v.doOpPop
	sequencer[bytecode.OpTrue] = v.doOpTrue
	sequencer[bytecode.OpFalse] = v.doOpFalse
	sequencer[bytecode.OpLNot] = v.doOpLNot
	sequencer[bytecode.OpBComplement] = v.doOpBComplement
	sequencer[bytecode.OpMinus] = v.doOpMinus
	sequencer[bytecode.OpJumpFalsy] = v.doOpJumpFalsy
	sequencer[bytecode.OpAndJump] = v.doOpAndJump
	sequencer[bytecode.OpOrJump] = v.doOpOrJump
	sequencer[bytecode.OpJump] = v.doOpJump
	sequencer[bytecode.OpSetGlobal] = v.doOpSetGlobal
	sequencer[bytecode.OpSetSelGlobal] = v.doOpSetSelGlobal
	sequencer[bytecode.OpGetGlobal] = v.doOpGetGlobal
	sequencer[bytecode.OpArray] = v.doOpArray
	sequencer[bytecode.OpMap] = v.doOpMap
	sequencer[bytecode.OpError] = v.doOpError
	sequencer[bytecode.OpImmutable] = v.doOpImmutable
	sequencer[bytecode.OpIndex] = v.doOpIndex
	sequencer[bytecode.OpSliceIndex] = v.doOpSliceIndex
	sequencer[bytecode.OpCall] = v.doOpCall
	sequencer[bytecode.OpReturn] = v.doOpReturn
	sequencer[bytecode.OpDefineLocal] = v.doOpDefineLocal
	sequencer[bytecode.OpSetLocal] = v.doOpSetLocal
	sequencer[bytecode.OpSetSelLocal] = v.doOpSetSelLocal
	sequencer[bytecode.OpGetLocal] = v.doOpGetLocal
	sequencer[bytecode.OpGetBuiltin] = v.doOpGetBuiltin
	sequencer[bytecode.OpClosure] = v.doOpClosure
	sequencer[bytecode.OpGetFreePtr] = v.doOpGetFreePtr
	sequencer[bytecode.OpGetFree] = v.doOpGetFree
	sequencer[bytecode.OpSetFree] = v.doOpSetFree
	sequencer[bytecode.OpGetLocalPtr] = v.doOpGetLocalPtr
	sequencer[bytecode.OpSetSelFree] = v.doOpSetSelFree
	sequencer[bytecode.OpIteratorInit] = v.doOpIteratorInit
	sequencer[bytecode.OpIteratorNext] = v.doOpIteratorNext
	sequencer[bytecode.OpIteratorKey] = v.doOpIteratorKey
	sequencer[bytecode.OpIteratorValue] = v.doOpIteratorValue
	sequencer[bytecode.OpSuspend] = v.doOpSuspend
	return sequencer
}
