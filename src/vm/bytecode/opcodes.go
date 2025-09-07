package bytecode

// Opcodes is a collection that manages and organizes Opcode instances, providing methods to create, retrieve, or compile them.
type Opcodes struct {
	container []*Opcode
	mask      int
}

var noOperands []int

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
	op.createOpcode(OpConstant, []int{Uint16Size}, "OpConstant", OpRelocatable)
	op.createOpcode(OpImport, []int{Uint16Size}, "OpImport", OpRelocatable)
	op.createOpcode(OpClosure, []int{Uint8Size, Uint16Size}, "OpClosure", OpRelocatable) //OpRelocatableFree)
	op.createOpcodeRepeated(OpCallImportGlobal, 16, Uint16Size, "OpCallImportGlobal", OpRelocatable)
	op.createOpcode(OpPop, noOperands, "OpPop", OpRelocatableNone)
	op.createOpcode(OpTrue, noOperands, "OpTrue", OpRelocatableNone)
	op.createOpcode(OpFalse, noOperands, "OpFalse", OpRelocatableNone)
	op.createOpcode(OpBitwiseComplement, noOperands, "OpBitwiseComplement", OpRelocatableNone)
	op.createOpcode(OpMinus, noOperands, "OpMinus", OpRelocatableNone)
	op.createOpcode(OpNot, noOperands, "OpNot", OpRelocatableNone)
	op.createOpcode(OpJumpFalsy, []int{Uint16Size}, "OpJumpFalsy", OpRelocatableNone)
	op.createOpcode(OpJumpTruthy, []int{Uint16Size}, "OpJumpTruthy", OpRelocatableNone)
	op.createOpcode(OpJumpAnd, []int{Uint16Size}, "OpJumpAnd", OpRelocatableNone)
	op.createOpcode(OpJumpOr, []int{Uint16Size}, "OpJumpOr", OpRelocatableNone)
	op.createOpcode(OpJump, []int{Uint16Size}, "OpJump", OpRelocatableNone)
	op.createOpcode(OpJumpNotError, []int{Uint16Size}, "OpJumpNotError", OpRelocatableNone)
	op.createOpcode(OpJumpIndirect, noOperands, "OpJumpIndirect", OpRelocatableNone)
	op.createOpcode(OpNull, noOperands, "OpNull", OpRelocatableNone)
	op.createOpcode(OpGlobalGet, []int{Uint16Size}, "OpGlobalGet", OpRelocatableNone)
	op.createOpcode(OpGlobalSet, []int{Uint16Size}, "OpGlobalSet", OpRelocatableNone)
	op.createOpcode(OpGlobalSelSet, []int{Uint16Size, Uint8Size}, "OpGlobalSelSet", OpRelocatableNone)
	op.createOpcode(OpGlobalCopy, []int{Uint16Size, Uint16Size}, "OpGlobalCopy", OpRelocatableNone)
	op.createOpcode(OpArray, []int{Uint16Size}, "OpArray", OpRelocatableNone)
	op.createOpcode(OpMap, []int{Uint16Size}, "OpMap", OpRelocatableNone)
	op.createOpcode(OpStruct, []int{Uint16Size}, "OpStruct", OpRelocatableNone)
	op.createOpcode(OpInterface, []int{Uint8Size}, "OpInterface", OpRelocatableNone)
	op.createOpcode(OpIndexGet, noOperands, "OpIndexGet", OpRelocatableNone)
	op.createOpcode(OpIndexSet, noOperands, "OpIndexSet", OpRelocatableNone)
	op.createOpcode(OpIndexSlice, noOperands, "OpIndexSlice", OpRelocatableNone)
	op.createOpcode(OpCall, []int{Uint8Size, Uint8Size}, "OpCall", OpRelocatableNone)
	op.createOpcode(OpCallMethod, []int{Uint16Size, Uint8Size}, "OpCallMethod", OpRelocatableNone)
	op.createOpcode(OpReturn, []int{Uint8Size}, "OpReturn", OpRelocatableNone)
	op.createOpcode(OpLocalGet, []int{Uint8Size}, "OpLocalGet", OpRelocatableNone)
	op.createOpcode(OpLocalSet, []int{Uint8Size}, "OpLocalSet", OpRelocatableNone)
	op.createOpcode(OpLocalDefine, []int{Uint8Size}, "OpLocalDefine", OpRelocatableNone)
	op.createOpcode(OpLocalPtrGet, []int{Uint8Size}, "OpLocalPtrGet", OpRelocatableNone)
	op.createOpcode(OpLocalSelSet, []int{Uint8Size, Uint8Size}, "OpLocalSelSet", OpRelocatableNone)
	op.createOpcode(OpFreeGet, []int{Uint8Size}, "OpFreeGet", OpRelocatableNone)
	op.createOpcode(OpFreeSet, []int{Uint8Size}, "OpFreeSet", OpRelocatableNone)
	op.createOpcode(OpFreePtrGet, []int{Uint8Size}, "OpFreePtrGet", OpRelocatableNone)
	op.createOpcode(OpIteratorInit, []int{Uint8Size}, "OpIteratorInit", OpRelocatableNone)
	op.createOpcode(OpIteratorNext, []int{Uint8Size}, "OpIteratorNext", OpRelocatableNone)
	op.createOpcode(OpIteratorKey, []int{Uint8Size}, "OpIteratorKey", OpRelocatableNone)
	op.createOpcode(OpIteratorValue, []int{Uint8Size}, "OpIteratorValue", OpRelocatableNone)
	op.createOpcode(OpLogical, []int{Uint8Size}, "OpLogical", OpRelocatableNone)
	op.createOpcode(OpArithmetic, []int{Uint8Size}, "OpArithmetic", OpRelocatableNone)
	op.createOpcode(OpIntLogical, []int{Uint16Size, Uint16Size, Uint16Size, Uint8Size}, "OpIntLogical", OpRelocatableNone)
	op.createOpcode(OpIntArithmetic, []int{Uint16Size, Uint16Size, Uint16Size, Uint8Size}, "OpIntArithmetic", OpRelocatableNone)
	op.createOpcode(OpDerefGet, noOperands, "OpDerefGet", OpRelocatableNone)
	op.createOpcode(OpDerefSet, noOperands, "OpDerefSet", OpRelocatableNone)
	op.createOpcode(OpTypeAssert, []int{Uint16Size}, "OpTypeAssert", OpRelocatableNone)
	op.createOpcode(OpIsType, []int{Uint16Size}, "OpIsType", OpRelocatableNone)
	op.createOpcode(OpAsType, []int{Uint16Size}, "OpAsType", OpRelocatableNone)
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
	compiler := NewCompiler()
	if err := compiler.Compile(opcode, operands); err != nil {
		return nil, err
	}
	return compiler.Instructions(), nil
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
func (op *Opcodes) createOpcodeRepeated(opcodeId OpcodeId, count int, operand int, name string, relocatable Relocatable) {
	var container []int
	if count > 0 {
		container = make([]int, count)
		for i := 0; i < count; i++ {
			container[i] = operand
		}
	}
	op.createOpcode(opcodeId, container, name, relocatable)
}

// createOpcode registers a new Opcode in the Opcodes container with its identifier, operands, and name.
func (op *Opcodes) createOpcode(opcodeId OpcodeId, operands []int, name string, relocatable Relocatable) {
	od := NewOpcode(opcodeId, operands, name, relocatable)
	op.container[int(od.OpcodeId())&op.mask] = od
}
