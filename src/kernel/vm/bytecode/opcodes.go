package bytecode

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// OpcodesLen defines the length of Opcodes as 256, calculated using a bitwise shift operation.
// OpcodesMask provides a bitmask for Opcodes of length 256 by subtracting 1 from OpcodesLen.
const (
	ByteSize   = 1
	Uint16Size = 2
	byteMask   = (ByteSize << 8) - 1
	uint16Mask = (1 << 16) - 1

	OpcodesLen  = ByteSize << 8
	OpcodesMask = byteMask
)

// Opcode is a type alias for byte, used to represent operation codes in instruction sets.
type Opcode = byte

// OpConstant loads a constant onto the stack.
// OpBitwiseComplement performs a bitwise complement operation.
// OpPop pops a value from the stack.
// OpTrue pushes the boolean value true onto the stack.
// OpFalse pushes the boolean value false onto the stack.
// OpEqual checks if two values are equal (==).
// OpNotEqual checks if two values are not equal (!=).
// OpMinus performs a subtraction or negation (-).
// OpNotLogical performs a logical NOT operation (!).
// OpJumpFalsy jumps if the top of the stack is falsy.
// OpJumpAnd performs a logical AND and jumps.
// OpJumpOr performs a logical OR and jumps.
// OpJump performs an unconditional jump.
// OpNull pushes a null value onto the stack.
// OpArray creates an array object.
// OpMap creates a map object.
// OpImmutable creates an immutable object.
// OpIndex performs an index operation.
// OpIndexSlice performs a slice operation.
// OpCall calls a function.
// OpReturn returns from a function call.
// OpGlobalGet retrieves a global variable.
// OpGlobalSet sets the value of a global variable.
// OpGlobalSelSet sets a global variable using selectors.
// OpLocalGet retrieves a local variable.
// OpLocalSet sets the value of a local variable.
// OpLocalDefine defines a new local variable.
// OpLocalSelSet sets a local variable using selectors.
// OpFreePtrGet retrieves a free variable pointer object.
// OpFreeGet retrieves a free variable.
// OpFreeSet sets the value of a free variable.
// OpLocalPtrGet retrieves a local variable as a pointer.
// OpFreeSelSet sets a free variable using selectors.
// OpGetBuiltin retrieves a builtin function.
// OpClosure creates a closure and pushes it onto the stack.
// OpIteratorInit initializes an iterator.
// OpIteratorNext advances an iterator to the next element.
// OpIteratorKey retrieves the key from the current iterator position.
// OpIteratorValue retrieves the value from the current iterator position.
// OpBinary performs a binary operation.
// OpReferences refers to a specific operation related to references.
// OpSuspend suspends the virtual machine.
// OpError creates an error object.
// OpUnknown represents an unknown operation.
const (
	OpConstant          Opcode = iota // Load constant
	OpBitwiseComplement               // bitwise complement
	OpPop                             // Pop
	OpTrue                            // Push true
	OpFalse                           // Push false
	OpEqual                           // Equal ==
	OpNotEqual                        // Not equal !=
	OpMinus                           // Minus -
	OpNotLogical                      // Logical not
	OpJumpFalsy                       // Jump if falsy
	OpJumpAnd                         // Logical AND jump
	OpJumpOr                          // Logical OR jump
	OpJump                            // Jump
	OpNull                            // Push null
	OpArray                           // Array object
	OpMap                             // Map object
	OpStruct                          // Struct object
	OpImmutable                       // Immutable object
	OpIndex                           // Index operation
	OpIndexSlice                      // Slice operation
	OpCall                            // Call function
	OpReturn                          // Return
	OpGlobalGet                       // Get global variable
	OpGlobalSet                       // Set global variable
	OpGlobalSelSet                    // Set global variable using selectors
	OpLocalGet                        // Get local variable
	OpLocalSet                        // Set local variable
	OpLocalDefine                     // Define local variable
	OpLocalSelSet                     // Set local variable using selectors
	OpLocalPtrGet                     // Get local variable as a pointer
	OpFreePtrGet                      // Get free variable pointer object
	OpFreeGet                         // Get free variables
	OpFreeSet                         // Set free variables
	OpFreeSelSet                      // Set free variables using selectors
	OpClosure                         // Push closure
	OpIteratorInit                    // Iterator init
	OpIteratorNext                    // Iterator next
	OpIteratorKey                     // Iterator key
	OpIteratorValue                   // Iterator value
	OpBinary                          // Binary operation
	OpReferences
	OpIntOp
	OpDeref   // Dereference a pointer
	OpSuspend // Suspend VM
	OpError   // Error object
	OpNoOp    // Push null
	OpUnknown
)

