package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// IOpExecutor defines an interface for executing specific bytecode instructions within a virtual machine context.
// Opcode returns the bytecode.Opcode associated with the operation.
// Name returns the name of the operation represented by the executor.
// Operands retrieves the operands required for the operation's execution.
// Execute performs the operation within the provided virtual machine instance.
type IOpExecutor interface {
	Opcode() bytecode.Opcode

	Factory() *objects.GateKeeper

	Name() string

	Operands() []int

	Execute(vm *VM, operands *[]int)
}

// OpConstant represents an operation used to load a constant onto the stack.
type OpConstant struct {
	*bytecode.OpcodeDetails
}

// NewOpConstant creates a new OpConstant instance with opcode details initialized for the OpConstant operation.
func NewOpConstant(op *bytecode.Opcodes) *OpConstant {
	return &OpConstant{OpcodeDetails: op.OpcodeToDetails(bytecode.OpConstant)}
}

// Execute executes the OpConstant instruction in the virtual machine, pushing a global constant onto the stack.
func (op *OpConstant) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	cIdx := (*operands)[0]
	//if cIdx != int(v.currFrame.Get16(v.ip)) {
	//	panic("invalid constant index")
	//}
	glObj := v.constants.Get(uint(cIdx))
	v.stack.Push(glObj)
}

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	*bytecode.OpcodeDetails
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull(op *bytecode.Opcodes) *OpNull {
	return &OpNull{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNull)}
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	val := op.Factory().UndefinedValue()
	v.stack.Push(val)
}

// OpBinary represents a type that performs binary operations by extending bytecode.OpcodeDetails.
type OpBinary struct {
	*bytecode.OpcodeDetails
}

// NewOpBinary creates a new instance of OpBinary with its corresponding OpcodeDetails initialized.
func NewOpBinary(op *bytecode.Opcodes) *OpBinary {
	return &OpBinary{OpcodeDetails: op.OpcodeToDetails(bytecode.OpBinary)}
}

// Execute performs a binary operation using operands from the stack, updates the instruction pointer, and handles errors.
func (op *OpBinary) Execute(v *VM, operands *[]int) {
	// Operands Offset  1
	opcode := (*operands)[0]
	//if opcode != int(v.currFrame.Get8(v.ip)) {
	//	panic("invalid opcode")
	//}
	right := v.stack.Pop()
	left := v.stack.Pop()
	operator := objects.Operator(opcode)
	res, err := left.BinaryOp(v.currFrame.ID(), operator, right)
	if err != nil {
		v.SetError(err)
		return
	}
	v.stack.Push(res)
}

// OpEqual represents an operation that checks if two values are equal and updates the stack accordingly.
type OpEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpEqual creates and returns an instance of OpEqual, initialized with its corresponding opcode details.
func NewOpEqual(op *bytecode.Opcodes) *OpEqual {
	return &OpEqual{OpcodeDetails: op.OpcodeToDetails(bytecode.OpEqual)}
}

// Execute performs the equality comparison between the top two stack values and pushes the result (true or false) back onto the stack.
func (op *OpEqual) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := op.Factory().TrueValue()
	if left.Equals(right) {
		val = op.Factory().FalseValue()
	}
	v.stack.Push(val)
}

// OpNotEqual is a structure representing the "not equal (!=)" opcode operation in the virtual machine.
// It embeds OpcodeDetails to provide information about the opcode, including its identifier and operands.
type OpNotEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpNotEqual creates and returns a new instance of OpNotEqual with OpcodeDetails initialized from bytecode.
func NewOpNotEqual(op *bytecode.Opcodes) *OpNotEqual {
	return &OpNotEqual{OpcodeDetails: op.OpcodeToDetails(bytecode.OpNotEqual)}
}

// Execute evaluates whether the top two stack elements are unequal and pushes the result as a boolean onto the stack.
func (op *OpNotEqual) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := op.Factory().FalseValue()
	if left.Equals(right) {
		val = op.Factory().TrueValue()
	}
	v.stack.Push(val)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	*bytecode.OpcodeDetails
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop(op *bytecode.Opcodes) *OpPop {
	return &OpPop{OpcodeDetails: op.OpcodeToDetails(bytecode.OpPop)}
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	v.stack.Decrement()
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	*bytecode.OpcodeDetails
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue(op *bytecode.Opcodes) *OpTrue {
	return &OpTrue{OpcodeDetails: op.OpcodeToDetails(bytecode.OpTrue)}
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	val := op.Factory().TrueValue()
	v.stack.Push(val)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	*bytecode.OpcodeDetails
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse(op *bytecode.Opcodes) *OpFalse {
	return &OpFalse{OpcodeDetails: op.OpcodeToDetails(bytecode.OpFalse)}
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	val := op.Factory().FalseValue()
	v.stack.Push(val)
}

// OpLNot represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpLNot struct {
	*bytecode.OpcodeDetails
}

// NewOpLNot creates a new instance of OpLNot, representing a logical NOT operation (!).
func NewOpLNot(op *bytecode.Opcodes) *OpLNot {
	return &OpLNot{OpcodeDetails: op.OpcodeToDetails(bytecode.OpLNot)}
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpLNot) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	operand := v.stack.Pop()
	val := op.Factory().FalseValue()
	if operand.Boolean() {
		val = op.Factory().TrueValue()
	}
	v.stack.Push(val)
}

