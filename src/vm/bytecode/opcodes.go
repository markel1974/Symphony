package bytecode

// Opcodes is a collection that manages and organizes Opcode instances, providing methods to create, retrieve, or compile them.
type Opcodes struct {
	container []*Opcode
	mask      int
}

// NewOpcodes initializes and returns a new Opcodes instance with predefined opcode mappings.
func NewOpcodes() *Opcodes {
	var noOperands []OperandFeature
	maskBits := 0
	for (1 << maskBits) <= int(OpUnknown) {
		maskBits++
	}
	mask := (1 << maskBits) - 1
	op := &Opcodes{
		container: make([]*Opcode, mask+1),
		mask:      mask,
	}
	for i := range op.container {
		op.container[i] = NewOpcode(OpUnknown, noOperands, "OpUnknown")
	}
	op.create(OpConstant, []OperandFeature{Relocatable}, "OpConstant")
	op.create(OpImport, []OperandFeature{Relocatable}, "OpImport")
	op.create(OpClosure, []OperandFeature{SzUint8, Relocatable}, "OpClosure") //OpRelocatableFree)
	op.create(OpCallImportGlobal, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, Relocatable}, "OpCallImportGlobal")
	op.create(OpGlobalDefine, []OperandFeature{Relocatable}, "OpGlobalDefine")
	op.create(OpGlobalGet, []OperandFeature{Relocatable}, "OpGlobalGet")
	op.create(OpGlobalSet, []OperandFeature{Relocatable}, "OpGlobalSet")
	op.create(OpGlobalCopy, []OperandFeature{Relocatable, Relocatable}, "OpGlobalCopy")
	op.create(OpGlobalIndex, []OperandFeature{SzUint8, Relocatable}, "OpGlobalIndex")
	op.create(OpLocalIndex, []OperandFeature{SzUint8, Relocatable}, "OpLocalIndex")
	op.create(OpPop, noOperands, "OpPop")
	op.create(OpTrue, noOperands, "OpTrue")
	op.create(OpFalse, noOperands, "OpFalse")
	op.create(OpBitwiseComplement, noOperands, "OpBitwiseComplement")
	op.create(OpMinus, noOperands, "OpMinus")
	op.create(OpNot, noOperands, "OpNot")
	op.create(OpJumpFalsy, []OperandFeature{SzUint16}, "OpJumpFalsy")
	op.create(OpJumpTruthy, []OperandFeature{SzUint16}, "OpJumpTruthy")
	op.create(OpJumpAnd, []OperandFeature{SzUint16}, "OpJumpAnd")
	op.create(OpJumpOr, []OperandFeature{SzUint16}, "OpJumpOr")
	op.create(OpJump, []OperandFeature{SzUint16}, "OpJump")
	op.create(OpJumpNotError, []OperandFeature{SzUint16}, "OpJumpNotError")
	op.create(OpJumpIndirect, noOperands, "OpJumpIndirect")
	op.create(OpNull, noOperands, "OpNull")
	op.create(OpArray, []OperandFeature{SzUint16}, "OpArray")
	op.create(OpMap, []OperandFeature{SzUint16}, "OpMap")
	op.create(OpStruct, []OperandFeature{SzUint16}, "OpStruct")
	op.create(OpInterface, []OperandFeature{SzUint8}, "OpInterface")
	op.create(OpIndexGet, noOperands, "OpIndexGet")
	op.create(OpIndexSet, noOperands, "OpIndexSet")
	op.create(OpIndexSlice, noOperands, "OpIndexSlice")
	op.create(OpCall, []OperandFeature{SzUint8, SzUint8}, "OpCall")
	op.create(OpCallMethod, []OperandFeature{SzUint16, SzUint8}, "OpCallMethod")
	op.create(OpReturn, []OperandFeature{SzUint8}, "OpReturn")
	op.create(OpLocalGet, []OperandFeature{SzUint8}, "OpLocalGet")
	op.create(OpLocalSet, []OperandFeature{SzUint8}, "OpLocalSet")
	op.create(OpLocalDefine, []OperandFeature{SzUint8}, "OpLocalDefine")
	op.create(OpLocalPtrGet, []OperandFeature{SzUint8}, "OpLocalPtrGet")
	op.create(OpFreeGet, []OperandFeature{SzUint8}, "OpFreeGet")
	op.create(OpFreeSet, []OperandFeature{SzUint8}, "OpFreeSet")
	op.create(OpFreePtrGet, []OperandFeature{SzUint8}, "OpFreePtrGet")
	op.create(OpIteratorInit, []OperandFeature{SzUint8}, "OpIteratorInit")
	op.create(OpIteratorNext, []OperandFeature{SzUint8}, "OpIteratorNext")
	op.create(OpIteratorKey, []OperandFeature{SzUint8}, "OpIteratorKey")
	op.create(OpIteratorValue, []OperandFeature{SzUint8}, "OpIteratorValue")
	op.create(OpLogical, []OperandFeature{SzUint8}, "OpLogical")
	op.create(OpArithmetic, []OperandFeature{SzUint8}, "OpArithmetic")
	op.create(OpIntLogical, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint8}, "OpIntLogical")
	op.create(OpIntArithmetic, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint8}, "OpIntArithmetic")
	op.create(OpDerefGet, noOperands, "OpDerefGet")
	op.create(OpDerefSet, noOperands, "OpDerefSet")
	op.create(OpTypeAssert, []OperandFeature{SzUint16}, "OpTypeAssert")
	op.create(OpIsType, []OperandFeature{SzUint16}, "OpIsType")
	op.create(OpAsType, []OperandFeature{SzUint16}, "OpAsType")
	op.create(OpSuspend, noOperands, "OpSuspend")
	op.create(OpError, noOperands, "OpError")
	return op
}

// Opcode retrieves the Opcode associated with the given OpcodeId and returns it along with a boolean indicating success.
func (op *Opcodes) Opcode(opcodeId OpcodeId) *Opcode {
	return op.container[int(opcodeId)&op.mask]
}

// CompileInstruction generates bytecode for a given opcode and its operands or returns an error if validation fails.
func (op *Opcodes) CompileInstruction(opcodeId OpcodeId, operands ...int) ([]byte, error) {
	opcode := op.Opcode(opcodeId)
	inst, err := opcode.Compile(operands)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// Mask returns the opcode mask used to determine the valid range or format of opcodes in the Opcodes collection.
func (op *Opcodes) Mask() int {
	return op.mask
}

// Len returns the number of elements in the Opcodes container.
func (op *Opcodes) Len() int {
	return len(op.container)
}

// create registers a new Opcode in the Opcodes container with its identifier, operands, and name.
func (op *Opcodes) create(opcodeId OpcodeId, operands []OperandFeature, name string) {
	od := NewOpcode(opcodeId, operands, name)
	op.container[int(od.OpcodeId())&op.mask] = od
}
