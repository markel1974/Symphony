package bytecode

// Opcode is a type alias for byte, representing an operation code in the bytecode instruction set.
type Opcode = byte

// OpConstant represents the opcode for loading a constant.
// OpBComplement represents the opcode for performing a bitwise complement operation.
// OpPop represents the opcode for popping the top of the stack.
// OpTrue represents the opcode for pushing a true value onto the stack.
// OpFalse represents the opcode for pushing a false value onto the stack.
// OpEqual represents the opcode for evaluating equality (==).
// OpNotEqual represents the opcode for evaluating inequality (!=).
// OpMinus represents the opcode for the subtraction operation (-).
// OpLNot represents the opcode for the logical NOT (!) operation.
// OpJumpFalsy represents the opcode for conditional jumping if the value is falsy.
// OpAndJump represents the opcode for a logical AND operation with jump.
// OpOrJump represents the opcode for a logical OR operation with jump.
// OpJump represents the opcode for an unconditional jump.
// OpNull represents the opcode for pushing a null value onto the stack.
// OpArray represents the opcode for creating an array object.
// OpMap represents the opcode for creating a map object.
// OpImmutable represents the opcode for creating an immutable object.
// OpIndex represents the opcode for performing index operations.
// OpSliceIndex represents the opcode for performing slice operations.
// OpCall represents the opcode for calling a function.
// OpReturn represents the opcode for returning from a function.
// OpGetGlobal represents the opcode for retrieving a global variable.
// OpSetGlobal represents the opcode for setting a global variable.
// OpSetSelGlobal represents the opcode for setting a global variable using a selector.
// OpGetLocal represents the opcode for retrieving a local variable.
// OpSetLocal represents the opcode for setting a local variable.
// OpDefineLocal represents the opcode for defining a local variable.
// OpSetSelLocal represents the opcode for setting a local variable using a selector.
// OpGetFreePtr represents the opcode for retrieving a free variable pointer object.
// OpGetFree represents the opcode for retrieving free variables.
// OpSetFree represents the opcode for setting free variables.
// OpGetLocalPtr represents the opcode for retrieving a local variable as a pointer.
// OpSetSelFree represents the opcode for setting free variables using a selector.
// OpGetBuiltin represents the opcode for retrieving a builtin function.
// OpClosure represents the opcode for pushing a closure onto the stack.
// OpIteratorInit represents the opcode for initializing an iterator.
// OpIteratorNext represents the opcode for advancing an iterator.
// OpIteratorKey represents the opcode for retrieving the current key of the iterator.
// OpIteratorValue represents the opcode for retrieving the current value of the iterator.
// OpBinaryOp represents the opcode for performing a binary operation.
// OpReferences represents the opcode for referencing an object.
// OpSuspend represents the opcode for suspending the virtual machine execution.
// OpError represents the opcode for creating an error object (must be the last opcode).
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
	OpBinaryOp                    // Binary operation
	OpReferences
	OpSuspend // Suspend VM
	OpError   // Error object - Must be last opcode
)

// OpcodeDetails represents details about an opcode, including its name and the operands it operates on.
type OpcodeDetails struct {
	operands []int
	name     string
}

// NewOpcodeDetails creates and returns a new OpcodeDetails object with the given operands and name.
func NewOpcodeDetails(operands []int, name string) *OpcodeDetails {
	od := &OpcodeDetails{
		operands: operands,
		name:     name,
	}
	return od
}

// _opcodesDetail holds details for each opcode, including operand specifications and their names.
var _opcodesDetail []*OpcodeDetails