// OpBComplement represents an operation for performing a bitwise complement on an operand.
// It extends OpcodeDetails, inheriting its metadata and behaviors.
type OpBComplement struct {
	*bytecode.OpcodeDetails
}

// NewOpBComplement initializes and returns an OpBComplement instance with the corresponding OpcodeDetails configuration.
func NewOpBComplement(op *bytecode.Opcodes) *OpBComplement {
	return &OpBComplement{OpcodeDetails: op.OpcodeToDetails(bytecode.OpBComplement)}
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpBComplement) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.Factory().NewInt(v.currFrame.ID(), ^x.Value())
		v.stack.Push(res)
	default:
		v.SetError(fmt.Errorf("invalid operation: ^%s", operand.TypeName()))
		return
	}
}

// OpMinus represents an operation for negating a numeric value.
// It embeds OpcodeDetails, providing details such as the opcode, operands, and name.
type OpMinus struct {
	*bytecode.OpcodeDetails
}

// NewOpMinus creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpMinus(op *bytecode.Opcodes) *OpMinus {
	return &OpMinus{OpcodeDetails: op.OpcodeToDetails(bytecode.OpMinus)}
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpMinus) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := op.Factory().NewInt(v.currFrame.ID(), -x.Value())
		v.stack.Push(res)
	case *objects.Float:
		res := op.Factory().NewFloat(v.currFrame.ID(), -x.Value())
		v.stack.Push(res)
	default:
		v.SetError(fmt.Errorf("invalid operation: -%s", operand.TypeName()))
	}
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	*bytecode.OpcodeDetails
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding OpcodeDetails.
func NewOpJumpFalsy(op *bytecode.Opcodes) *OpJumpFalsy {
	return &OpJumpFalsy{OpcodeDetails: op.OpcodeToDetails(bytecode.OpJumpFalsy)}
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	obj := v.stack.Pop()
	if obj.Boolean() {
		pos := (*operands)[0]
		//if pos != int(v.currFrame.Get16(v.ip)) {
		//	panic("invalid jump position")
		//}
		v.SetIp(pos - 1)
	}
}

// OpAndJump represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpAndJump struct {
	*bytecode.OpcodeDetails
}

// NewOpAndJump creates and returns a new instance of OpAndJump, initializing it with details for the OpAndJump opcode.
func NewOpAndJump(op *bytecode.Opcodes) *OpAndJump {
	return &OpAndJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpAndJump)}
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpAndJump) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	obj := v.stack.Peek()
	if obj.Boolean() {
		pos := (*operands)[0]
		//if pos != int(v.currFrame.Get16(v.ip)) {
		//	panic("invalid jump position")
		//}
		v.SetIp(pos - 1)
	} else {
		v.stack.Decrement()
	}
}

// OpOrJump represents an operation that performs a logical OR and jumps based on the result.
type OpOrJump struct {
	*bytecode.OpcodeDetails
}

// NewOpOrJump creates and returns a new instance of OpOrJump, associated with the OpOrJump opcode and its details.
func NewOpOrJump(op *bytecode.Opcodes) *OpOrJump {
	return &OpOrJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpOrJump)}
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpOrJump) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		v.stack.Decrement()
	} else {
		pos := (*operands)[0]
		//if pos !=  int(v.currFrame.Get16(v.ip)) {
		//	panic("invalid jump position")
		//}
		v.SetIp(pos - 1)
	}
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	*bytecode.OpcodeDetails
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump(op *bytecode.Opcodes) *OpJump {
	return &OpJump{OpcodeDetails: op.OpcodeToDetails(bytecode.OpJump)}
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`VM`) to a calculated position in the frame.
func (op *OpJump) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	pos := (*operands)[0]
	//if pos !=  int(v.currFrame.Get16(v.ip)) {
	//	panic("invalid jump position")
	//}
	v.SetIp(pos - 1)
}

