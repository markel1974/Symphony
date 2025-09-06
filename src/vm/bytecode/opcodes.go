package bytecode

import (
	"fmt"
)

const (
	// Uint8Size defines the size in bytes of a uint8 type, which is commonly used for single-byte operand encoding.
	Uint8Size = 1
	// Uint16Size represents the size in bytes of a 16-bit unsigned integer, which is 2 bytes.
	Uint16Size = 2
	// Uint32Size represents the size in bytes of a 32-bit unsigned integer.
	Uint32Size = 4
	// Uint64Size represents the size in bytes of a 64-bit unsigned integer, which is 8.
	Uint64Size = 8
)

const (
	// uint8Mask defines the maximum value that can be represented within 1 byte (8 bits), calculated as (1 << 8) - 1.
	uint8Mask = (1 << (8 * Uint8Size)) - 1
	// uint16Mask represents the maximum value that can be stored in a 16-bit unsigned integer (65535).
	uint16Mask = (1 << (8 * Uint16Size)) - 1

	uint32Mask = (1 << (8 * Uint32Size)) - 1 // Aggiunto
	// OpcodesLen defines the total number of opcodes available in the bytecode system, typically calculated as Uint8Size << 8.
	OpcodesLen = 1 << (8 * Uint8Size)
	// OpcodesMask defines the mask applied to OpcodeId values to ensure they fit within the allowable range.
	OpcodesMask = uint8Mask
)

// OpcodeId is a type alias for byte, used to represent operation codes in instruction sets.
type OpcodeId = byte

const (
	// OpConstant represents an operation that loads a constant value onto the stack.
	OpConstant OpcodeId = iota

	// OpBitwiseComplement represents the opcode for performing a bitwise complement operation on a value.
	OpBitwiseComplement

	// OpPop is an OpcodeId used to remove and discard the top value from the stack.
	OpPop

	// OpTrue represents the opcode for pushing the boolean value true onto the stack.
	OpTrue

	// OpFalse represents an operation code for pushing the boolean value 'false' onto the stack in the virtual machine.
	OpFalse

	// OpMinus represents the opcode for performing unary negation operations.
	OpMinus

	// OpNot represents the logical NOT (!) operation in the opcode set.
	OpNot

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

	OpJumpIndirect

	// OpNull represents a null operation or a placeholder indicating a null value in the opcode sequence.
	OpNull

	// OpArray represents the opcode for creating a new array with a specified number of elements from the operand.
	OpArray

	// OpMap defines an opcode representing the creation of a map structure with a specified number of key-value pairs.
	OpMap

	// OpStruct represents an opcode for initializing a struct with a specified number of key-value pairs.
	OpStruct

	// OpInterface represents the opcode used to construct an interface object with required method bindings on the stack.
	OpInterface

	// OpJumpNotError handles the 'if err != nil' pattern, skipping the if block if the top stack object is null or not an error.
	OpJumpNotError

	// OpTypeAssert implements type assertion 'val, ok := i.(Type)'.
	OpTypeAssert

	// OpIsType is a helper for type switch, checks type and pushes a boolean.
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

	// OpGlobalGet retrieves a value from the globals scope by its index in the constants pool.
	OpGlobalGet

	// OpGlobalSet is an opcode used to assign a value to a globals variable within a globals scope.
	OpGlobalSet

	// OpGlobalSelSet represents an operation for setting a value in a globally selected field or property.
	OpGlobalSelSet

	// OpGlobalCopy is an opcode used to copy a value from one global variable to another.
	OpGlobalCopy

	// OpLocalGet is an opcode used to retrieve the value of a local variable from the current scope by its index.
	OpLocalGet

	// OpLocalSet is an opcode representing the operation of setting a value to a variable in the local scope.
	OpLocalSet

	// OpLocalDefine is an opcode used to define a new local variable within the local scope of the current function.
	OpLocalDefine

	// OpLocalSelSet represents the opcode for setting a value in a local variable with a selector (e.g., struct field or map key).
	OpLocalSelSet

	// OpLocalPtrGet represents an opcode used to retrieve a pointer to a local variable in the current execution scope.
	OpLocalPtrGet

	// OpFreePtrGet retrieves the pointer to a variable from the free variables scope for further operations.
	OpFreePtrGet

	// OpFreeGet represents the opcode for retrieving a value from a free variable in a closure context.
	OpFreeGet

	// OpFreeSet is an opcode used to set the value of a free variable in an enclosing scope.
	OpFreeSet

	// OpClosure represents the opcode used to create a function closure with constants and free variables.
	OpClosure

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

	// OpFuncImport represents an opcode for handling imports, typically operating with two associated operands.
	OpFuncImport

	// OpFuncInternal represents an opcode for internal operations within the virtual machine.
	//OpFuncInternal

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

	// OpError is a constant OpcodeId representing an error operation in the instruction set.
	OpError

	// OpNoOp represents a no-operation opcode, often used as a placeholder or for instruction alignment.
	OpNoOp

	// OpUnknown represents an undefined or placeholder opcode in the instruction set, typically used as a default value.
	OpUnknown
)

