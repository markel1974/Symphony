package bytecode

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// OpcodesLen defines the length of Opcodes as 256, calculated using a bitwise shift operation.
// OpcodesMask provides a bitmask for Opcodes of length 256 by subtracting 1 from OpcodesLen.
const (
	OpcodesLen  = 1 << 8
	OpcodesMask = OpcodesLen - 1
)

// Opcode is a type alias for byte, used to represent operation codes in instruction sets.
type Opcode = byte

// OpConstant loads a constant onto the stack.
// OpBComplement performs a bitwise complement operation.
// OpPop pops a value from the stack.
// OpTrue pushes the boolean value true onto the stack.
// OpFalse pushes the boolean value false onto the stack.
// OpEqual checks if two values are equal (==).
// OpNotEqual checks if two values are not equal (!=).
// OpMinus performs a subtraction or negation (-).
// OpLNot performs a logical NOT operation (!).
// OpJumpFalsy jumps if the top of the stack is falsy.
// OpAndJump performs a logical AND and jumps.
// OpOrJump performs a logical OR and jumps.
// OpJump performs an unconditional jump.
// OpNull pushes a null value onto the stack.
// OpArray creates an array object.
// OpMap creates a map object.
// OpImmutable creates an immutable object.
// OpIndex performs an index operation.
// OpSliceIndex performs a slice operation.
// OpCall calls a function.
// OpReturn returns from a function call.
// OpGetGlobal retrieves a global variable.
// OpSetGlobal sets the value of a global variable.
// OpSetSelGlobal sets a global variable using selectors.
// OpGetLocal retrieves a local variable.
// OpSetLocal sets the value of a local variable.
// OpDefineLocal defines a new local variable.
// OpSetSelLocal sets a local variable using selectors.
// OpGetFreePtr retrieves a free variable pointer object.
// OpGetFree retrieves a free variable.
// OpSetFree sets the value of a free variable.
// OpGetLocalPtr retrieves a local variable as a pointer.
// OpSetSelFree sets a free variable using selectors.
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
	OpConstant      Opcode = iota // Load constant
	OpBComplement                 // bitwise complement
	OpPop                         // Pop
	OpTrue                        // Push true
	OpFalse                       // Push false
	OpEqual                       // Equal ==
	OpNotEqual                    // Not equal !=
	OpMinus                       // Minus -
	OpLNot                        // Logical not
	OpJumpFalsy                   // Jump if falsy
	OpAndJump                     // Logical AND jump
	OpOrJump                      // Logical OR jump
	OpJump                        // Jump
	OpNull                        // Push null
	OpArray                       // Array object
	OpMap                         // Map object
	OpStruct                      // Struct object
	OpImmutable                   // Immutable object
	OpIndex                       // Index operation
	OpSliceIndex                  // Slice operation
	OpCall                        // Call function
	OpReturn                      // Return
	OpGetGlobal                   // Get global variable
	OpSetGlobal                   // Set global variable
	OpSetSelGlobal                // Set global variable using selectors
	OpGetLocal                    // Get local variable
	OpSetLocal                    // Set local variable
	OpDefineLocal                 // Define local variable
	OpSetSelLocal                 // Set local variable using selectors
	OpGetFreePtr                  // Get free variable pointer object
	OpGetFree                     // Get free variables
	OpSetFree                     // Set free variables
	OpGetLocalPtr                 // Get local variable as a pointer
	OpSetSelFree                  // Set free variables using selectors
	OpClosure                     // Push closure
	OpIteratorInit                // Iterator init
	OpIteratorNext                // Iterator next
	OpIteratorKey                 // Iterator key
	OpIteratorValue               // Iterator value
	OpBinary                      // Binary operation
	OpReferences
	OpIntOp
	OpDeref   // Dereference a pointer
	OpSuspend // Suspend VM
	OpError   // Error object
	OpUnknown
)

// OpcodeDetails represents the details of an opcode, including its identifier, its operands, and its name.
type OpcodeDetails struct {
	factory  objects.IGateKeeper
	opcode   Opcode
	operands []int
	name     string
}

