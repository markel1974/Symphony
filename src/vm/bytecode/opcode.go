package bytecode

const OpcodeWidth = 1

// TODO use OperandFeature
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

type OperandFeature int

const (
	SizeMask OperandFeature = 0x0F
	Size1    OperandFeature = 1
	Size2    OperandFeature = 2
	Size4    OperandFeature = 4
	Size8    OperandFeature = 8

	IsRelocatable  OperandFeature = 1 << 4
	IsSigned       OperandFeature = 1 << 5
	IsPointer      OperandFeature = 1 << 6
	HintForGC      OperandFeature = 1 << 7
	HintForJIT     OperandFeature = 1 << 8
	IsRegisterHint OperandFeature = 1 << 9
)

const (
	ImmediateUint8  = Size1
	ImmediateInt16  = Size2 | IsSigned
	ConstantIndex16 = Size2 | IsRelocatable
	PointerToLocal  = Size1 | IsPointer | HintForGC
)

const (
	// uint8Mask defines the maximum value that can be represented within 1 byte (8 bits), calculated as (1 << 8) - 1.
	uint8Mask = (1 << (8 * Uint8Size)) - 1
	// uint16Mask represents the maximum value that can be stored in a 16-bit unsigned integer (65535).
	uint16Mask = (1 << (8 * Uint16Size)) - 1
	// uint32Mask represents the maximum value that can be stored in a 32-bit unsigned integer (4294967295).
	uint32Mask = (1 << (8 * Uint32Size)) - 1 // Aggiunto
)

type Relocatable int

const (
	OpRelocatableNone Relocatable = iota
	OpRelocatable                 // first operand relocatableId(16bit)
	//OpRelocatableFree             // first operand relocatableId(16bit) second operand numFree(8bit)
)

const (
// OpcodesLen defines the total number of opcodes available in the bytecode system, typically calculated as Uint8Size << 8.
// OpcodesLen = 1 << (8)
// OpcodesMask defines the mask applied to OpcodeId values to ensure they fit within the allowable range.
// OpcodesMask = OpcodesLen - 1
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

	// OpImport represents an opcode for handling imports, typically operating with two associated operands.
	OpImport

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

	// OpUnknown represents an undefined and latest opcode in the instruction set
	OpUnknown
)

// Opcode represents the opcode of an opcode, including its identifier, its operand, and its name.
type Opcode struct {
	opcodeId      OpcodeId
	relocatable   Relocatable
	operands      []int
	name          string
	operandsWidth int
	compiler      *Compiler
}

// NewOpcode creates a new Opcode instance, initializing its opcode, operands, and name fields.
func NewOpcode(opcodeId OpcodeId, operands []int, name string, relocatable Relocatable) *Opcode {
	od := &Opcode{
		opcodeId:      opcodeId,
		operands:      operands,
		relocatable:   relocatable,
		name:          name,
		operandsWidth: 0,
	}
	for _, w := range od.operands {
		od.operandsWidth += w
	}
	od.compiler = NewCompiler(od)
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

// OperandsWidth returns the total width of the operands for the Opcode instance.
func (od *Opcode) OperandsWidth() int {
	return od.operandsWidth
}

// FullWidth returns the total width of the Opcode instance, including the opcode and operands.
func (od *Opcode) FullWidth() int {
	return OpcodeWidth + od.OperandsWidth()
}

// Relocatable returns the relocatable value associated with the Opcode instance.
func (od *Opcode) Relocatable() Relocatable {
	return od.relocatable
}

// Compile compiles the opcode into a sequence of bytes.
func (od *Opcode) Compile(operands []int) ([]byte, error) {
	if err := od.compiler.Compile(operands); err != nil {
		return nil, err
	}
	return od.compiler.Instructions(), nil
}

// Decompile decompiles the opcode into a sequence of integers.
func (od *Opcode) Decompile(bytecode []byte) ([]int, error) {
	od.compiler.SetInstructions(bytecode)
	return od.compiler.Decompile()
}
