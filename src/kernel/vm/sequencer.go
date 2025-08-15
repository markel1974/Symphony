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
func (ds *Sequencer) Create() []IOpExecutor {
	opUnknown := &OpUnknown{}
	opConstant := &OpConstant{}
	opNull := &OpNull{}
	opBinary := &OpBinary{}
	opReferences := &OpReferences{}
	opEqual := &OpEqual{}
	opNotEqual := &OpNotEqual{}
	opPop := &OpPop{}
	opTrue := &OpTrue{}
	opFalse := &OpFalse{}
	opLNot := &OpLNot{}
	opBComplement := &OpBComplement{}
	opMinus := &OpMinus{}
	opJumpFalsy := &OpJumpFalsy{}
	opAndJump := &OpAndJump{}
	opOrJump := &OpOrJump{}
	opJump := &OpJump{}
	opSetGlobal := &OpSetGlobal{}
	opSetSelGlobal := &OpSetSelGlobal{}
	opGetGlobal := &OpGetGlobal{}
	opArray := &OpArray{}
	opMap := &OpMap{}
	opError := &OpError{}
	opImmutable := &OpImmutable{}
	opIndex := &OpIndex{}
	opSliceIndex := &OpSliceIndex{}
	opCall := &OpCall{}
	opReturn := &OpReturn{}
	opDefineLocal := &OpDefineLocal{}
	opSetLocal := &OpSetLocal{}
	opSetSelLocal := &OpSetSelLocal{}
	opGetLocal := &OpGetLocal{}
	opGetBuiltin := &OpGetBuiltin{}
	opClosure := &OpClosure{}
	opGetFreePtr := &OpGetFreePtr{}
	opGetFree := &OpGetFree{}
	opSetFree := &OpSetFree{}
	opGetLocalPtr := &OpGetLocalPtr{}
	opSetSelFree := &OpSetSelFree{}
	opIteratorInit := &OpIteratorInit{}
	opIteratorNext := &OpIteratorNext{}
	opIteratorKey := &OpIteratorKey{}
	opIteratorValue := &OpIteratorValue{}
	opSuspend := &OpSuspend{}

	sequencer := make([]IOpExecutor, sequenceLen)
	for idx := range sequencer {
		sequencer[idx] = opUnknown
	}
	sequencer[bytecode.OpConstant] = opConstant
	sequencer[bytecode.OpNull] = opNull
	sequencer[bytecode.OpBinaryOp] = opBinary
	sequencer[bytecode.OpReferences] = opReferences
	sequencer[bytecode.OpEqual] = opEqual
	sequencer[bytecode.OpNotEqual] = opNotEqual
	sequencer[bytecode.OpPop] = opPop
	sequencer[bytecode.OpTrue] = opTrue
	sequencer[bytecode.OpFalse] = opFalse
	sequencer[bytecode.OpLNot] = opLNot
	sequencer[bytecode.OpBComplement] = opBComplement
	sequencer[bytecode.OpMinus] = opMinus
	sequencer[bytecode.OpJumpFalsy] = opJumpFalsy
	sequencer[bytecode.OpAndJump] = opAndJump
	sequencer[bytecode.OpOrJump] = opOrJump
	sequencer[bytecode.OpJump] = opJump
	sequencer[bytecode.OpSetGlobal] = opSetGlobal
	sequencer[bytecode.OpSetSelGlobal] = opSetSelGlobal
	sequencer[bytecode.OpGetGlobal] = opGetGlobal
	sequencer[bytecode.OpArray] = opArray
	sequencer[bytecode.OpMap] = opMap
	sequencer[bytecode.OpError] = opError
	sequencer[bytecode.OpImmutable] = opImmutable
	sequencer[bytecode.OpIndex] = opIndex
	sequencer[bytecode.OpSliceIndex] = opSliceIndex
	sequencer[bytecode.OpCall] = opCall
	sequencer[bytecode.OpReturn] = opReturn
	sequencer[bytecode.OpDefineLocal] = opDefineLocal
	sequencer[bytecode.OpSetLocal] = opSetLocal
	sequencer[bytecode.OpSetSelLocal] = opSetSelLocal
	sequencer[bytecode.OpGetLocal] = opGetLocal
	sequencer[bytecode.OpGetBuiltin] = opGetBuiltin
	sequencer[bytecode.OpClosure] = opClosure
	sequencer[bytecode.OpGetFreePtr] = opGetFreePtr
	sequencer[bytecode.OpGetFree] = opGetFree
	sequencer[bytecode.OpSetFree] = opSetFree
	sequencer[bytecode.OpGetLocalPtr] = opGetLocalPtr
	sequencer[bytecode.OpSetSelFree] = opSetSelFree
	sequencer[bytecode.OpIteratorInit] = opIteratorInit
	sequencer[bytecode.OpIteratorNext] = opIteratorNext
	sequencer[bytecode.OpIteratorKey] = opIteratorKey
	sequencer[bytecode.OpIteratorValue] = opIteratorValue
	sequencer[bytecode.OpSuspend] = opSuspend
	return sequencer
}
