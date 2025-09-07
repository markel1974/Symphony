package bytecode

// Opcodes is a collection that manages and organizes Opcode instances, providing methods to create, retrieve, or compile them.
type Opcodes struct {
	container []*Opcode
	mask      int
}

var noOperands []OperandFeature

// NewOpcodes initializes and returns a new Opcodes instance with predefined opcode mappings.
func NewOpcodes() *Opcodes {
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
		op.container[i] = NewOpcode(OpUnknown, noOperands, "OpUnknown", OpRelocatableNone)
	}
	op.createOpcode(OpConstant, []OperandFeature{Size2}, "OpConstant", OpRelocatable)
	op.createOpcode(OpImport, []OperandFeature{Size2}, "OpImport", OpRelocatable)
	op.createOpcode(OpClosure, []OperandFeature{Size1, Size2}, "OpClosure", OpRelocatable) //OpRelocatableFree)
	op.createOpcodeRepeated(OpCallImportGlobal, 16, Size2, "OpCallImportGlobal", OpRelocatable)
	op.createOpcode(OpPop, noOperands, "OpPop", OpRelocatableNone)
	op.createOpcode(OpTrue, noOperands, "OpTrue", OpRelocatableNone)
	op.createOpcode(OpFalse, noOperands, "OpFalse", OpRelocatableNone)
	op.createOpcode(OpBitwiseComplement, noOperands, "OpBitwiseComplement", OpRelocatableNone)
	op.createOpcode(OpMinus, noOperands, "OpMinus", OpRelocatableNone)
	op.createOpcode(OpNot, noOperands, "OpNot", OpRelocatableNone)
	op.createOpcode(OpJumpFalsy, []OperandFeature{Size2}, "OpJumpFalsy", OpRelocatableNone)
	op.createOpcode(OpJumpTruthy, []OperandFeature{Size2}, "OpJumpTruthy", OpRelocatableNone)
	op.createOpcode(OpJumpAnd, []OperandFeature{Size2}, "OpJumpAnd", OpRelocatableNone)
	op.createOpcode(OpJumpOr, []OperandFeature{Size2}, "OpJumpOr", OpRelocatableNone)
	op.createOpcode(OpJump, []OperandFeature{Size2}, "OpJump", OpRelocatableNone)
	op.createOpcode(OpJumpNotError, []OperandFeature{Size2}, "OpJumpNotError", OpRelocatableNone)
	op.createOpcode(OpJumpIndirect, noOperands, "OpJumpIndirect", OpRelocatableNone)
	op.createOpcode(OpNull, noOperands, "OpNull", OpRelocatableNone)
	op.createOpcode(OpGlobalGet, []OperandFeature{Size2}, "OpGlobalGet", OpRelocatableNone)
	op.createOpcode(OpGlobalSet, []OperandFeature{Size2}, "OpGlobalSet", OpRelocatableNone)
	op.createOpcode(OpGlobalSelSet, []OperandFeature{Size2, Size1}, "OpGlobalSelSet", OpRelocatableNone)
	op.createOpcode(OpGlobalCopy, []OperandFeature{Size2, Size2}, "OpGlobalCopy", OpRelocatableNone)
	op.createOpcode(OpArray, []OperandFeature{Size2}, "OpArray", OpRelocatableNone)
	op.createOpcode(OpMap, []OperandFeature{Size2}, "OpMap", OpRelocatableNone)
	op.createOpcode(OpStruct, []OperandFeature{Size2}, "OpStruct", OpRelocatableNone)
	op.createOpcode(OpInterface, []OperandFeature{Size1}, "OpInterface", OpRelocatableNone)
	op.createOpcode(OpIndexGet, noOperands, "OpIndexGet", OpRelocatableNone)
	op.createOpcode(OpIndexSet, noOperands, "OpIndexSet", OpRelocatableNone)
	op.createOpcode(OpIndexSlice, noOperands, "OpIndexSlice", OpRelocatableNone)
	op.createOpcode(OpCall, []OperandFeature{Size1, Size1}, "OpCall", OpRelocatableNone)
	op.createOpcode(OpCallMethod, []OperandFeature{Size2, Size1}, "OpCallMethod", OpRelocatableNone)
	op.createOpcode(OpReturn, []OperandFeature{Size1}, "OpReturn", OpRelocatableNone)
	op.createOpcode(OpLocalGet, []OperandFeature{Size1}, "OpLocalGet", OpRelocatableNone)
	op.createOpcode(OpLocalSet, []OperandFeature{Size1}, "OpLocalSet", OpRelocatableNone)
	op.createOpcode(OpLocalDefine, []OperandFeature{Size1}, "OpLocalDefine", OpRelocatableNone)
	op.createOpcode(OpLocalPtrGet, []OperandFeature{Size1}, "OpLocalPtrGet", OpRelocatableNone)
	op.createOpcode(OpLocalSelSet, []OperandFeature{Size1, Size1}, "OpLocalSelSet", OpRelocatableNone)
	op.createOpcode(OpFreeGet, []OperandFeature{Size1}, "OpFreeGet", OpRelocatableNone)
	op.createOpcode(OpFreeSet, []OperandFeature{Size1}, "OpFreeSet", OpRelocatableNone)
	op.createOpcode(OpFreePtrGet, []OperandFeature{Size1}, "OpFreePtrGet", OpRelocatableNone)
	op.createOpcode(OpIteratorInit, []OperandFeature{Size1}, "OpIteratorInit", OpRelocatableNone)
	op.createOpcode(OpIteratorNext, []OperandFeature{Size1}, "OpIteratorNext", OpRelocatableNone)
	op.createOpcode(OpIteratorKey, []OperandFeature{Size1}, "OpIteratorKey", OpRelocatableNone)
	op.createOpcode(OpIteratorValue, []OperandFeature{Size1}, "OpIteratorValue", OpRelocatableNone)
	op.createOpcode(OpLogical, []OperandFeature{Size1}, "OpLogical", OpRelocatableNone)
	op.createOpcode(OpArithmetic, []OperandFeature{Size1}, "OpArithmetic", OpRelocatableNone)
	op.createOpcode(OpIntLogical, []OperandFeature{Size2, Size2, Size2, Size1}, "OpIntLogical", OpRelocatableNone)
	op.createOpcode(OpIntArithmetic, []OperandFeature{Size2, Size2, Size2, Size1}, "OpIntArithmetic", OpRelocatableNone)
	op.createOpcode(OpDerefGet, noOperands, "OpDerefGet", OpRelocatableNone)
	op.createOpcode(OpDerefSet, noOperands, "OpDerefSet", OpRelocatableNone)
	op.createOpcode(OpTypeAssert, []OperandFeature{Size2}, "OpTypeAssert", OpRelocatableNone)
	op.createOpcode(OpIsType, []OperandFeature{Size2}, "OpIsType", OpRelocatableNone)
	op.createOpcode(OpAsType, []OperandFeature{Size2}, "OpAsType", OpRelocatableNone)
	op.createOpcode(OpSuspend, noOperands, "OpSuspend", OpRelocatableNone)
	op.createOpcode(OpError, noOperands, "OpError", OpRelocatableNone)
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

// createOpcodeRepeated initializes and registers an opcode with a repeated operand pattern and a specified name.
func (op *Opcodes) createOpcodeRepeated(opcodeId OpcodeId, count int, operand OperandFeature, name string, relocatable Relocatable) {
	var container []OperandFeature
	if count > 0 {
		container = make([]OperandFeature, count)
		for i := 0; i < count; i++ {
			container[i] = operand
		}
	}
	op.createOpcode(opcodeId, container, name, relocatable)
}

// createOpcode registers a new Opcode in the Opcodes container with its identifier, operands, and name.
func (op *Opcodes) createOpcode(opcodeId OpcodeId, operands []OperandFeature, name string, relocatable Relocatable) {
	od := NewOpcode(opcodeId, operands, name, relocatable)
	op.container[int(od.OpcodeId())&op.mask] = od
}
