package native

import "github.com/markel1974/c64emu/src/vm/opcodes"

const (
	// OpConstantId represents an operation that loads a constant value onto the stack.
	OpConstantId opcodes.OpcodeId = iota

	// OpPopId is an OpcodeId used to remove and discard the top value from the stack.
	OpPopId

	// OpTrueId represents the opcode for pushing the boolean value true onto the stack.
	OpTrueId

	// OpFalseId represents an operation code for pushing the boolean value 'false' onto the stack in the virtual machine.
	OpFalseId

	// OpUnarySubId represents the opcode for performing unary negation operations.
	OpUnarySubId

	// OpUnaryAddId represents an operation for applying a unary addition to a numeric value.
	OpUnaryAddId

	// OpUnaryNotId represents the logical NOT (!) operation in the opcode set.
	OpUnaryNotId

	// OpUnaryBitwiseComplementId represents the opcode for performing a bitwise complement operation on a value.
	OpUnaryBitwiseComplementId

	// OpJumpFalsyId represents a conditional jump instruction that redirects execution if the top stack value is falsy.
	OpJumpFalsyId

	// OpJumpTruthyId represents a conditional jump instruction that redirects execution if the top stack value is truthy.
	OpJumpTruthyId

	// OpJumpAndId is an opcode used to perform a conditional jump based on the evaluation of a logical AND operation.
	OpJumpAndId

	// OpJumpOrId represents an operation code used to perform a conditional jump if a logical OR condition evaluates to true.
	OpJumpOrId

	// OpJumpId is a constant representing an unconditional jump operation in the bytecode instruction set.
	OpJumpId

	// OpJumpIndirectId is a constant representing an opcode for performing an indirect jump. It does not require operands.
	OpJumpIndirectId

	// OpNullId represents a null operation or a placeholder indicating a null value in the opcode sequence.
	OpNullId

	// OpCreateArrayId represents the opcode for creating a new array with a specified number of elements from the operand.
	OpCreateArrayId

	// OpCreateMapId defines an opcode representing the creation of a map structure with a specified number of key-value pairs.
	OpCreateMapId

	// OpCreateStructId represents an opcode for initializing a struct with a specified number of key-value pairs.
	OpCreateStructId

	// OpCreateInterfaceId represents the opcode used to construct an interface object with required method bindings on the stack.
	OpCreateInterfaceId

	// OpJumpNotErrorId handles the 'if err != nil' pattern, skipping the if block if the top stack object is null or not an error.
	OpJumpNotErrorId

	// OpTypeAssertId implements type assertion 'val, ok := i.(Type)'.
	OpTypeAssertId

	// OpTypeCheckId is a helper for type switch, checks type, and pushes a boolean.
	OpTypeCheckId

	// OpAsTypeId is a helper for type switch, performs type casting without checks.
	OpAsTypeId

	// OpIndexGetId represents an operation code for indexing operations on arrays, maps, or slices within the virtual machine.
	OpIndexGetId

	// OpIndexSetId represents an operation code for setting a value in an array, map, or slice.
	OpIndexSetId

	// OpIndexSliceId is a constant representing the operation code for slice-based indexing in bytecode execution.
	OpIndexSliceId

	// OpCallId represents the opcode for function or method invocation with specified argument and receiver counts.
	OpCallId

	// OpCallAsyncId represents the opcode for invoking asynchronous function calls in the instruction set.
	OpCallAsyncId

	// OpCallInterfaceId represents an opcode for invoking a method directly on an object with specified arguments.
	OpCallInterfaceId

	// OpCallImportGlobalId represents an opcode for invoking a global function imported from another module.
	OpCallImportGlobalId

	// OpReturnId represents the opcode for returning from a function or operation, potentially with a value.
	OpReturnId

	// OpDeferId represents an opcode for deferring the execution of a function call until the surrounding function returns.
	OpDeferId

	// OpGlobalDefineId represents an opcode for defining a new global variable in the globals scope.
	OpGlobalDefineId

	// OpGlobalGetId retrieves a value from the globals scope by its index in the constants pool.
	OpGlobalGetId

	// OpGlobalSetId is an opcode used to assign a value to a globals variable within a globals scope.
	OpGlobalSetId

	// OpGlobalIndexId represents an operation for setting a value in a globally selected field or property.
	OpGlobalIndexId

	// OpGlobalCopyId is an opcode used to copy a value from one global variable to another.
	OpGlobalCopyId

	// OpGlobalPtrGetId retrieves a pointer to a global variable from the globals scope.
	OpGlobalPtrGetId

	// OpLocalGetId is an opcode used to retrieve the value of a local variable from the current scope by its index.
	OpLocalGetId

	// OpLocalSetId is an opcode representing the operation of setting a value to a variable in the local scope.
	OpLocalSetId

	// OpLocalDefineId is an opcode used to define a new local variable within the local scope of the current function.
	OpLocalDefineId

	// OpLocalIndexId represents the opcode for setting a value in a local variable with a selector (e.g., struct field or map key).
	OpLocalIndexId

	// OpLocalPtrGetId represents an opcode used to retrieve a pointer to a local variable in the current execution scope.
	OpLocalPtrGetId

	// OpFreePtrGetId retrieves the pointer to a variable from the free variables scope for further operations.
	OpFreePtrGetId

	// OpFreeGetId represents the opcode for retrieving a value from a free variable in a closure context.
	OpFreeGetId

	// OpFreeSetId is an opcode used to set the value of a free variable in an enclosing scope.
	OpFreeSetId

	// OpCreateClosureId represents the opcode used to create a function closure with constants and free variables.
	OpCreateClosureId

	// OpIteratorInitId initializes an iterator for iterating over a collection or data structure.
	OpIteratorInitId

	// OpIteratorNextId is a constant representing the operation to move the iterator to the next element in a collection.
	OpIteratorNextId

	// OpIteratorKeyId represents the operation to retrieve the current key from an iterator.
	OpIteratorKeyId

	// OpIteratorValueId represents the OpcodeId used to retrieve the current value during an iterator operation.
	OpIteratorValueId

	// OpArithmeticId represents an operation code for performing arithmetic operations between operands.
	OpArithmeticId

	// OpLogicalId represents an opcode for performing logical operations within the instruction set.
	OpLogicalId

	// OpImportId represents an opcode for handling imports, typically operating with two associated operands.
	OpImportId

	// OpIntLogicalId performs integer-specific logical operations such as AND, OR, or XOR for the appropriate operands.
	OpIntLogicalId

	// OpIntArithmeticId represents an operation code for performing integer arithmetic instructions.
	OpIntArithmeticId

	// OpDerefGetId is an opcode that dereferences a pointer or reference to retrieve its value.
	OpDerefGetId

	// OpDerefSetId represents an operation that assigns a value to the memory location pointed to by a dereferenced pointer.
	OpDerefSetId

	// OpCreateErrorId is a constant OpcodeId representing an error operation in the instruction set.
	OpCreateErrorId

	// OpNoOpId represents a no-operation opcode, often used as a placeholder or for instruction alignment.
	OpNoOpId

	// OpUnknownId represents an undefined and latest opcode in the instruction set
	OpUnknownId
)