// OpSetGlobal represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpSetGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetGlobal creates and returns a new instance of OpSetGlobal with initialized OpcodeDetails.
func NewOpSetGlobal(op *bytecode.Opcodes) *OpSetGlobal {
	return &OpSetGlobal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetGlobal)}
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpSetGlobal) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	pos := (*operands)[0]
	//if pos !=  int(v.currFrame.Get16(v.ip)) {
	//	panic("invalid jump position")
	//}
	val := v.stack.Peek()
	v.global.Set(uint(pos), val)
}

// OpSetSelGlobal represents an operation for setting a global variable's value using selectors for indexing or access.
type OpSetSelGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelGlobal creates a new instance of OpSetSelGlobal with its corresponding OpcodeDetails initialized.
func NewOpSetSelGlobal(op *bytecode.Opcodes) *OpSetSelGlobal {
	return &OpSetSelGlobal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetSelGlobal)}
}

// Execute performs the operation defined by OpSetSelGlobal, updating the VM state and handling global index assignment.
func (op *OpSetSelGlobal) Execute(v *VM, operands *[]int) {
	// Operands Offset  3
	numSelectors := (*operands)[0]
	//if numSelectors !=  int(v.currFrame.Get8(v.ip)) {
	//	panic("invalid jump position")
	//}
	globalIndex := (*operands)[1]
	//if globalIndex !=  uint(v.currFrame.Get16(v.ip - 1)) {
	//	panic("invalid jump position")
	//}
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(numSelectors + 1)
	glObj := v.constants.Get(uint(globalIndex))
	if err := v.IndexAssign(v.currFrame.ID(), glObj, val, selectors); err != nil {
		v.SetError(err)
		return
	}
}

// OpGetGlobal represents an operation to retrieve a global variable in the virtual machine.
// It embeds OpcodeDetails for detailed opcode information.
type OpGetGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpGetGlobal creates a new instance of OpGetGlobal with its associated opcode details.
func NewOpGetGlobal(op *bytecode.Opcodes) *OpGetGlobal {
	return &OpGetGlobal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetGlobal)}
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGetGlobal) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	glIndex := (*operands)[0]
	//if glIndex != v.currFrame.Get16(v.ip) {
	//	panic("invalid jump position")
	//}
	glObj := v.global.Get(uint(glIndex))
	if glObj == nil {
		v.SetError(fmt.Errorf("undefined global: %d", glIndex))
		return
	}
	v.stack.Push(glObj)
}

// OpArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base OpcodeDetails for opcode, operands, and name information.
type OpArray struct {
	*bytecode.OpcodeDetails
}

// NewOpArray creates and returns a new instance of OpArray, initialized with details for the OpArray operation.
func NewOpArray(op *bytecode.Opcodes) *OpArray {
	return &OpArray{OpcodeDetails: op.OpcodeToDetails(bytecode.OpArray)}
}

// Execute processes the OpArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpArray) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	numElements := (*operands)[0]
	//if numElements != v.currFrame.Get16(v.ip {
	//	panic("invalid OpArray position")
	//}
	elements := v.stack.PopArrayElements(numElements)
	arr := op.Factory().NewArray(v.currFrame.ID(), elements)
	v.stack.Push(arr)
}

// OpMap is a wrapper around bytecode.OpcodeDetails, representing a map creation operation in bytecode execution.
type OpMap struct {
	*bytecode.OpcodeDetails
}

// NewOpMap initializes and returns a new instance of OpMap with its OpcodeDetails set to OpMap details.
func NewOpMap(op *bytecode.Opcodes) *OpMap {
	return &OpMap{OpcodeDetails: op.OpcodeToDetails(bytecode.OpMap)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpMap) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	numElements := (*operands)[0]
	//if numElements != int(v.currFrame.Get16(v.ip)) {
	//	panic("invalid OpMap position")
	//}
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(op.Factory().NewMap(v.currFrame.ID(), mElem))
}

// OpStruct is a wrapper around bytecode.OpcodeDetails, representing a struct creation operation in bytecode execution.
type OpStruct struct {
	*bytecode.OpcodeDetails
}

// NewOpStruct initializes and returns a new instance of OpStruct with its OpcodeDetails set to OpMap details.
func NewOpStruct(op *bytecode.Opcodes) *OpStruct {
	return &OpStruct{OpcodeDetails: op.OpcodeToDetails(bytecode.OpStruct)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpStruct) Execute(v *VM, operands *[]int) {
	// Operands Offset  2
	numElements := (*operands)[0]
	//if numElements != int(v.currFrame.Get16(v.ip)) {
	//	panic("invalid OpMap position")
	//}
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(op.Factory().NewStruct(v.currFrame.ID(), mElem))
}

// OpError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpError struct {
	*bytecode.OpcodeDetails
}

