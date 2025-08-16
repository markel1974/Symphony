package bytecode

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
	OpLNot                        // Logical not !
	OpJumpFalsy                   // Jump if falsy
	OpAndJump                     // Logical AND jump
	OpOrJump                      // Logical OR jump
	OpJump                        // Jump
	OpNull                        // Push null
	OpArray                       // Array object
	OpMap                         // Map object
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
	OpGetBuiltin                  // Get builtin function
	OpClosure                     // Push closure
	OpIteratorInit                // Iterator init
	OpIteratorNext                // Iterator next
	OpIteratorKey                 // Iterator key
	OpIteratorValue               // Iterator value
	OpBinary                      // Binary operation
	OpReferences
	OpSuspend // Suspend VM
	OpError   // Error object
	OpUnknown
)

// OpcodeDetails represents the details of an opcode, including its identifier, its operands, and its name.
type OpcodeDetails struct {
	opcode   Opcode
	operands []int
	name     string
}

// NewOpcodeDetails creates a new OpcodeDetails instance, initializing its opcode, operands, and name fields.
func NewOpcodeDetails(opcode Opcode, operands []int, name string) *OpcodeDetails {
	od := &OpcodeDetails{
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

// Name returns the name of the opcode as a string.
func (od *OpcodeDetails) Name() string {
	return od.name
}

// Operands retrieves the list of integer operands associated with the OpcodeDetails instance.
func (od *OpcodeDetails) Operands() []int {
	return od.operands
}

// _opcodesDetail holds details for all defined opcodes, including their operands and names, indexed by the opcode value.
var _opcodesDetail []*OpcodeDetails

// init initializes the opcode details by creating and populating the `_opcodesDetail` slice with `OpcodeDetails` instances.
func init() {
	_opcodesDetail = make([]*OpcodeDetails, OpcodesLen)
	for i := range _opcodesDetail {
		_opcodesDetail[i] = NewOpcodeDetails(OpUnknown, []int{}, "OpUnknown")
	}
	createOpcodeDetails(OpConstant, []int{2}, "OpConstant")
	createOpcodeDetails(OpPop, []int{}, "OpPop")
	createOpcodeDetails(OpTrue, []int{}, "OpTrue")
	createOpcodeDetails(OpFalse, []int{}, "OpFalse")
	createOpcodeDetails(OpBComplement, []int{}, "OpBComplement")
	createOpcodeDetails(OpEqual, []int{}, "OpEqual")
	createOpcodeDetails(OpNotEqual, []int{}, "OpNotEqual")
	createOpcodeDetails(OpMinus, []int{}, "OpMinus")
	createOpcodeDetails(OpLNot, []int{}, "OpLNot")
	createOpcodeDetails(OpJumpFalsy, []int{2}, "OpJumpFalsy")
	createOpcodeDetails(OpAndJump, []int{2}, "OpAndJump")
	createOpcodeDetails(OpOrJump, []int{2}, "OpOrJump")
	createOpcodeDetails(OpJump, []int{2}, "OpJump")
	createOpcodeDetails(OpNull, []int{}, "OpNull")
	createOpcodeDetails(OpGetGlobal, []int{2}, "OpGetGlobal")
	createOpcodeDetails(OpSetGlobal, []int{2}, "OpSetGlobal")
	createOpcodeDetails(OpSetSelGlobal, []int{2, 1}, "OpSetSelGlobal")
	createOpcodeDetails(OpArray, []int{2}, "OpArray")
	createOpcodeDetails(OpMap, []int{2}, "OpMap")
	createOpcodeDetails(OpImmutable, []int{}, "OpImmutable")
	createOpcodeDetails(OpIndex, []int{}, "OpIndex")
	createOpcodeDetails(OpSliceIndex, []int{}, "OpSliceIndex")
	createOpcodeDetails(OpCall, []int{1, 1}, "OpCall")
	createOpcodeDetails(OpReturn, []int{1}, "OpReturn")
	createOpcodeDetails(OpGetLocal, []int{1}, "OpGetLocal")
	createOpcodeDetails(OpSetLocal, []int{1}, "OpSetLocal")
	createOpcodeDetails(OpDefineLocal, []int{1}, "OpDefineLocal")
	createOpcodeDetails(OpSetSelLocal, []int{1, 1}, "OpSetSelLocal")
	createOpcodeDetails(OpGetBuiltin, []int{1}, "OpGetBuiltin")
	createOpcodeDetails(OpClosure, []int{2, 1}, "OpClosure")
	createOpcodeDetails(OpGetFreePtr, []int{1}, "OpGetFreePtr")
	createOpcodeDetails(OpGetFree, []int{1}, "OpGetFree")
	createOpcodeDetails(OpSetFree, []int{1}, "OpSetFree")
	createOpcodeDetails(OpGetLocalPtr, []int{1}, "OpGetLocalPtr")
	createOpcodeDetails(OpSetSelFree, []int{1, 1}, "OpSetSelFree")
	createOpcodeDetails(OpIteratorInit, []int{1}, "OpIteratorInit")
	createOpcodeDetails(OpIteratorNext, []int{1}, "OpIteratorNext")
	createOpcodeDetails(OpIteratorKey, []int{1}, "OpIteratorKey")
	createOpcodeDetails(OpIteratorValue, []int{1}, "OpIteratorValue")
	createOpcodeDetails(OpBinary, []int{1}, "OpBinary")
	createOpcodeDetails(OpReferences, []int{2}, "OpReferences")
	createOpcodeDetails(OpSuspend, []int{}, "OpSuspend")
	createOpcodeDetails(OpError, []int{}, "OpError")
}

// createOpcodeDetails associates opcode details with a specific opcode, storing it in a global lookup by applying a mask.
func createOpcodeDetails(opcode Opcode, operands []int, name string) {
	od := NewOpcodeDetails(opcode, operands, name)
	_opcodesDetail[od.opcode&OpcodesMask] = od
}

// OpcodeToDetails retrieves the OpcodeDetails corresponding to the given opcode by applying the OpcodesMask.
func OpcodeToDetails(opcode Opcode) *OpcodeDetails {
	return _opcodesDetail[opcode&OpcodesMask]
}

// OpcodeToOperands retrieves the operand widths for a given opcode from its details.
func OpcodeToOperands(opcode Opcode) []int {
	details := OpcodeToDetails(opcode)
	return details.Operands()
}

// OpcodeToOperandsOffset calculates the total byte offset for the operands of a given opcode.
func OpcodeToOperandsOffset(opcode Opcode) int {
	details := OpcodeToDetails(opcode)
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
func OpcodeToOperandsDetails(opcode Opcode, ins []byte) ([]int, []int, int) {
	details := OpcodeToDetails(opcode)
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
func OpcodeNames(opcode Opcode) string {
	details := OpcodeToDetails(opcode)
	return details.Name()
}
