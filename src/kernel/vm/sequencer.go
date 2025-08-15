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
func (ds *Sequencer) Create() []func(vm *VM) {
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

	sequencer := make([]func(vm *VM), sequenceLen)
	for idx := range sequencer {
		sequencer[idx] = opUnknown.Execute
	}

	sequencer[bytecode.OpConstant] = opConstant.Execute
	sequencer[bytecode.OpNull] = opNull.Execute
	sequencer[bytecode.OpBinaryOp] = opBinary.Execute
	sequencer[bytecode.OpReferences] = opReferences.Execute
	sequencer[bytecode.OpEqual] = opEqual.Execute
	sequencer[bytecode.OpNotEqual] = opNotEqual.Execute
	sequencer[bytecode.OpPop] = opPop.Execute
	sequencer[bytecode.OpTrue] = opTrue.Execute
	sequencer[bytecode.OpFalse] = opFalse.Execute
	sequencer[bytecode.OpLNot] = opLNot.Execute
	sequencer[bytecode.OpBComplement] = opBComplement.Execute
	sequencer[bytecode.OpMinus] = opMinus.Execute
	sequencer[bytecode.OpJumpFalsy] = opJumpFalsy.Execute
	sequencer[bytecode.OpAndJump] = opAndJump.Execute
	sequencer[bytecode.OpOrJump] = opOrJump.Execute
	sequencer[bytecode.OpJump] = opJump.Execute
	sequencer[bytecode.OpSetGlobal] = opSetGlobal.Execute
	sequencer[bytecode.OpSetSelGlobal] = opSetSelGlobal.Execute
	sequencer[bytecode.OpGetGlobal] = opGetGlobal.Execute
	sequencer[bytecode.OpArray] = opArray.Execute
	sequencer[bytecode.OpMap] = opMap.Execute
	sequencer[bytecode.OpError] = opError.Execute
	sequencer[bytecode.OpImmutable] = opImmutable.Execute
	sequencer[bytecode.OpIndex] = opIndex.Execute
	sequencer[bytecode.OpSliceIndex] = opSliceIndex.Execute
	sequencer[bytecode.OpCall] = opCall.Execute
	sequencer[bytecode.OpReturn] = opReturn.Execute
	sequencer[bytecode.OpDefineLocal] = opDefineLocal.Execute
	sequencer[bytecode.OpSetLocal] = opSetLocal.Execute
	sequencer[bytecode.OpSetSelLocal] = opSetSelLocal.Execute
	sequencer[bytecode.OpGetLocal] = opGetLocal.Execute
	sequencer[bytecode.OpGetBuiltin] = opGetBuiltin.Execute
	sequencer[bytecode.OpClosure] = opClosure.Execute
	sequencer[bytecode.OpGetFreePtr] = opGetFreePtr.Execute
	sequencer[bytecode.OpGetFree] = opGetFree.Execute
	sequencer[bytecode.OpSetFree] = opSetFree.Execute
	sequencer[bytecode.OpGetLocalPtr] = opGetLocalPtr.Execute
	sequencer[bytecode.OpSetSelFree] = opSetSelFree.Execute
	sequencer[bytecode.OpIteratorInit] = opIteratorInit.Execute
	sequencer[bytecode.OpIteratorNext] = opIteratorNext.Execute
	sequencer[bytecode.OpIteratorKey] = opIteratorKey.Execute
	sequencer[bytecode.OpIteratorValue] = opIteratorValue.Execute
	sequencer[bytecode.OpSuspend] = opSuspend.Execute
	return sequencer

	/*

		sequencer := make([]func(), sequenceLen)
		for idx := range sequencer {
			sequencer[idx] = v.doOpUnknown
		}
		sequencer[bytecode.OpConstant] = v.doOpConstant
		sequencer[bytecode.OpNull] = v.doOpNull
		sequencer[bytecode.OpBinaryOp] = v.doOpBinary
		sequencer[bytecode.OpReferences] = v.doOpReferences
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

	*/
}