// NewOpcodeDetails creates a new OpcodeDetails instance, initializing its opcode, operands, and name fields.
func NewOpcodeDetails(factory objects.IGateKeeper, opcode Opcode, operands []int, name string) *OpcodeDetails {
	od := &OpcodeDetails{
		factory:  factory,
		opcode:   opcode,
		operands: operands,
		name:     name,
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
	op.createOpcodeDetails(factory, OpBComplement, []int{}, "OpBComplement")
	op.createOpcodeDetails(factory, OpEqual, []int{}, "OpEqual")
	op.createOpcodeDetails(factory, OpNotEqual, []int{}, "OpNotEqual")
	op.createOpcodeDetails(factory, OpMinus, []int{}, "OpMinus")
	op.createOpcodeDetails(factory, OpLNot, []int{}, "OpLNot")
	op.createOpcodeDetails(factory, OpJumpFalsy, []int{2}, "OpJumpFalsy")
	op.createOpcodeDetails(factory, OpAndJump, []int{2}, "OpAndJump")
	op.createOpcodeDetails(factory, OpOrJump, []int{2}, "OpOrJump")
	op.createOpcodeDetails(factory, OpJump, []int{2}, "OpJump")
	op.createOpcodeDetails(factory, OpNull, []int{}, "OpNull")
	op.createOpcodeDetails(factory, OpGetGlobal, []int{2}, "OpGetGlobal")
	op.createOpcodeDetails(factory, OpSetGlobal, []int{2}, "OpSetGlobal")
	op.createOpcodeDetails(factory, OpSetSelGlobal, []int{2, 1}, "OpSetSelGlobal")
	op.createOpcodeDetails(factory, OpArray, []int{2}, "OpArray")
	op.createOpcodeDetails(factory, OpMap, []int{2}, "OpMap")
	op.createOpcodeDetails(factory, OpStruct, []int{2}, "OpStruct")
	op.createOpcodeDetails(factory, OpImmutable, []int{}, "OpImmutable")
	op.createOpcodeDetails(factory, OpIndex, []int{}, "OpIndex")
	op.createOpcodeDetails(factory, OpSliceIndex, []int{}, "OpSliceIndex")
	op.createOpcodeDetails(factory, OpCall, []int{1, 1}, "OpCall")
	op.createOpcodeDetails(factory, OpReturn, []int{1}, "OpReturn")
	op.createOpcodeDetails(factory, OpGetLocal, []int{1}, "OpGetLocal")
	op.createOpcodeDetails(factory, OpSetLocal, []int{1}, "OpSetLocal")
	op.createOpcodeDetails(factory, OpDefineLocal, []int{1}, "OpDefineLocal")
	op.createOpcodeDetails(factory, OpSetSelLocal, []int{1, 1}, "OpSetSelLocal")
	op.createOpcodeDetails(factory, OpClosure, []int{2, 1}, "OpClosure")
	op.createOpcodeDetails(factory, OpGetFreePtr, []int{1}, "OpGetFreePtr")
	op.createOpcodeDetails(factory, OpGetFree, []int{1}, "OpGetFree")
	op.createOpcodeDetails(factory, OpSetFree, []int{1}, "OpSetFree")
	op.createOpcodeDetails(factory, OpGetLocalPtr, []int{1}, "OpGetLocalPtr")
	op.createOpcodeDetails(factory, OpSetSelFree, []int{1, 1}, "OpSetSelFree")
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
	if len(details.Operands()) == 0 {
		return 0
	}
	offset := 0
	for _, width := range details.Operands() {
		offset += width
	}
	return offset
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
		case 1:
			if offset >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset]))
		case 2:
			if offset+1 >= len(ins) {
				return nil, nil, 0
			}
			retOperands = append(retOperands, int(ins[offset+1])|int(ins[offset])<<8)
		}
		offset += width
	}
	return details.Operands(), retOperands, offset
}

// OpcodeNames returns the name of the provided opcode as a string.
func (op *Opcodes) OpcodeNames(opcode Opcode) string {
	details := op.OpcodeToDetails(opcode)
	return details.Name()
}

// CompileInstruction returns a bytecode for an opcode and the operands.
func (op *Opcodes) CompileInstruction(opcode Opcode, operands ...int) []byte {
	numOperands := op.OpcodeToOperands(opcode)
	totalLen := 1
	for _, w := range numOperands {
		totalLen += w
	}
	instruction := make([]byte, totalLen)
	instruction[0] = opcode
	offset := 1
	for i, o := range operands {
		width := numOperands[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			n := uint16(o)
			instruction[offset] = byte(n >> 8)
			instruction[offset+1] = byte(n)
		}
		offset += width
	}
	return instruction
}