// NewOpError creates and returns a new instance of OpError with associated OpcodeDetails for the OpError opcode.
func NewOpError(op *bytecode.Opcodes) *OpError {
	return &OpError{OpcodeDetails: op.OpcodeToDetails(bytecode.OpError)}
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpError) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	value := v.stack.Peek()
	e := op.Factory().NewError(v.currFrame.ID(), value.String())
	v.stack.Set(e)
}

// OpImmutable represents an operation that creates immutable objects, inheriting details from OpcodeDetails.
type OpImmutable struct {
	*bytecode.OpcodeDetails
}

// NewOpImmutable creates a new instance of OpImmutable with details loaded from bytecode.OpcodeToDetails.
func NewOpImmutable(op *bytecode.Opcodes) *OpImmutable {
	return &OpImmutable{OpcodeDetails: op.OpcodeToDetails(bytecode.OpImmutable)}
}

// Execute processes the top element on the stack and converts it into an immutable version if it's an array or map.
func (op *OpImmutable) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	val := v.stack.Peek()
	switch value := val.(type) {
	case *objects.Array:
		obj := op.Factory().NewArrayImmutable(v.currFrame.ID(), value.Values())
		v.stack.Set(obj)
	case *objects.Map:
		obj := op.Factory().NewMapImmutable(v.currFrame.ID(), value.Values())
		v.stack.Set(obj)
	}
}

// OpIndex represents the operation for performing an indexing operation on a value.
type OpIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpIndex creates and returns a new instance of OpIndex initialized with its associated OpcodeDetails.
func NewOpIndex(op *bytecode.Opcodes) *OpIndex {
	return &OpIndex{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIndex)}
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndex) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	index := v.stack.Pop()
	left := v.stack.Pop()
	val, err := left.IndexGet(v.currFrame.ID(), index)
	if err != nil {
		if objects.Is(err, objects.ErrNotIndexable) {
			v.SetError(fmt.Errorf("not indexable: %s", index.TypeName()))
			return
		}
		if objects.Is(err, objects.ErrInvalidIndexType) {
			v.SetError(fmt.Errorf("invalid index type: %s", index.TypeName()))
			return
		}
		v.SetError(err)
		return
	}
	if val == nil {
		val = op.Factory().UndefinedValue()
	}
	v.stack.Push(val)
}

// OpSliceIndex represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds OpcodeDetails to inherit opcode, operand, and name information for execution and identification.
type OpSliceIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpSliceIndex creates a new instance of OpSliceIndex containing details for the slice indexing bytecode operation.
func NewOpSliceIndex(op *bytecode.Opcodes) *OpSliceIndex {
	return &OpSliceIndex{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSliceIndex)}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpSliceIndex) Execute(v *VM, _ *[]int) {
	// Operands Offset  0
	highStack := v.stack.Pop()
	lowStack := v.stack.Pop()
	leftStack := v.stack.Pop()
	lowIdx, highIdx, err := v.BoundsCheck(lowStack, highStack, int64(leftStack.Length()))
	if err != nil {
		v.SetError(err)
		return
	}
	var val objects.IObject = nil
	switch left := leftStack.(type) {
	case *objects.Array:
		val = op.Factory().NewArray(v.currFrame.ID(), left.Values()[lowIdx:highIdx])
	case *objects.ArrayImmutable:
		val = op.Factory().NewArray(v.currFrame.ID(), left.Values()[lowIdx:highIdx])
	case *objects.String:
		val = op.Factory().NewString(v.currFrame.ID(), left.Value()[lowIdx:highIdx])
	case *objects.Bytes:
		val = op.Factory().NewBytes(v.currFrame.ID(), left.Value()[lowIdx:highIdx])
	}
	if val != nil {
		v.stack.Push(val)
	}
}

// OpCall represents an operation code for invoking a function call in the virtual machine.
type OpCall struct {
	*bytecode.OpcodeDetails
}