// Opcode represents the details of an opcode, including its identifier, its operands, and its name.
type Opcode struct {
	opcodeId OpcodeId
	operands []int
	name     string
	offset   int
}

// NewOpcode creates a new Opcode instance, initializing its opcode, operands, and name fields.
func NewOpcode(opcodeId OpcodeId, operands []int, name string) *Opcode {
	od := &Opcode{
		opcodeId: opcodeId,
		operands: operands,
		name:     name,
		offset:   0,
	}
	for _, w := range od.operands {
		od.offset += w
	}
	return od
}

// OpcodeId returns the opcode associated with the Opcode instance.
func (od *Opcode) OpcodeId() OpcodeId {
	return od.opcodeId
}

// Name returns the name of the opcode as a string.
func (od *Opcode) Name() string {
	return od.name
}

// Operands returns the operand widths for the Opcode as a slice of integers.
func (od *Opcode) Operands() []int {
	return od.operands
}

// Offset returns the byte offset of the opcode within an instruction.
func (od *Opcode) Offset() int {
	return od.offset
}

// Opcodes is a collection that manages and organizes Opcode instances, providing methods to create, retrieve, or compile them.
type Opcodes struct {
	container []*Opcode
	maxLen    int
}

// NewOpcodes initializes and returns a new Opcodes instance with predefined opcode mappings.
func NewOpcodes() *Opcodes {
	op := &Opcodes{
		container: make([]*Opcode, OpcodesLen),
	}
	for i := range op.container {
		op.container[i] = NewOpcode(OpUnknown, []int{}, "OpUnknown")
	}
	op.createOpcode(OpConstant, []int{Uint16Size}, "OpConstant")
	op.createOpcode(OpPop, []int{}, "OpPop")
	op.createOpcode(OpTrue, []int{}, "OpTrue")
	op.createOpcode(OpFalse, []int{}, "OpFalse")
	op.createOpcode(OpBitwiseComplement, []int{}, "OpBitwiseComplement")
	op.createOpcode(OpMinus, []int{}, "OpMinus")
	op.createOpcode(OpNot, []int{}, "OpNot")
	op.createOpcode(OpJumpFalsy, []int{Uint16Size}, "OpJumpFalsy")
	op.createOpcode(OpJumpTruthy, []int{Uint16Size}, "OpJumpTruthy")
	op.createOpcode(OpJumpAnd, []int{Uint16Size}, "OpJumpAnd")
	op.createOpcode(OpJumpOr, []int{Uint16Size}, "OpJumpOr")
	op.createOpcode(OpJump, []int{Uint16Size}, "OpJump")
	op.createOpcode(OpJumpNotError, []int{Uint16Size}, "OpJumpNotError")
	op.createOpcode(OpJumpIndirect, []int{}, "OpJumpIndirect")
	op.createOpcode(OpNull, []int{}, "OpNull")
	op.createOpcode(OpGlobalGet, []int{Uint16Size}, "OpGlobalGet")
	op.createOpcode(OpGlobalSet, []int{Uint16Size}, "OpGlobalSet")
	op.createOpcode(OpGlobalSelSet, []int{Uint16Size, Uint8Size}, "OpGlobalSelSet")
	op.createOpcode(OpGlobalCopy, []int{Uint16Size, Uint16Size}, "OpGlobalCopy")
	op.createOpcode(OpArray, []int{Uint16Size}, "OpArray")
	op.createOpcode(OpMap, []int{Uint16Size}, "OpMap")
	op.createOpcode(OpStruct, []int{Uint16Size}, "OpStruct")
	op.createOpcode(OpInterface, []int{Uint8Size}, "OpInterface")
	op.createOpcode(OpIndexGet, []int{}, "OpIndexGet")
	op.createOpcode(OpIndexSet, []int{}, "OpIndexSet")
	op.createOpcode(OpIndexSlice, []int{}, "OpIndexSlice")
	op.createOpcode(OpCall, []int{Uint8Size, Uint8Size}, "OpCall")
	op.createOpcode(OpCallMethod, []int{Uint16Size, Uint8Size}, "OpCallMethod")
	op.createOpcode(OpCallImportGlobal, []int{Uint16Size, Uint16Size, Uint16Size, Uint16Size, Uint16Size, Uint16Size, Uint16Size, Uint16Size}, "OpCallImportGlobal")
	op.createOpcode(OpReturn, []int{Uint8Size}, "OpReturn")
	op.createOpcode(OpLocalGet, []int{Uint8Size}, "OpLocalGet")
	op.createOpcode(OpLocalSet, []int{Uint8Size}, "OpLocalSet")
	op.createOpcode(OpLocalDefine, []int{Uint8Size}, "OpLocalDefine")
	op.createOpcode(OpLocalPtrGet, []int{Uint8Size}, "OpLocalPtrGet")
	op.createOpcode(OpLocalSelSet, []int{Uint8Size, Uint8Size}, "OpLocalSelSet")
	op.createOpcode(OpClosure, []int{Uint16Size, Uint8Size}, "OpClosure")
	op.createOpcode(OpFreeGet, []int{Uint8Size}, "OpFreeGet")
	op.createOpcode(OpFreeSet, []int{Uint8Size}, "OpFreeSet")
	op.createOpcode(OpFreePtrGet, []int{Uint8Size}, "OpFreePtrGet")
	op.createOpcode(OpIteratorInit, []int{Uint8Size}, "OpIteratorInit")
	op.createOpcode(OpIteratorNext, []int{Uint8Size}, "OpIteratorNext")
	op.createOpcode(OpIteratorKey, []int{Uint8Size}, "OpIteratorKey")
	op.createOpcode(OpIteratorValue, []int{Uint8Size}, "OpIteratorValue")
	op.createOpcode(OpLogical, []int{Uint8Size}, "OpLogical")
	op.createOpcode(OpArithmetic, []int{Uint8Size}, "OpArithmetic")
	op.createOpcode(OpFuncImport, []int{Uint16Size}, "OpFuncImport")
	//op.createOpcode(OpFuncInternal, []int{Uint16Size}, "OpFuncInternal")
	op.createOpcode(OpIntLogical, []int{Uint16Size, Uint16Size, Uint16Size, Uint8Size}, "OpIntLogical")
	op.createOpcode(OpIntArithmetic, []int{Uint16Size, Uint16Size, Uint16Size, Uint8Size}, "OpIntArithmetic")
	op.createOpcode(OpDerefGet, []int{}, "OpDerefGet")
	op.createOpcode(OpDerefSet, []int{}, "OpDerefSet")
	op.createOpcode(OpTypeAssert, []int{Uint16Size}, "OpTypeAssert")
	op.createOpcode(OpIsType, []int{Uint16Size}, "OpIsType")
	op.createOpcode(OpAsType, []int{Uint16Size}, "OpAsType")
	op.createOpcode(OpSuspend, []int{}, "OpSuspend")
	op.createOpcode(OpError, []int{}, "OpError")
	return op
}

