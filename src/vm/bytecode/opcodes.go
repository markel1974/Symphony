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
	op.createOpcode(OpConstant, []OperandFeature{Relocatable}, "OpConstant")
	op.createOpcode(OpImport, []OperandFeature{Relocatable}, "OpImport")
	op.createOpcode(OpClosure, []OperandFeature{SzUint8, Relocatable}, "OpClosure") //OpRelocatableFree)
	op.createOpcode(OpCallImportGlobal, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, Relocatable}, "OpCallImportGlobal")
	op.createOpcode(OpPop, noOperands, "OpPop")
	op.createOpcode(OpTrue, noOperands, "OpTrue")
	op.createOpcode(OpFalse, noOperands, "OpFalse")
	op.createOpcode(OpBitwiseComplement, noOperands, "OpBitwiseComplement")
	op.createOpcode(OpMinus, noOperands, "OpMinus")
	op.createOpcode(OpNot, noOperands, "OpNot")
	op.createOpcode(OpJumpFalsy, []OperandFeature{SzUint16}, "OpJumpFalsy")
	op.createOpcode(OpJumpTruthy, []OperandFeature{SzUint16}, "OpJumpTruthy")
	op.createOpcode(OpJumpAnd, []OperandFeature{SzUint16}, "OpJumpAnd")
	op.createOpcode(OpJumpOr, []OperandFeature{SzUint16}, "OpJumpOr")
	op.createOpcode(OpJump, []OperandFeature{SzUint16}, "OpJump")
	op.createOpcode(OpJumpNotError, []OperandFeature{SzUint16}, "OpJumpNotError")
	op.createOpcode(OpJumpIndirect, noOperands, "OpJumpIndirect")
	op.createOpcode(OpNull, noOperands, "OpNull")
	op.createOpcode(OpGlobalGet, []OperandFeature{SzUint16}, "OpGlobalGet")
	op.createOpcode(OpGlobalSet, []OperandFeature{SzUint16}, "OpGlobalSet")
	op.createOpcode(OpGlobalSelSet, []OperandFeature{SzUint16, SzUint8}, "OpGlobalSelSet")
	op.createOpcode(OpGlobalCopy, []OperandFeature{SzUint16, SzUint16}, "OpGlobalCopy")
	op.createOpcode(OpArray, []OperandFeature{SzUint16}, "OpArray")
	op.createOpcode(OpMap, []OperandFeature{SzUint16}, "OpMap")
	op.createOpcode(OpStruct, []OperandFeature{SzUint16}, "OpStruct")
	op.createOpcode(OpInterface, []OperandFeature{SzUint8}, "OpInterface")
	op.createOpcode(OpIndexGet, noOperands, "OpIndexGet")
	op.createOpcode(OpIndexSet, noOperands, "OpIndexSet")
	op.createOpcode(OpIndexSlice, noOperands, "OpIndexSlice")
	op.createOpcode(OpCall, []OperandFeature{SzUint8, SzUint8}, "OpCall")
	op.createOpcode(OpCallMethod, []OperandFeature{SzUint16, SzUint8}, "OpCallMethod")
	op.createOpcode(OpReturn, []OperandFeature{SzUint8}, "OpReturn")
	op.createOpcode(OpLocalGet, []OperandFeature{SzUint8}, "OpLocalGet")
	op.createOpcode(OpLocalSet, []OperandFeature{SzUint8}, "OpLocalSet")
	op.createOpcode(OpLocalDefine, []OperandFeature{SzUint8}, "OpLocalDefine")
	op.createOpcode(OpLocalPtrGet, []OperandFeature{SzUint8}, "OpLocalPtrGet")
	op.createOpcode(OpLocalSelSet, []OperandFeature{SzUint8, SzUint8}, "OpLocalSelSet")
	op.createOpcode(OpFreeGet, []OperandFeature{SzUint8}, "OpFreeGet")
	op.createOpcode(OpFreeSet, []OperandFeature{SzUint8}, "OpFreeSet")
	op.createOpcode(OpFreePtrGet, []OperandFeature{SzUint8}, "OpFreePtrGet")
	op.createOpcode(OpIteratorInit, []OperandFeature{SzUint8}, "OpIteratorInit")
	op.createOpcode(OpIteratorNext, []OperandFeature{SzUint8}, "OpIteratorNext")
	op.createOpcode(OpIteratorKey, []OperandFeature{SzUint8}, "OpIteratorKey")
	op.createOpcode(OpIteratorValue, []OperandFeature{SzUint8}, "OpIteratorValue")
	op.createOpcode(OpLogical, []OperandFeature{SzUint8}, "OpLogical")
	op.createOpcode(OpArithmetic, []OperandFeature{SzUint8}, "OpArithmetic")
	op.createOpcode(OpIntLogical, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint8}, "OpIntLogical")
	op.createOpcode(OpIntArithmetic, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint8}, "OpIntArithmetic")
	op.createOpcode(OpDerefGet, noOperands, "OpDerefGet")
	op.createOpcode(OpDerefSet, noOperands, "OpDerefSet")
	op.createOpcode(OpTypeAssert, []OperandFeature{SzUint16}, "OpTypeAssert")
	op.createOpcode(OpIsType, []OperandFeature{SzUint16}, "OpIsType")
	op.createOpcode(OpAsType, []OperandFeature{SzUint16}, "OpAsType")
	op.createOpcode(OpSuspend, noOperands, "OpSuspend")
	op.createOpcode(OpError, noOperands, "OpError")
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
func (op *Opcodes) createOpcodeRepeated(opcodeId OpcodeId, count int, operand OperandFeature, name string) {
	var container []OperandFeature
	if count > 0 {
		container = make([]OperandFeature, count)
		for i := 0; i < count; i++ {
			container[i] = operand
		}
	}
	op.createOpcode(opcodeId, container, name)
}

// createOpcode registers a new Opcode in the Opcodes container with its identifier, operands, and name.
func (op *Opcodes) createOpcode(opcodeId OpcodeId, operands []OperandFeature, name string) {
	od := NewOpcode(opcodeId, operands, name)
	op.container[int(od.OpcodeId())&op.mask] = od
}