// NewOpCall creates and returns a new instance of OpCall with initialized OpcodeDetails for the OpCall opcode.
func NewOpCall(op *bytecode.Opcodes) *OpCall {
	return &OpCall{OpcodeDetails: op.OpcodeToDetails(bytecode.OpCall)}
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	//fmt.Println("call", operands)
	numArgs := (*operands)[1]
	//if numArgs != int(v.currFrame.Get8(v.ip-1)) {
	//	panic("num args mismatch")
	//}
	value := v.stack.PeekOffset(-1 - numArgs)
	if !value.CanCall() {
		v.SetError(fmt.Errorf("%s is not callable: %s", value.String(), value.TypeName()))
		return
	}
	spread := v.currFrame.Get8(v.GetIp() + 2)
	if spread == 1 {
		arrObj := v.stack.Pop()
		switch z := arrObj.(type) {
		case *objects.Array:
			for _, item := range z.Values() {
				v.stack.Push(item)
			}
			numArgs += z.Length() - 1
		case *objects.ArrayImmutable:
			for _, item := range z.Values() {
				v.stack.Push(item)
			}
			numArgs += z.Length() - 1
		default:
			v.SetError(fmt.Errorf("not an array: %s", arrObj.TypeName()))
			return
		}
	}

	if callee, ok := value.(*objects.FuncCompiled); ok {
		if callee.VarArgs() {
			v.stack.PushVarArgs(v.currFrame.ID(), numArgs, callee.NumParameters()-1)
			numArgs = callee.NumParameters()
		}
		if numArgs != callee.NumParameters() {
			numParams := callee.NumParameters()
			if callee.VarArgs() {
				numParams--
			}
			v.SetError(fmt.Errorf("%s wrong number of arguments: want>=%d, got=%d", callee.Name(), numParams, numArgs))
			return
		}
		// Frame setup
		v.currFrame = v.frames.Get()
		v.frames.Next()
		bp := v.stack.StackPointer() - numArgs
		v.currFrame.Bind(v.GetIp(), callee, bp)
		// Reserve space for *all* local variables of the new function
		// by simply advancing the stack pointer.
		// This ensures that space for temporary calculations starts *after*
		// the space reserved for local variables, avoiding collisions.
		v.stack.SetStackPointer(v.stack.StackPointer() + callee.NumLocals())
		v.ReseIp()
	} else {
		var args []objects.IObject
		args = append(args, v.stack.PeekArrayObject(numArgs)...)
		ret, err := value.Call(v.currFrame.ID(), args...)
		// Cleans the stack from the function and its arguments
		v.stack.DecrementCount(numArgs + 1)
		if err != nil {
			if objects.Is(err, objects.ErrWrongNumArguments) {
				v.SetError(fmt.Errorf("wrong number of arguments in call to '%s'", value.TypeName()))
			} else {
				v.SetError(err)
			}
			return
		}
		if ret == nil {
			v.stack.Push(op.Factory().UndefinedValue())
		} else {
			v.stack.Push(ret)
		}
	}
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.OpcodeDetails.
type OpReturn struct {
	*bytecode.OpcodeDetails
}

// NewOpReturn creates a new instance of OpReturn with its OpcodeDetails initialized for the OpReturn operation.
func NewOpReturn(op *bytecode.Opcodes) *OpReturn {
	return &OpReturn{OpcodeDetails: op.OpcodeToDetails(bytecode.OpReturn)}
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the VM.
func (op *OpReturn) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	//if numReturnVals != int(v.currFrame.Get8(v.ip)) {
	//	panic("num return vals mismatch")
	//}

	// collect return values from the stack using Pop(),
	// this is necessary to uncover the underlying values.
	var returnValues []objects.IObject
	if numReturnVals := (*operands)[0]; numReturnVals > 0 {
		returnValues = make([]objects.IObject, (*operands)[0])
		for i := 0; i < numReturnVals; i++ {
			returnValues[i] = v.stack.Pop()
		}
	}

	shutdown := false
	prevIp := v.currFrame.SavedIP()
	leavingFrameBasePointer := v.currFrame.BasePointer()
	//v.stack.ReleaseObjects(leavingFrameBasePointer, v.stack.StackPointer())
	if v.frames.Index() > 1 {
		v.frames.Previous()
		v.currFrame = v.frames.GetPrev()
		v.SetIp(prevIp)
	} else {
		shutdown = true
	}
	v.stack.SetStackPointer(leavingFrameBasePointer)
	// push return values onto the new stack (caller's stack).
	if lRet := len(returnValues); lRet > 0 {
		// iterate over the slice in reverse to restore the original order.
		for i := lRet - 1; i >= 0; i-- {
			v.stack.Push(returnValues[i])
		}
	} else {
		v.stack.Push(op.Factory().UndefinedValue())
	}

	if shutdown {
		v.Shutdown()
	}
}

// OpDefineLocal represents the opcode for defining a new local variable within the current frame's scope.
type OpDefineLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpDefineLocal creates a new instance of OpDefineLocal with its associated opcode details.
func NewOpDefineLocal(op *bytecode.Opcodes) *OpDefineLocal {
	return &OpDefineLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpDefineLocal)}
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpDefineLocal) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("invalid OpDefineLocal position")
	//}
	val := v.stack.Peek()
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, val)
}

