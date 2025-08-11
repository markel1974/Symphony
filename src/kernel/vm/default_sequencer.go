package vm

import "github.com/markel1974/c64emu/src/kernel/vm/opcodes"

// DefaultSequencer is a type that provides functionality to create a sequence of operations for a virtual machine.
type DefaultSequencer struct {
}

// Create initializes a sequencer array with operation functions mapped to their corresponding opcodes.
func (ds *DefaultSequencer) Create(v *VM) []func() {
	sequencer := make([]func(), sequenceLen)
	for idx := range sequencer {
		sequencer[idx] = v.doOpUnknown
	}
	sequencer[opcodes.OpConstant] = v.doOpConstant
	sequencer[opcodes.OpNull] = v.doOpNull
	sequencer[opcodes.OpBinaryOp] = v.doOpBinary
	sequencer[opcodes.OpEqual] = v.doOpEqual
	sequencer[opcodes.OpNotEqual] = v.doOpNotEqual
	sequencer[opcodes.OpPop] = v.doOpPop
	sequencer[opcodes.OpTrue] = v.doOpTrue
	sequencer[opcodes.OpFalse] = v.doOpFalse
	sequencer[opcodes.OpLNot] = v.doOpLNot
	sequencer[opcodes.OpBComplement] = v.doOpBComplement
	sequencer[opcodes.OpMinus] = v.doOpMinus
	sequencer[opcodes.OpJumpFalsy] = v.doOpJumpFalsy
	sequencer[opcodes.OpAndJump] = v.doOpAndJump
	sequencer[opcodes.OpOrJump] = v.doOpOrJump
	sequencer[opcodes.OpJump] = v.doOpJump
	sequencer[opcodes.OpSetGlobal] = v.doOpSetGlobal
	sequencer[opcodes.OpSetSelGlobal] = v.doOpSetSelGlobal
	sequencer[opcodes.OpGetGlobal] = v.doOpGetGlobal
	sequencer[opcodes.OpArray] = v.doOpArray
	sequencer[opcodes.OpMap] = v.doOpMap
	sequencer[opcodes.OpError] = v.doOpError
	sequencer[opcodes.OpImmutable] = v.doOpImmutable
	sequencer[opcodes.OpIndex] = v.doOpIndex
	sequencer[opcodes.OpSliceIndex] = v.doOpSliceIndex
	sequencer[opcodes.OpCall] = v.doOpCall
	sequencer[opcodes.OpReturn] = v.doOpReturn
	sequencer[opcodes.OpDefineLocal] = v.doOpDefineLocal
	sequencer[opcodes.OpSetLocal] = v.doOpSetLocal
	sequencer[opcodes.OpSetSelLocal] = v.doOpSetSelLocal
	sequencer[opcodes.OpGetLocal] = v.doOpGetLocal
	sequencer[opcodes.OpGetBuiltin] = v.doOpGetBuiltin
	sequencer[opcodes.OpClosure] = v.doOpClosure
	sequencer[opcodes.OpGetFreePtr] = v.doOpGetFreePtr
	sequencer[opcodes.OpGetFree] = v.doOpGetFree
	sequencer[opcodes.OpSetFree] = v.doOpSetFree
	sequencer[opcodes.OpGetLocalPtr] = v.doOpGetLocalPtr
	sequencer[opcodes.OpSetSelFree] = v.doOpSetSelFree
	sequencer[opcodes.OpIteratorInit] = v.doOpIteratorInit
	sequencer[opcodes.OpIteratorNext] = v.doOpIteratorNext
	sequencer[opcodes.OpIteratorKey] = v.doOpIteratorKey
	sequencer[opcodes.OpIteratorValue] = v.doOpIteratorValue
	sequencer[opcodes.OpSuspend] = v.doOpSuspend
	return sequencer
}