// init initializes the opcode details mapping for various operation codes and their corresponding operand specifications.
func init() {
	_opcodesDetail = make([]*OpcodeDetails, OpError+1)
	_opcodesDetail[OpConstant] = NewOpcodeDetails([]int{2}, "OpConstant")
	_opcodesDetail[OpPop] = NewOpcodeDetails([]int{}, "OpPop")
	_opcodesDetail[OpTrue] = NewOpcodeDetails([]int{}, "OpTrue")
	_opcodesDetail[OpFalse] = NewOpcodeDetails([]int{}, "OpFalse")
	_opcodesDetail[OpBComplement] = NewOpcodeDetails([]int{}, "OpBComplement")
	_opcodesDetail[OpEqual] = NewOpcodeDetails([]int{}, "OpEqual")
	_opcodesDetail[OpNotEqual] = NewOpcodeDetails([]int{}, "OpNotEqual")
	_opcodesDetail[OpMinus] = NewOpcodeDetails([]int{}, "OpMinus")
	_opcodesDetail[OpLNot] = NewOpcodeDetails([]int{}, "OpLNot")
	_opcodesDetail[OpJumpFalsy] = NewOpcodeDetails([]int{2}, "OpJumpFalsy")
	_opcodesDetail[OpAndJump] = NewOpcodeDetails([]int{2}, "OpAndJump")
	_opcodesDetail[OpOrJump] = NewOpcodeDetails([]int{2}, "OpOrJump")
	_opcodesDetail[OpJump] = NewOpcodeDetails([]int{2}, "OpJump")
	_opcodesDetail[OpNull] = NewOpcodeDetails([]int{}, "OpNull")
	_opcodesDetail[OpGetGlobal] = NewOpcodeDetails([]int{2}, "OpGetGlobal")
	_opcodesDetail[OpSetGlobal] = NewOpcodeDetails([]int{2}, "OpSetGlobal")
	_opcodesDetail[OpSetSelGlobal] = NewOpcodeDetails([]int{2, 1}, "OpSetSelGlobal")
	_opcodesDetail[OpArray] = NewOpcodeDetails([]int{2}, "OpArray")
	_opcodesDetail[OpMap] = NewOpcodeDetails([]int{2}, "OpMap")
	_opcodesDetail[OpImmutable] = NewOpcodeDetails([]int{}, "OpImmutable")
	_opcodesDetail[OpIndex] = NewOpcodeDetails([]int{}, "OpIndex")
	_opcodesDetail[OpSliceIndex] = NewOpcodeDetails([]int{}, "OpSliceIndex")
	_opcodesDetail[OpCall] = NewOpcodeDetails([]int{1, 1}, "OpCall")
	_opcodesDetail[OpReturn] = NewOpcodeDetails([]int{1}, "OpReturn")
	_opcodesDetail[OpGetLocal] = NewOpcodeDetails([]int{1}, "OpGetLocal")
	_opcodesDetail[OpSetLocal] = NewOpcodeDetails([]int{1}, "OpSetLocal")
	_opcodesDetail[OpDefineLocal] = NewOpcodeDetails([]int{1}, "OpDefineLocal")
	_opcodesDetail[OpSetSelLocal] = NewOpcodeDetails([]int{1, 1}, "OpSetSelLocal")
	_opcodesDetail[OpGetBuiltin] = NewOpcodeDetails([]int{1}, "OpGetBuiltin")
	_opcodesDetail[OpClosure] = NewOpcodeDetails([]int{2, 1}, "OpClosure")
	_opcodesDetail[OpGetFreePtr] = NewOpcodeDetails([]int{1}, "OpGetFreePtr")
	_opcodesDetail[OpGetFree] = NewOpcodeDetails([]int{1}, "OpGetFree")
	_opcodesDetail[OpSetFree] = NewOpcodeDetails([]int{1}, "OpSetFree")
	_opcodesDetail[OpGetLocalPtr] = NewOpcodeDetails([]int{1}, "OpGetLocalPtr")
	_opcodesDetail[OpSetSelFree] = NewOpcodeDetails([]int{1, 1}, "OpSetSelFree")
	_opcodesDetail[OpIteratorInit] = NewOpcodeDetails([]int{1}, "OpIteratorInit")
	_opcodesDetail[OpIteratorNext] = NewOpcodeDetails([]int{1}, "OpIteratorNext")
	_opcodesDetail[OpIteratorKey] = NewOpcodeDetails([]int{1}, "OpIteratorKey")
	_opcodesDetail[OpIteratorValue] = NewOpcodeDetails([]int{1}, "OpIteratorValue")
	_opcodesDetail[OpBinaryOp] = NewOpcodeDetails([]int{1}, "OpBinaryOp")
	_opcodesDetail[OpReferences] = NewOpcodeDetails([]int{2}, "OpReferences")
	_opcodesDetail[OpSuspend] = NewOpcodeDetails([]int{}, "OpSuspend")
	_opcodesDetail[OpError] = NewOpcodeDetails([]int{}, "OpError")
}

func OpcodeToDetails(opcode Opcode) *OpcodeDetails {
	if opcode < 0 || int(opcode) > len(_opcodesDetail) {
		return nil
	}
	return _opcodesDetail[opcode]
}

// OpcodeToOperands returns the operand widths associated with the given opcode.
// If the opcode is invalid, it returns nil.
func OpcodeToOperands(opcode Opcode) []int {
	details := OpcodeToDetails(opcode)
	if details == nil {
		return nil
	}
	opOperands := details.operands
	return opOperands
}

// OpcodeToOperandsOffset calculates the cumulative operand width for a given opcode using OpcodeToOperands to determine widths.
func OpcodeToOperandsOffset(opcode Opcode) int {
	details := OpcodeToDetails(opcode)
	if details == nil {
		return 0
	}
	if len(details.operands) == 0 {
		return 0
	}
	offset := 0
	for _, width := range details.operands {
		offset += width
	}
	return offset
}

// OpcodeToOperandsDetails extracts operand widths, operand values, and total offset from an instruction byte slice.
func OpcodeToOperandsDetails(opcode Opcode, ins []byte) ([]int, []int, int) {
	details := OpcodeToDetails(opcode)
	if details == nil {
		return nil, nil, 0
	}
	if len(details.operands) == 0 {
		return nil, nil, 0
	}
	var retOperands []int
	var offset int
	for _, width := range details.operands {
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
	return details.operands, retOperands, offset
}

// OpcodeNames returns the name of the provided opcode as a string, or an empty string if the opcode is invalid.
func OpcodeNames(opcode Opcode) string {
	if opcode < 0 || int(opcode) > len(_opcodesDetail) {
		return ""
	}
	return _opcodesDetail[opcode].name
}