// OpSetLocal represents an operation to set the value of a local variable within the current frame.
// It embeds OpcodeDetails for opcode-specific information such as name, operands, and code.
type OpSetLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetLocal initializes and returns a new instance of OpSetLocal with associated opcode details.
func NewOpSetLocal(op *bytecode.Opcodes) *OpSetLocal {
	return &OpSetLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetLocal)}
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpSetLocal) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local index mismatch")
	//}

	val := v.stack.Peek()
	destSlot := v.currFrame.BasePointer() + localIndex
	existingValue := v.stack.PeekAbsolute(destSlot)
	if obj, ok := existingValue.(*objects.ObjectPointer); ok {
		obj.SetValue(val)
	} else {
		v.stack.SetAbsolute(destSlot, val)
	}
}

// OpSetSelLocal represents an operation for setting a local variable using selectors in the virtual machine.
// It embeds OpcodeDetails to utilize its properties like opcode, name, and operands.
type OpSetSelLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelLocal creates and returns a new instance of the OpSetSelLocal operation executor.
func NewOpSetSelLocal(op *bytecode.Opcodes) *OpSetSelLocal {
	return &OpSetSelLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetSelLocal)}
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpSetSelLocal) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	numSelectors := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpSetSelLocal mismatch")
	//}
	localIndex := (*operands)[1]
	//if localIndex != int(v.currFrame.Get8(v.ip - 1)) {
	//	panic("local OpSetSelLocal mismatch")
	//}
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(numSelectors + 1)
	dst := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := dst.(*objects.ObjectPointer); ok {
		dst = *obj.Value()
	}
	if err := v.IndexAssign(v.currFrame.ID(), dst, val, selectors); err != nil {
		v.SetError(err)
		return
	}
}

// OpGetLocal represents an operation to retrieve a local variable from the stack using its index.
type OpGetLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocal creates a new OpGetLocal instance and initializes it with details for the OpGetLocal opcode.
func NewOpGetLocal(op *bytecode.Opcodes) *OpGetLocal {
	return &OpGetLocal{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetLocal)}
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpGetLocal) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpGetLocal mismatch")
	//}
	val := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.stack.Push(val)
}

// OpClosure represents a closure operation that creates a new closure in the virtual machine.
type OpClosure struct {
	*bytecode.OpcodeDetails
}

// NewOpClosure returns a new instance of OpClosure initialized with the details of the OpClosure opcode.
func NewOpClosure(op *bytecode.Opcodes) *OpClosure {
	return &OpClosure{OpcodeDetails: op.OpcodeToDetails(bytecode.OpClosure)}
}

// Execute performs the operation associated with the OpClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpClosure) Execute(v *VM, operands *[]int) {
	// Operands Offset 3
	numFree := (*operands)[0]
	//if numFree != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpClosure mismatch")
	//}
	constIndex := (*operands)[1]
	//if constIndex != int(v.currFrame.Get16(v.ip - 1)) {
	//	panic("local OpClosure mismatch")
	//}
	glObj := v.constants.Get(uint(constIndex))
	fn, ok := glObj.(*objects.FuncCompiled)
	if !ok {
		v.SetError(fmt.Errorf("not a function: %s", fn.TypeName()))
		return
	}
	free := make([]*objects.ObjectPointer, numFree)
	for i := 0; i < numFree; i++ {
		o := v.stack.PeekOffset(-numFree + i)
		switch freeVar := o.(type) {
		case *objects.ObjectPointer:
			free[i] = freeVar
		default:
			t := v.stack.PeekOffset(-numFree + i)
			obj := op.Factory().NewObjectPointer(v.currFrame.ID(), &t)
			ptr, ok := obj.(*objects.ObjectPointer)
			if !ok {
				v.SetError(fmt.Errorf("not a pointer: %s", t.TypeName()))
				return
			}
			free[i] = ptr
		}
	}
	v.stack.DecrementCount(numFree)
	cl := op.Factory().NewFuncCompiled(v.currFrame.ID(), "closure", fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	v.stack.Push(cl)
}

// OpGetFreePtr represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds OpcodeDetails, which provides opcode metadata such as identifier, operands, and name.
type OpGetFreePtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFreePtr creates a new instance of OpGetFreePtr initialized with the corresponding OpcodeDetails.
func NewOpGetFreePtr(op *bytecode.Opcodes) *OpGetFreePtr {
	return &OpGetFreePtr{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetFreePtr)}
}