// OpcodeDetails represents the details of an opcode, including its identifier, its operands, and its name.
type OpcodeDetails struct {
	factory  objects.IGateKeeper
	opcode   Opcode
	operands []int
	name     string
	offset   int
}

// NewOpcodeDetails creates a new OpcodeDetails instance, initializing its opcode, operands, and name fields.
func NewOpcodeDetails(factory objects.IGateKeeper, opcode Opcode, operands []int, name string) *OpcodeDetails {
	od := &OpcodeDetails{
		factory:  factory,
		opcode:   opcode,
		operands: operands,
		name:     name,
		offset:   0,
	}
	for _, w := range od.operands {
		od.offset += w
	}
	return od
}

// Opcode returns the opcode associated with the OpcodeDetails instance.
func (od *OpcodeDetails) Opcode() Opcode {
	return od.opcode
}

func (od *OpcodeDetails) Factory() objects.IGateKeeper {
	return od.factory
}

// Name returns the name of the opcode as a string.
func (od *OpcodeDetails) Name() string {
	return od.name
}

// Operands retrieve the list of integer operands associated with the OpcodeDetails instance.
func (od *OpcodeDetails) Operands() []int {
	return od.operands
}
func (od *OpcodeDetails) Offset() int {
	return od.offset
}

type Opcodes struct {
	factory objects.IGateKeeper
	details []*OpcodeDetails
}

func NewOpcodes(factory objects.IGateKeeper) *Opcodes {
	op := &Opcodes{
		factory: factory,
		details: make([]*OpcodeDetails, OpcodesLen),
	}
	for i := range op.details {
		op.details[i] = NewOpcodeDetails(factory, OpUnknown, []int{}, "OpUnknown")
	}
	op.createOpcodeDetails(factory, OpConstant, []int{2}, "OpConstant")
	op.createOpcodeDetails(factory, OpPop, []int{}, "OpPop")
	op.createOpcodeDetails(factory, OpTrue, []int{}, "OpTrue")
	op.createOpcodeDetails(factory, OpFalse, []int{}, "OpFalse")
	op.createOpcodeDetails(factory, OpBitwiseComplement, []int{}, "OpBitwiseComplement")
	op.createOpcodeDetails(factory, OpEqual, []int{}, "OpEqual")
	op.createOpcodeDetails(factory, OpNotEqual, []int{}, "OpNotEqual")
	op.createOpcodeDetails(factory, OpMinus, []int{}, "OpMinus")
	op.createOpcodeDetails(factory, OpNotLogical, []int{}, "OpNotLogical")
	op.createOpcodeDetails(factory, OpJumpFalsy, []int{2}, "OpJumpFalsy")
	op.createOpcodeDetails(factory, OpJumpAnd, []int{2}, "OpJumpAnd")
	op.createOpcodeDetails(factory, OpJumpOr, []int{2}, "OpJumpOr")
	op.createOpcodeDetails(factory, OpJump, []int{2}, "OpJump")
	op.createOpcodeDetails(factory, OpNull, []int{}, "OpNull")
	op.createOpcodeDetails(factory, OpGlobalGet, []int{2}, "OpGlobalGet")
	op.createOpcodeDetails(factory, OpGlobalSet, []int{2}, "OpGlobalSet")
	op.createOpcodeDetails(factory, OpGlobalSelSet, []int{2, 1}, "OpGlobalSelSet")
	op.createOpcodeDetails(factory, OpArray, []int{2}, "OpArray")
	op.createOpcodeDetails(factory, OpMap, []int{2}, "OpMap")
	op.createOpcodeDetails(factory, OpStruct, []int{2}, "OpStruct")
	op.createOpcodeDetails(factory, OpImmutable, []int{}, "OpImmutable")
	op.createOpcodeDetails(factory, OpIndex, []int{}, "OpIndex")
	op.createOpcodeDetails(factory, OpIndexSlice, []int{}, "OpIndexSlice")
	op.createOpcodeDetails(factory, OpCall, []int{1, 1}, "OpCall")
	op.createOpcodeDetails(factory, OpReturn, []int{1}, "OpReturn")
	op.createOpcodeDetails(factory, OpLocalGet, []int{1}, "OpLocalGet")
	op.createOpcodeDetails(factory, OpLocalSet, []int{1}, "OpLocalSet")
	op.createOpcodeDetails(factory, OpLocalDefine, []int{1}, "OpLocalDefine")
	op.createOpcodeDetails(factory, OpLocalSelSet, []int{1, 1}, "OpLocalSelSet")
	op.createOpcodeDetails(factory, OpClosure, []int{2, 1}, "OpClosure")
	op.createOpcodeDetails(factory, OpFreePtrGet, []int{1}, "OpFreePtrGet")
	op.createOpcodeDetails(factory, OpFreeGet, []int{1}, "OpFreeGet")
	op.createOpcodeDetails(factory, OpFreeSet, []int{1}, "OpFreeSet")
	op.createOpcodeDetails(factory, OpLocalPtrGet, []int{1}, "OpLocalPtrGet")
	op.createOpcodeDetails(factory, OpFreeSelSet, []int{1, 1}, "OpFreeSelSet")
	op.createOpcodeDetails(factory, OpIteratorInit, []int{1}, "OpIteratorInit")
	op.createOpcodeDetails(factory, OpIteratorNext, []int{1}, "OpIteratorNext")
	op.createOpcodeDetails(factory, OpIteratorKey, []int{1}, "OpIteratorKey")
	op.createOpcodeDetails(factory, OpIteratorValue, []int{1}, "OpIteratorValue")
	op.createOpcodeDetails(factory, OpBinary, []int{1}, "OpBinary")
	op.createOpcodeDetails(factory, OpReferences, []int{2}, "OpReferences")
	op.createOpcodeDetails(factory, OpIntOp, []int{2, 1}, "OpIntOp")
	op.createOpcodeDetails(factory, OpDeref, []int{}, "OpDeref")
	op.createOpcodeDetails(factory, OpSuspend, []int{}, "OpSuspend")
	op.createOpcodeDetails(factory, OpError, []int{}, "OpError")
	return op
}