// createOpcode registers a new Opcode in the Opcodes container with its identifier, operands, and name.
func (op *Opcodes) createOpcode(opcodeId OpcodeId, operands []int, name string) {
	od := NewOpcode(opcodeId, operands, name)
	op.container[od.opcodeId&OpcodesMask] = od
	if len(operands) > op.maxLen {
		op.maxLen = len(operands)
	}
}

// Opcode retrieves the *Opcode instance corresponding to the given OpcodeId from the container, applying a mask.
func (op *Opcodes) Opcode(opcodeId OpcodeId) *Opcode {
	return op.container[opcodeId&OpcodesMask]
}

// CompileInstruction generates bytecode for a given opcode and its operands or returns an error if validation fails.
func (op *Opcodes) CompileInstruction(opcode OpcodeId, operands ...int) ([]byte, error) {
	details := op.Opcode(opcode)
	numOperands := details.Operands()
	if len(operands) != len(numOperands) {
		return nil, fmt.Errorf("wrong number of operands for %s: want %d, got %d", details.Name(), len(numOperands), len(operands))
	}
	totalLen := 1
	totalLen += details.Offset()
	instruction := make([]byte, totalLen)
	instruction[0] = opcode
	offset := 1
	for i, o := range operands {
		width := numOperands[i]
		switch width {
		case Uint8Size:
			if o < 0 || o > uint8Mask {
				return nil, fmt.Errorf("operand %d value %d out of 1-byte range", i, o)
			}
			instruction[offset] = byte(o)
		case Uint16Size:
			if o < 0 || o > uint16Mask {
				return nil, fmt.Errorf("operand %d value %d out of 2-byte range", i, o)
			}
			n := uint16(o)
			instruction[offset] = byte(n >> 8) // Most significant byte (Big Endian)
			instruction[offset+1] = byte(n)    // Least significant byte
		case Uint32Size:
			if o < 0 || o > uint32Mask {
				return nil, fmt.Errorf("operand %d value %d out of 4-byte range", i, o)
			}
			n := uint32(o)
			instruction[offset] = byte(n >> 24)
			instruction[offset+1] = byte(n >> 16)
			instruction[offset+2] = byte(n >> 8)
			instruction[offset+3] = byte(n)
		case Uint64Size:
			if o < 0 {
				return nil, fmt.Errorf("operand %d value %d out of 8-byte range (negative)", i, o)
			}
			n := uint64(o)
			instruction[offset] = byte(n >> 56)
			instruction[offset+1] = byte(n >> 48)
			instruction[offset+2] = byte(n >> 40)
			instruction[offset+3] = byte(n >> 32)
			instruction[offset+4] = byte(n >> 24)
			instruction[offset+5] = byte(n >> 16)
			instruction[offset+6] = byte(n >> 8)
			instruction[offset+7] = byte(n)
		}
		offset += width
	}
	return instruction, nil
}

// MaxLen returns the maximum number of operands encountered across all compiled instructions in the Opcodes instance.
func (op *Opcodes) MaxLen() int {
	return op.maxLen
}