// Execute executes the OpGetFreePtr operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpGetFreePtr) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	freeIndex := (*operands)[0]
	//if freeIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpGetFreePtr mismatch")
	//}
	val := v.currFrame.FreeVarsIndex(freeIndex)
	v.stack.Push(val)
}

// OpGetFree represents an operation to retrieve a free variable in a closure during execution.
type OpGetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFree creates and returns a new instance of OpGetFree, initializing its OpcodeDetails using bytecode metadata.
func NewOpGetFree(op *bytecode.Opcodes) *OpGetFree {
	return &OpGetFree{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetFree)}
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpGetFree) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	freeIndex := (*operands)[0]
	//if freeIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpGetFree mismatch")
	//}
	val := *v.currFrame.FreeVarsIndex(freeIndex).Value()
	v.stack.Push(val)
}

// OpSetFree represents an operation to set the value of a free variable within a closure's environment.
type OpSetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpSetFree creates and returns a new instance of OpSetFree initialized with its corresponding OpcodeDetails.
func NewOpSetFree(op *bytecode.Opcodes) *OpSetFree {
	return &OpSetFree{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetFree)}
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpSetFree) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	freeIndex := (*operands)[0]
	//if freeIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpSetFree mismatch")
	//}
	o := v.stack.Pop()
	v.currFrame.FreeVarsIndex(freeIndex).SetValue(o)
}

// OpGetLocalPtr retrieves a local variable as a pointer using its index within the current frame.
type OpGetLocalPtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocalPtr creates and returns a new instance of OpGetLocalPtr, initializing its OpcodeDetails.
func NewOpGetLocalPtr(op *bytecode.Opcodes) *OpGetLocalPtr {
	return &OpGetLocalPtr{OpcodeDetails: op.OpcodeToDetails(bytecode.OpGetLocalPtr)}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpGetLocalPtr) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpGetLocalPtr mismatch")
	//}
	sp := v.currFrame.BasePointer() + localIndex
	val := v.stack.PeekAbsolute(sp)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		v.stack.Push(obj)
		return
	}
	freeVar := op.Factory().NewObjectPointer(v.currFrame.ID(), &val)
	v.stack.SetAbsolute(sp, freeVar)
	v.stack.Push(freeVar)
}

// OpSetSelFree represents an operation to set a free variable's value using selectors.
type OpSetSelFree struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelFree creates a new instance of OpSetSelFree with initialized OpcodeDetails referencing OpSetSelFree.
func NewOpSetSelFree(op *bytecode.Opcodes) *OpSetSelFree {
	return &OpSetSelFree{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSetSelFree)}
}

// Execute updates the instruction pointer, retrieves operands, processes selectors, and performs indexed assignment in the VM.
func (op *OpSetSelFree) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	numSelectors := (*operands)[0]
	//if numSelectors != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpSetSelFree mismatch")
	//}
	freeIndex := (*operands)[1]
	//if freeIndex != int(v.currFrame.Get8(v.ip - 1)) {
	//	panic("local OpSetSelFree mismatch")
	//}
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(numSelectors + 1)
	fvi := v.currFrame.FreeVarsIndex(freeIndex)
	if err := v.IndexAssign(v.currFrame.ID(), *fvi.Value(), val, selectors); err != nil {
		v.SetError(err)
		return
	}
}

// OpIteratorInit represents an operation that initializes an iterator over an iterable object.
// It embeds OpcodeDetails for additional opcode-specific metadata.
type OpIteratorInit struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit(op *bytecode.Opcodes) *OpIteratorInit {
	return &OpIteratorInit{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorInit)}
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpIteratorInit mismatch")
	//}
	iterable := v.stack.Pop()
	if !iterable.CanIterate() {
		v.SetError(fmt.Errorf("not iterable: %s", iterable.TypeName()))
		return
	}
	iterator := iterable.Iterate(v.currFrame.ID())
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, iterator)
}

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext(op *bytecode.Opcodes) *OpIteratorNext {
	return &OpIteratorNext{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorNext)}
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpIteratorNext mismatch")
	//}

	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	if iterator.Next() {
		v.stack.Push(op.Factory().TrueValue())
	} else {
		v.stack.Push(op.Factory().FalseValue())
	}
}

// OpIteratorKey wraps bytecode.OpcodeDetails to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey(op *bytecode.Opcodes) *OpIteratorKey {
	return &OpIteratorKey{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorKey)}
}

// Execute processes the "iterator key" operation, retrieves the iterator key, and pushes it onto the VM stack.
func (op *OpIteratorKey) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpIteratorKey mismatch")
	//}

	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.stack.Push(iterator.Key(v.currFrame.ID()))
}

