package bytecode

import "fmt"

// OpcodeId is a type alias for byte, used to represent operation codes in instruction sets.
type OpcodeId = byte

const (
	// OpConstant represents an operation that loads a constant value onto the stack.
	OpConstant OpcodeId = iota

	// OpPop is an OpcodeId used to remove and discard the top value from the stack.
	OpPop

	// OpTrue represents the opcode for pushing the boolean value true onto the stack.
	OpTrue

	// OpFalse represents an operation code for pushing the boolean value 'false' onto the stack in the virtual machine.
	OpFalse

	// OpUnarySub represents the opcode for performing unary negation operations.
	OpUnarySub

	// OpUnaryAdd represents an operation for applying a unary addition to a numeric value.
	OpUnaryAdd

	// OpUnaryNot represents the logical NOT (!) operation in the opcode set.
	OpUnaryNot

	// OpUnaryBitwiseComplement represents the opcode for performing a bitwise complement operation on a value.
	OpUnaryBitwiseComplement

	// OpJumpFalsy represents a conditional jump instruction that redirects execution if the top stack value is falsy.
	OpJumpFalsy

	// OpJumpTruthy represents a conditional jump instruction that redirects execution if the top stack value is truthy.
	OpJumpTruthy

	// OpJumpAnd is an opcode used to perform a conditional jump based on the evaluation of a logical AND operation.
	OpJumpAnd

	// OpJumpOr represents an operation code used to perform a conditional jump if a logical OR condition evaluates to true.
	OpJumpOr

	// OpJump is a constant representing an unconditional jump operation in the bytecode instruction set.
	OpJump

	// OpJumpIndirect is a constant representing an opcode for performing an indirect jump. It does not require operands.
	OpJumpIndirect

	// OpNull represents a null operation or a placeholder indicating a null value in the opcode sequence.
	OpNull

	// OpCreateArray represents the opcode for creating a new array with a specified number of elements from the operand.
	OpCreateArray

	// OpCreateMap defines an opcode representing the creation of a map structure with a specified number of key-value pairs.
	OpCreateMap

	// OpCreateStruct represents an opcode for initializing a struct with a specified number of key-value pairs.
	OpCreateStruct

	// OpCreateInterface represents the opcode used to construct an interface object with required method bindings on the stack.
	OpCreateInterface

	// OpJumpNotError handles the 'if err != nil' pattern, skipping the if block if the top stack object is null or not an error.
	OpJumpNotError

	// OpTypeAssert implements type assertion 'val, ok := i.(Type)'.
	OpTypeAssert

	// OpIsType is a helper for type switch, checks type, and pushes a boolean.
	OpIsType

	// OpAsType is a helper for type switch, performs type casting without checks.
	OpAsType

	// OpIndexGet represents an operation code for indexing operations on arrays, maps, or slices within the virtual machine.
	OpIndexGet

	// OpIndexSet represents an operation code for setting a value in an array, map, or slice.
	OpIndexSet

	// OpIndexSlice is a constant representing the operation code for slice-based indexing in bytecode execution.
	OpIndexSlice

	// OpCall represents the opcode for function or method invocation with specified argument and receiver counts.
	OpCall

	// OpCallMethod represents an opcode for invoking a method directly on an object with specified arguments.
	OpCallMethod

	// OpCallImportGlobal represents an opcode for invoking a global function imported from another module.
	OpCallImportGlobal

	// OpReturn represents the opcode for returning from a function or operation, potentially with a value.
	OpReturn

	// OpDefer represents an opcode for deferring the execution of a function call until the surrounding function returns.
	OpDefer

	// OpGlobalDefine represents an opcode for defining a new global variable in the globals scope.
	OpGlobalDefine

	// OpGlobalGet retrieves a value from the globals scope by its index in the constants pool.
	OpGlobalGet

	// OpGlobalSet is an opcode used to assign a value to a globals variable within a globals scope.
	OpGlobalSet

	// OpGlobalIndex represents an operation for setting a value in a globally selected field or property.
	OpGlobalIndex

	// OpGlobalCopy is an opcode used to copy a value from one global variable to another.
	OpGlobalCopy

	// OpGlobalPtrGet retrieves a pointer to a global variable from the globals scope.
	OpGlobalPtrGet

	// OpLocalGet is an opcode used to retrieve the value of a local variable from the current scope by its index.
	OpLocalGet

	// OpLocalSet is an opcode representing the operation of setting a value to a variable in the local scope.
	OpLocalSet

	// OpLocalDefine is an opcode used to define a new local variable within the local scope of the current function.
	OpLocalDefine

	// OpLocalIndex represents the opcode for setting a value in a local variable with a selector (e.g., struct field or map key).
	OpLocalIndex

	// OpLocalPtrGet represents an opcode used to retrieve a pointer to a local variable in the current execution scope.
	OpLocalPtrGet

	// OpFreePtrGet retrieves the pointer to a variable from the free variables scope for further operations.
	OpFreePtrGet

	// OpFreeGet represents the opcode for retrieving a value from a free variable in a closure context.
	OpFreeGet

	// OpFreeSet is an opcode used to set the value of a free variable in an enclosing scope.
	OpFreeSet

	// OpCreateClosure represents the opcode used to create a function closure with constants and free variables.
	OpCreateClosure

	// OpIteratorInit initializes an iterator for iterating over a collection or data structure.
	OpIteratorInit

	// OpIteratorNext is a constant representing the operation to move the iterator to the next element in a collection.
	OpIteratorNext

	// OpIteratorKey represents the operation to retrieve the current key from an iterator.
	OpIteratorKey

	// OpIteratorValue represents the OpcodeId used to retrieve the current value during an iterator operation.
	OpIteratorValue

	// OpArithmetic represents an operation code for performing arithmetic operations between operands.
	OpArithmetic

	// OpLogical represents an opcode for performing logical operations within the instruction set.
	OpLogical

	// OpImport represents an opcode for handling imports, typically operating with two associated operands.
	OpImport

	// OpIntLogical performs integer-specific logical operations such as AND, OR, or XOR for the appropriate operands.
	OpIntLogical

	// OpIntArithmetic represents an operation code for performing integer arithmetic instructions.
	OpIntArithmetic

	// OpDerefGet is an opcode that dereferences a pointer or reference to retrieve its value.
	OpDerefGet

	// OpDerefSet represents an operation that assigns a value to the memory location pointed to by a dereferenced pointer.
	OpDerefSet

	// OpSuspend represents an opcode used to pause the execution of a process or coroutine until it is resumed.
	OpSuspend

	// OpCreateError is a constant OpcodeId representing an error operation in the instruction set.
	OpCreateError

	// OpNoOp represents a no-operation opcode, often used as a placeholder or for instruction alignment.
	OpNoOp

	// OpUnknown represents an undefined and latest opcode in the instruction set
	OpUnknown
)

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
	op.create(OpCreateClosure, []OperandFeature{SzUint8, Relocatable}, "OpCreateClosure")
	op.create(OpCallImportGlobal, []OperandFeature{SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, SzUint16, Relocatable}, "OpCallImportGlobal")
	op.create(OpGlobalDefine, []OperandFeature{Relocatable}, "OpGlobalDefine")
	op.create(OpGlobalGet, []OperandFeature{Relocatable}, "OpGlobalGet")
	op.create(OpGlobalSet, []OperandFeature{Relocatable}, "OpGlobalSet")
	op.create(OpGlobalPtrGet, []OperandFeature{Relocatable}, "OpGlobalPtrGet")
	op.create(OpGlobalCopy, []OperandFeature{Relocatable, Relocatable}, "OpGlobalCopy")
	op.create(OpGlobalIndex, []OperandFeature{SzUint8, Relocatable}, "OpGlobalIndex")
	op.create(OpPop, noOperands, "OpPop")
	op.create(OpTrue, noOperands, "OpTrue")
	op.create(OpFalse, noOperands, "OpFalse")
	op.create(OpUnaryAdd, noOperands, "OpUnaryAdd")
	op.create(OpUnarySub, noOperands, "OpUnarySub")
	op.create(OpUnaryNot, noOperands, "OpUnaryNot")
	op.create(OpUnaryBitwiseComplement, noOperands, "OpUnaryBitwiseComplement")
	op.create(OpJumpFalsy, []OperandFeature{SzUint16}, "OpJumpFalsy")
	op.create(OpJumpTruthy, []OperandFeature{SzUint16}, "OpJumpTruthy")
	op.create(OpJumpAnd, []OperandFeature{SzUint16}, "OpJumpAnd")
	op.create(OpJumpOr, []OperandFeature{SzUint16}, "OpJumpOr")
	op.create(OpJump, []OperandFeature{SzUint16}, "OpJump")
	op.create(OpJumpNotError, []OperandFeature{SzUint16}, "OpJumpNotError")
	op.create(OpJumpIndirect, noOperands, "OpJumpIndirect")
	op.create(OpNull, noOperands, "OpNull")
	op.create(OpCreateArray, []OperandFeature{SzUint16}, "OpCreateArray")
	op.create(OpCreateMap, []OperandFeature{SzUint16}, "OpCreateMap")
	op.create(OpCreateStruct, []OperandFeature{SzUint16}, "OpCreateStruct")
	op.create(OpCreateInterface, []OperandFeature{SzUint8}, "OpCreateInterface")
	op.create(OpIndexGet, noOperands, "OpIndexGet")
	op.create(OpIndexSet, noOperands, "OpIndexSet")
	op.create(OpIndexSlice, noOperands, "OpIndexSlice")
	op.create(OpCall, []OperandFeature{SzUint8, SzUint8}, "OpCall")
	op.create(OpCallMethod, []OperandFeature{SzUint16, SzUint8}, "OpCallMethod")
	op.create(OpReturn, []OperandFeature{SzUint8}, "OpReturn")
	op.create(OpDefer, noOperands, "OpDefer")
	op.create(OpLocalGet, []OperandFeature{SzUint16}, "OpLocalGet")
	op.create(OpLocalSet, []OperandFeature{SzUint16}, "OpLocalSet")
	op.create(OpLocalDefine, []OperandFeature{SzUint16}, "OpLocalDefine")
	op.create(OpLocalPtrGet, []OperandFeature{SzUint16}, "OpLocalPtrGet")
	op.create(OpLocalIndex, []OperandFeature{SzUint8, SzUint16}, "OpLocalIndex")
	op.create(OpFreeGet, []OperandFeature{SzUint16}, "OpFreeGet")
	op.create(OpFreeSet, []OperandFeature{SzUint16}, "OpFreeSet")
	op.create(OpFreePtrGet, []OperandFeature{SzUint16}, "OpFreePtrGet")
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
	op.create(OpCreateError, noOperands, "OpCreateError")
	return op
}

// Opcode retrieves the Opcode associated with the given OpcodeId and returns it along with a boolean indicating success.
func (op *Opcodes) Opcode(opcodeId OpcodeId) *Opcode {
	return op.container[int(opcodeId)&op.mask]
}

// Compile generates bytecode for a given opcode and its operands or returns an error if validation fails.
func (op *Opcodes) Compile(opcodeId OpcodeId, operands ...int) ([]byte, error) {
	opcode := op.Opcode(opcodeId)
	if opcode.OpcodeId() == OpUnknown {
		return nil, fmt.Errorf("compile: Unknown opcode: %d", opcodeId)
	}
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