// Factory returns the factory associated with the Opcodes instance.
func (op *Opcodes) Factory() objects.IGateKeeper {
	return op.factory
}

// createOpcodeDetails associates opcode details with a specific opcode, storing it in a global lookup by applying a mask.
func (op *Opcodes) createOpcodeDetails(factory objects.IGateKeeper, opcode Opcode, operands []int, name string) {
	od := NewOpcodeDetails(factory, opcode, operands, name)
	op.details[od.opcode&OpcodesMask] = od
}

// OpcodeToDetails retrieves the OpcodeDetails corresponding to the given opcode by applying the OpcodesMask.
func (op *Opcodes) OpcodeToDetails(opcode Opcode) *OpcodeDetails {
	return op.details[opcode&OpcodesMask]
}

// OpcodeToOperands retrieves the operand widths for a given opcode from its details.
func (op *Opcodes) OpcodeToOperands(opcode Opcode) []int {
	details := op.OpcodeToDetails(opcode)
	return details.Operands()
}

// OpcodeToOperandsOffset calculates the total byte offset for the operands of a given opcode.
func (op *Opcodes) OpcodeToOperandsOffset(opcode Opcode) int {
	details := op.OpcodeToDetails(opcode)
	return details.Offset()
}

// OpcodeToOperandsDetails extracts operand details from a given opcode and instruction sequence, returning operand widths, values, and bytes read.
func (op *Opcodes) OpcodeToOperandsDetails(opcode Opcode, ins []byte) ([]int, []int, int) {
	details := op.OpcodeToDetails(opcode)
	if len(details.Operands()) == 0 {
		return nil, nil, 0
	}
	var retOperands []int
	var offset int
	for _, width := range details.Operands() {
		switch width {
		case ByteSize:
			if offset >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset]))
		case Uint16Size:
			if offset+1 >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset+1])|int(ins[offset])<<8)
		}
		offset += width
	}
	return details.Operands(), retOperands, offset
}

// OpcodeName returns the name of the provided opcode as a string.
func (op *Opcodes) OpcodeName(opcode Opcode) string {
	details := op.OpcodeToDetails(opcode)
	return details.Name()
}

// CompileInstruction generates a byte-encoded instruction using the given opcode and operands, verifying operand widths.
func (op *Opcodes) CompileInstruction(opcode Opcode, operands ...int) ([]byte, error) {
	details := op.OpcodeToDetails(opcode)
	numOperands := details.Operands()
	if len(operands) != len(numOperands) {
		return nil, fmt.Errorf(
			"wrong number of operands for %s: want %d, got %d", op.OpcodeName(opcode), len(numOperands), len(operands))
	}
	totalLen := 1
	totalLen += details.Offset()
	instruction := make([]byte, totalLen)
	instruction[0] = opcode
	offset := 1
	for i, o := range operands {
		width := numOperands[i]
		switch width {
		case ByteSize:
			if o < 0 || o > byteMask {
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
		}
		offset += width
	}
	return instruction, nil
}