// OpIteratorValue retrieves the value from the current iterator position.
// It embeds OpcodeDetails, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated OpcodeDetails initialized.
func NewOpIteratorValue(op *bytecode.Opcodes) *OpIteratorValue {
	return &OpIteratorValue{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIteratorValue)}
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(v *VM, operands *[]int) {
	// Operands Offset 1
	localIndex := (*operands)[0]
	//if localIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("local OpIteratorValue mismatch")
	//}
	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.SetError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.stack.Push(iterator.Value(v.currFrame.ID()))
}

// OpReferences extends OpcodeDetails to represent operations specifically related to reference handling in the bytecode.
type OpReferences struct {
	*bytecode.OpcodeDetails
}

// NewOpReferences initializes a new OpReferences instance with corresponding OpcodeDetails from the bytecode package.
func NewOpReferences(op *bytecode.Opcodes) *OpReferences {
	return &OpReferences{OpcodeDetails: op.OpcodeToDetails(bytecode.OpReferences)}
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpReferences) Execute(v *VM, operands *[]int) {
	// Operands Offset 2
	nameIndex := (*operands)[0]
	//if nameIndex != int(v.currFrame.Get16(v.ip)) {
	//	panic("name index mismatch: %d != %d", nameIndex, int(v.currFrame.Get16(v.ip)))
	//}
	//nameIndex := int(v.currFrame.Get16(v.ip))
	if nameIndex < 0 || nameIndex >= len(v.references) {
		v.SetError(fmt.Errorf("invalid attribute index %d", nameIndex))
		return
	}
	symbol := v.references[nameIndex]
	v.stack.Push(symbol)
}

// OpIntOp extends OpcodeDetails and represents integer operations performed on a virtual machine.
type OpIntOp struct {
	*bytecode.OpcodeDetails
}

// NewOpIntOp initializes and returns a new instance of OpIntOp with relevant opcode details provided by bytecode.Opcodes.
func NewOpIntOp(op *bytecode.Opcodes) *OpIntOp {
	return &OpIntOp{OpcodeDetails: op.OpcodeToDetails(bytecode.OpIntOp)}
}

// Execute performs a specified binary operation between two integers from the stack and stores the result in a destination slot.
// It retrieves operands, validates types, and executes the operation, setting an error on unsupported cases or type mismatches.
func (op *OpIntOp) Execute(v *VM, operands *[]int) {
	dstObj := v.stack.PeekAbsolute((*operands)[0])
	binaryOp := objects.Operator((*operands)[1])
	dst, ok := dstObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("dst expected int, got %s", dstObj.TypeName()))
		return
	}
	rhsObj := v.stack.Pop()
	rhs, ok := rhsObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("rhs expected int, got %s", rhsObj.TypeName()))
		return
	}
	lhsObj := v.stack.Pop()
	lhs, ok := lhsObj.(*objects.Int)
	if !ok {
		v.SetError(fmt.Errorf("lhs expected int, got %s", lhsObj.TypeName()))
		return
	}
	result, err := v.BinaryOpInt64(binaryOp, lhs.Value(), rhs.Value())
	if err != nil {
		v.SetError(err)
	}
	dst.SetValue(result)
}

// OpSuspend represents an operation that suspends the execution of the virtual machine.
type OpSuspend struct {
	*bytecode.OpcodeDetails
}

// NewOpSuspend creates and returns a new OpSuspend instance with opcode details initialized for the suspend operation.
func NewOpSuspend(op *bytecode.Opcodes) *OpSuspend {
	return &OpSuspend{OpcodeDetails: op.OpcodeToDetails(bytecode.OpSuspend)}
}

// Execute performs the suspend operation on the given virtual machine by setting its shutdown state to true.
func (op *OpSuspend) Execute(v *VM, _ *[]int) {
	// Operands Offset 0
	v.shutdown = true
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*bytecode.OpcodeDetails
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding OpcodeDetails configuration set.
func NewOpUnknown(op *bytecode.Opcodes) *OpUnknown {
	return &OpUnknown{OpcodeDetails: op.OpcodeToDetails(bytecode.OpUnknown)}
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(v *VM, operands *[]int) {
	// Operands Offset 0
	//currentIp := int(v.currFrame.Get8(v.ip))
	//if nameIndex != int(v.currFrame.Get8(v.ip)) {
	//	panic("name index mismatch: %d != %d", nameIndex, int(v.currFrame.Get16(v.ip)))
	//}
	v.SetError(fmt.Errorf("unknown opcode at: %d, %v", v.GetIp(), operands))
}
