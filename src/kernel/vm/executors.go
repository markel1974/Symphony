package vm

import (
	"fmt"
	"log"

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
	Name() string
	Operands() []int
	Execute(v *VM)
}

// OpConstant represents an operation used to load a constant onto the stack.
type OpConstant struct {
	*bytecode.OpcodeDetails
}

// NewOpConstant creates a new OpConstant instance with opcode details initialized for the OpConstant operation.
func NewOpConstant() *OpConstant {
	return &OpConstant{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpConstant)}
}

// Execute executes the OpConstant instruction in the virtual machine, pushing a global constant onto the stack.
func (op *OpConstant) Execute(v *VM) {
	v.ip += 2
	cIdx := v.currFrame.Pos(v.ip, v.ip-1)
	glObj := v.globals.Get(uint(cIdx))
	v.stack.Push(glObj)
}

// OpNull represents a virtual machine operation to push a null value onto the stack.
type OpNull struct {
	*bytecode.OpcodeDetails
}

// NewOpNull creates a new OpNull instance with details mapped from the OpNull opcode.
func NewOpNull() *OpNull {
	return &OpNull{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpNull)}
}

// Execute pushes an undefined value onto the virtual machine's stack.
func (op *OpNull) Execute(v *VM) {
	v.stack.Push(objects.UndefinedValue)
}

// OpBinary represents a type that performs binary operations by extending bytecode.OpcodeDetails.
type OpBinary struct {
	*bytecode.OpcodeDetails
}

// NewOpBinary creates a new instance of OpBinary with its corresponding OpcodeDetails initialized.
func NewOpBinary() *OpBinary {
	return &OpBinary{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpBinary)}
}

// Execute performs a binary operation using operands from the stack, updates the instruction pointer, and handles errors.
func (op *OpBinary) Execute(v *VM) {
	v.ip++
	right := v.stack.Pop()
	left := v.stack.Pop()
	opcode := v.currFrame.Get(v.ip)
	operator := objects.Operator(opcode)
	res, err := left.BinaryOp(operator, right)
	if err != nil {
		v.setError(err)
		return
	}
	v.stack.Push(res)
}

// OpEqual represents an operation that checks if two values are equal and updates the stack accordingly.
type OpEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpEqual creates and returns an instance of OpEqual, initialized with its corresponding opcode details.
func NewOpEqual() *OpEqual {
	return &OpEqual{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpEqual)}
}

// Execute performs the equality comparison between the top two stack values and pushes the result (true or false) back onto the stack.
func (op *OpEqual) Execute(v *VM) {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.FalseValue
	if left.Equals(right) {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

// OpNotEqual is a structure representing the "not equal (!=)" opcode operation in the virtual machine.
// It embeds OpcodeDetails to provide information about the opcode, including its identifier and operands.
type OpNotEqual struct {
	*bytecode.OpcodeDetails
}

// NewOpNotEqual creates and returns a new instance of OpNotEqual with OpcodeDetails initialized from bytecode.
func NewOpNotEqual() *OpNotEqual {
	return &OpNotEqual{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpNotEqual)}
}

// Execute evaluates whether the top two stack elements are unequal and pushes the result as a boolean onto the stack.
func (op *OpNotEqual) Execute(v *VM) {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.TrueValue
	if left.Equals(right) {
		val = objects.FalseValue
	}
	v.stack.Push(val)
}

// OpPop represents an operation that removes the top value from the virtual machine stack.
type OpPop struct {
	*bytecode.OpcodeDetails
}

// NewOpPop creates and returns a new instance of OpPop, initializing it with details corresponding to the OpPop opcode.
func NewOpPop() *OpPop {
	return &OpPop{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpPop)}
}

// Execute performs the operation defined by OpPop, which decreases the stack pointer of the VM.
func (op *OpPop) Execute(v *VM) {
	v.stack.Decrement()
}

// OpTrue represents the opcode for pushing the boolean value true onto the stack.
type OpTrue struct {
	*bytecode.OpcodeDetails
}

// NewOpTrue initializes a new instance of OpTrue, representing the opcode that pushes the boolean value true onto the stack.
func NewOpTrue() *OpTrue {
	return &OpTrue{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpTrue)}
}

// Execute pushes the constant true value onto the virtual machine's stack.
func (op *OpTrue) Execute(v *VM) {
	v.stack.Push(objects.TrueValue)
}

// OpFalse represents an opcode structure for pushing the boolean value false onto the stack.
type OpFalse struct {
	*bytecode.OpcodeDetails
}

// NewOpFalse creates a new instance of OpFalse, representing the operation to push the boolean value false onto the stack.
func NewOpFalse() *OpFalse {
	return &OpFalse{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpFalse)}
}

// Execute pushes a predefined `FalseValue` onto the virtual machine's stack.
func (op *OpFalse) Execute(v *VM) {
	v.stack.Push(objects.FalseValue)
}

// OpLNot represents the logical NOT (!) operation opcode in the virtual machine's instruction set.
type OpLNot struct {
	*bytecode.OpcodeDetails
}

// NewOpLNot creates a new instance of OpLNot, representing a logical NOT operation (!).
func NewOpLNot() *OpLNot {
	return &OpLNot{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpLNot)}
}

// Execute performs a logical NOT operation on the operand at the top of the stack, pushing the result back onto the stack.
func (op *OpLNot) Execute(v *VM) {
	operand := v.stack.Pop()
	val := objects.FalseValue
	if operand.Boolean() {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

// OpBComplement represents an operation for performing a bitwise complement on an operand.
// It extends OpcodeDetails, inheriting its metadata and behaviors.
type OpBComplement struct {
	*bytecode.OpcodeDetails
}

// NewOpBComplement initializes and returns an OpBComplement instance with the corresponding OpcodeDetails configuration.
func NewOpBComplement() *OpBComplement {
	return &OpBComplement{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpBComplement)}
}

// Execute performs the bitwise complement operation on the top stack value. Sets an error if the value is not an integer.
func (op *OpBComplement) Execute(v *VM) {
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := objects.NewInt(^x.Value())
		v.stack.Push(res)
	default:
		v.setError(fmt.Errorf("invalid operation: ^%s", operand.TypeName()))
		return
	}
}

// OpMinus represents an operation for negating a numeric value.
// It embeds OpcodeDetails, providing details such as the opcode, operands, and name.
type OpMinus struct {
	*bytecode.OpcodeDetails
}

// NewOpMinus creates and returns a new OpMinus instance, initializing it with the details of the OpMinus bytecode.
func NewOpMinus() *OpMinus {
	return &OpMinus{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpMinus)}
}

// Execute performs a subtraction operation by negating the top stack element, supporting integers and floats.
// Pushes the result back to the stack or sets an error for unsupported types.
func (op *OpMinus) Execute(v *VM) {
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		res := objects.NewInt(-x.Value())
		v.stack.Push(res)
	case *objects.Float:
		res := objects.NewFloat(-x.Value())
		v.stack.Push(res)
	default:
		v.setError(fmt.Errorf("invalid operation: -%s", operand.TypeName()))
	}
}

// OpJumpFalsy represents an instruction that performs a conditional jump if the stack's top value evaluates to falsy.
type OpJumpFalsy struct {
	*bytecode.OpcodeDetails
}

// NewOpJumpFalsy creates and returns a new instance of OpJumpFalsy initialized with its corresponding OpcodeDetails.
func NewOpJumpFalsy() *OpJumpFalsy {
	return &OpJumpFalsy{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpJumpFalsy)}
}

// Execute advances the instruction pointer, evaluates the stack's top element, and updates the pointer if false.
func (op *OpJumpFalsy) Execute(v *VM) {
	v.ip += 2
	obj := v.stack.Pop()
	if obj.Boolean() {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	}
}

// OpAndJump represents a logical AND operation followed by a conditional jump in the bytecode execution process.
type OpAndJump struct {
	*bytecode.OpcodeDetails
}

// NewOpAndJump creates and returns a new instance of OpAndJump, initializing it with details for the OpAndJump opcode.
func NewOpAndJump() *OpAndJump {
	return &OpAndJump{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpAndJump)}
}

// Execute updates the instruction pointer, evaluates a condition, and adjusts or decrements the stack based on the result.
func (op *OpAndJump) Execute(v *VM) {
	v.ip += 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	} else {
		v.stack.Decrement()
	}
}

// OpOrJump represents an operation that performs a logical OR and jumps based on the result.
type OpOrJump struct {
	*bytecode.OpcodeDetails
}

// NewOpOrJump creates and returns a new instance of OpOrJump, associated with the OpOrJump opcode and its details.
func NewOpOrJump() *OpOrJump {
	return &OpOrJump{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpOrJump)}
}

// Execute advances the instruction pointer, evaluates the stack's top object, and updates the IP based on its boolean value.
func (op *OpOrJump) Execute(v *VM) {
	v.ip += 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		v.stack.Decrement()
	} else {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	}
}

// OpJump represents an unconditional jump operation in the virtual machine, utilizing associated opcode details.
type OpJump struct {
	*bytecode.OpcodeDetails
}

// NewOpJump creates and returns a new instance of OpJump with details initialized for the OpJump opcode.
func NewOpJump() *OpJump {
	return &OpJump{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpJump)}
}

// Execute updates the instruction pointer (`ip`) in the virtual machine (`VM`) to a calculated position in the frame.
func (op *OpJump) Execute(v *VM) {
	pos := v.currFrame.Pos(v.ip+2, v.ip+1)
	v.ip = pos - 1
}

// OpSetGlobal represents a bytecode operation for setting a global variable's value in the virtual machine.
type OpSetGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetGlobal creates and returns a new instance of OpSetGlobal with initialized OpcodeDetails.
func NewOpSetGlobal() *OpSetGlobal {
	return &OpSetGlobal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetGlobal)}
}

// Execute updates the instruction pointer, calculates a global variable position, and sets its value from the stack.
func (op *OpSetGlobal) Execute(v *VM) {
	v.ip += 2
	pos := v.currFrame.Pos(v.ip, v.ip-1)
	val := v.stack.Peek()
	v.globals.Set(uint(pos), val)
}

// OpSetSelGlobal represents an operation for setting a global variable's value using selectors for indexing or access.
type OpSetSelGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelGlobal creates a new instance of OpSetSelGlobal with its corresponding OpcodeDetails initialized.
func NewOpSetSelGlobal() *OpSetSelGlobal {
	return &OpSetSelGlobal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetSelGlobal)}
}

// Execute performs the operation defined by OpSetSelGlobal, updating the VM state and handling global index assignment.
func (op *OpSetSelGlobal) Execute(v *VM) {
	v.ip += 3
	globalIndex := v.currFrame.Pos(v.ip-1, v.ip-2)
	numSelectors := v.currFrame.Get(v.ip)
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(int(numSelectors) + 1)
	glObj := v.globals.Get(uint(globalIndex))
	if err := objects.IndexAssign(glObj, val, selectors); err != nil {
		v.setError(err)
		return
	}
}

// OpGetGlobal represents an operation to retrieve a global variable in the virtual machine.
// It embeds OpcodeDetails for detailed opcode information.
type OpGetGlobal struct {
	*bytecode.OpcodeDetails
}

// NewOpGetGlobal creates a new instance of OpGetGlobal with its associated opcode details.
func NewOpGetGlobal() *OpGetGlobal {
	return &OpGetGlobal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetGlobal)}
}

// Execute retrieves a global object using its index, pushes it onto the stack, and advances the instruction pointer.
func (op *OpGetGlobal) Execute(v *VM) {
	v.ip += 2
	glIndex := v.currFrame.Pos(v.ip, v.ip-1)
	glObj := v.globals.Get(uint(glIndex))
	if glObj == nil {
		//v.setError(fmt.Errorf("undefined global: %d", globalIndex))
		//return
	}
	v.stack.Push(glObj)
}

// OpArray represents a bytecode operation for creating an array object in the virtual machine.
// Extends base OpcodeDetails for opcode, operands, and name information.
type OpArray struct {
	*bytecode.OpcodeDetails
}

// NewOpArray creates and returns a new instance of OpArray, initialized with details for the OpArray operation.
func NewOpArray() *OpArray {
	return &OpArray{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpArray)}
}

// Execute processes the OpArray instruction, constructing an array from stack elements and pushing it onto the stack.
func (op *OpArray) Execute(v *VM) {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	elements := v.stack.PopArrayElements(numElements)
	arr := objects.NewArray(elements)
	v.stack.Push(arr)
}

// OpMap is a wrapper around bytecode.OpcodeDetails, representing a map creation operation in bytecode execution.
type OpMap struct {
	*bytecode.OpcodeDetails
}

// NewOpMap initializes and returns a new instance of OpMap with its OpcodeDetails set to OpMap details.
func NewOpMap() *OpMap {
	return &OpMap{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpMap)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpMap) Execute(v *VM) {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(objects.NewMap(mElem))
}

// OpStruct is a wrapper around bytecode.OpcodeDetails, representing a struct creation operation in bytecode execution.
type OpStruct struct {
	*bytecode.OpcodeDetails
}

// NewOpStruct initializes and returns a new instance of OpStruct with its OpcodeDetails set to OpMap details.
func NewOpStruct() *OpStruct {
	return &OpStruct{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpStruct)}
}

// Execute processes the OpMap instruction, adjusts the instruction pointer, and pushes a new map object onto the stack.
func (op *OpStruct) Execute(v *VM) {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(objects.NewMap(mElem))
}

// OpError represents an operation that creates and assigns an error object in a virtual machine's runtime environment.
type OpError struct {
	*bytecode.OpcodeDetails
}

// NewOpError creates and returns a new instance of OpError with associated OpcodeDetails for the OpError opcode.
func NewOpError() *OpError {
	return &OpError{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpError)}
}

// Execute converts the top value on the VM stack into an error object and replaces it on the stack.
func (op *OpError) Execute(v *VM) {
	value := v.stack.Peek()
	e := objects.NewError(value)
	v.stack.Set(e)
}

// OpImmutable represents an operation that creates immutable objects, inheriting details from OpcodeDetails.
type OpImmutable struct {
	*bytecode.OpcodeDetails
}

// NewOpImmutable creates a new instance of OpImmutable with details loaded from bytecode.OpcodeToDetails.
func NewOpImmutable() *OpImmutable {
	return &OpImmutable{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpImmutable)}
}

// Execute processes the top element on the stack and converts it into an immutable version if it's an array or map.
func (op *OpImmutable) Execute(v *VM) {
	val := v.stack.Peek()
	switch value := val.(type) {
	case *objects.Array:
		obj := objects.NewArrayImmutable(value.Values())
		v.stack.Set(obj)
	case *objects.Map:
		obj := objects.NewMapImmutable(value.Values())
		v.stack.Set(obj)
	}
}

// OpIndex represents the operation for performing an indexing operation on a value.
type OpIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpIndex creates and returns a new instance of OpIndex initialized with its associated OpcodeDetails.
func NewOpIndex() *OpIndex {
	return &OpIndex{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpIndex)}
}

// Execute processes the index operation on the stack, retrieving a value or setting an error if indexing is invalid.
func (op *OpIndex) Execute(v *VM) {
	index := v.stack.Pop()
	left := v.stack.Pop()
	val, err := left.IndexGet(index)
	if err != nil {
		if objects.Is(err, objects.ErrNotIndexable) {
			v.setError(fmt.Errorf("not indexable: %s", index.TypeName()))
			return
		}
		if objects.Is(err, objects.ErrInvalidIndexType) {
			v.setError(fmt.Errorf("invalid index type: %s", index.TypeName()))
			return
		}
		v.setError(err)
		return
	}
	if val == nil {
		val = objects.UndefinedValue
	}
	v.stack.Push(val)
}

// OpSliceIndex represents an operation that performs a slicing action on an array, string, or bytes within a virtual machine.
// It embeds OpcodeDetails to inherit opcode, operand, and name information for execution and identification.
type OpSliceIndex struct {
	*bytecode.OpcodeDetails
}

// NewOpSliceIndex creates a new instance of OpSliceIndex containing details for the slice indexing bytecode operation.
func NewOpSliceIndex() *OpSliceIndex {
	return &OpSliceIndex{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSliceIndex)}
}

// Execute processes the slice operation on the stack, adjusting bounds and supporting various object types like arrays and strings.
func (op *OpSliceIndex) Execute(v *VM) {
	highStack := v.stack.Pop()
	lowStack := v.stack.Pop()
	leftStack := v.stack.Pop()
	lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(leftStack.Length()))
	if err != nil {
		v.setError(err)
		return
	}
	var val objects.IObject = nil
	switch left := leftStack.(type) {
	case *objects.Array:
		val = objects.NewArray(left.Values()[lowIdx:highIdx])
	case *objects.ArrayImmutable:
		val = objects.NewArray(left.Values()[lowIdx:highIdx])
	case *objects.String:
		if val, err = objects.NewString(left.Value()[lowIdx:highIdx]); err != nil {
			v.setError(err)
			return
		}
	case *objects.Bytes:
		val = objects.NewBytes(left.Value()[lowIdx:highIdx])
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
func NewOpCall() *OpCall {
	return &OpCall{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpCall)}
}

// Execute processes the OpCall instruction, invoking the callable or handling array spreads, and manages the stack state.
func (op *OpCall) Execute(v *VM) {
	numArgs := v.currFrame.Get(v.ip + 1)
	v.ip += 2
	value := v.stack.PeekOffset(-1 - numArgs)
	if !value.CanCall() {
		v.setError(fmt.Errorf("%s is not callable: %s", value.String(), value.TypeName()))
		return
	}
	spread := v.currFrame.Get(v.ip + 2)
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
			v.setError(fmt.Errorf("not an array: %s", arrObj.TypeName()))
			return
		}
	}

	if callee, ok := value.(*objects.FunctionCompiled); ok {
		if callee.VarArgs() {
			v.stack.PushVarArgs(numArgs, callee.NumParameters()-1)
			numArgs = callee.NumParameters()
		}
		if numArgs != callee.NumParameters() {
			numParams := callee.NumParameters()
			if callee.VarArgs() {
				numParams--
			}
			v.setError(fmt.Errorf("%s wrong number of arguments: want>=%d, got=%d", callee.Name(), numParams, numArgs))
			return
		}
		// Frame setup
		v.currFrame = v.frames.Get()
		v.frames.Next()
		bp := v.stack.StackPointer() - numArgs
		v.currFrame.Bind(v.ip, callee, bp)
		// Si riserva lo spazio per *tutte* le variabili locali della nuova funzione
		// semplicemente avanzando il puntatore dello stack.
		// Questo garantisce che lo spazio per i calcoli temporanei inizi *dopo*
		// lo spazio riservato per le variabili locali, evitando collisioni.
		v.stack.SetStackPointer(v.stack.StackPointer() + callee.NumLocals())
		v.ip = resetIp
	} else {
		var args []objects.IObject
		args = append(args, v.stack.PeekArrayObject(numArgs)...)
		ret, err := value.Call(args...)
		// Pulisce lo stack dalla funzione e dai suoi argomenti
		v.stack.DecrementCount(numArgs + 1)
		if err != nil {
			if objects.Is(err, objects.ErrWrongNumArguments) {
				v.setError(fmt.Errorf("wrong number of arguments in call to '%s'", value.TypeName()))
			} else {
				v.setError(err)
			}
			return
		}
		if ret != nil && ret != objects.UndefinedValue {
			v.stack.Push(ret)
		}
		if ret == nil {
			ret = objects.UndefinedValue
		}
		v.stack.Push(ret)
	}
}

// OpReturn represents a specialized operation that extends the behavior of bytecode.OpcodeDetails.
type OpReturn struct {
	*bytecode.OpcodeDetails
}

// NewOpReturn creates a new instance of OpReturn with its OpcodeDetails initialized for the OpReturn operation.
func NewOpReturn() *OpReturn {
	return &OpReturn{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpReturn)}
}

// Execute performs the return operation for the current frame, manages the stack, and transitions between frames in the VM.
func (op *OpReturn) Execute(v *VM) {
	v.ip++
	var retVal objects.IObject
	if opcode := v.currFrame.Get(v.ip); opcode == 1 {
		retVal = v.stack.Peek()
	} else {
		retVal = objects.UndefinedValue
	}
	if v.frames.Index() > 1 {
		leavingFrameBasePointer := v.currFrame.BasePointer()
		prevIp := v.currFrame.StartIP()
		v.frames.Previous()
		v.currFrame = v.frames.GetPrev()
		v.ip = prevIp
		v.stack.SetStackPointer(leavingFrameBasePointer)
		v.stack.Push(retVal)
	} else {
		log.Println("returning from the root frame")
		log.Println("stack:")
		v.stack.Print()
		v.Shutdown()
	}
}

// OpDefineLocal represents the opcode for defining a new local variable within the current frame's scope.
type OpDefineLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpDefineLocal creates a new instance of OpDefineLocal with its associated opcode details.
func NewOpDefineLocal() *OpDefineLocal {
	return &OpDefineLocal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpDefineLocal)}
}

// Execute increments the instruction pointer, retrieves a local index, and assigns a stack value to a designated slot.
func (op *OpDefineLocal) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
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
func NewOpSetLocal() *OpSetLocal {
	return &OpSetLocal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetLocal)}
}

// Execute updates a local variable in the current frame using the stack's top value and the local index from instructions.
func (op *OpSetLocal) Execute(v *VM) {
	localIndex := v.currFrame.Get(v.ip + 1)
	v.ip++
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
func NewOpSetSelLocal() *OpSetSelLocal {
	return &OpSetSelLocal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetSelLocal)}
}

// Execute performs the operation of retrieving, modifying, and reassigning a value using selectors in the local scope.
func (op *OpSetSelLocal) Execute(v *VM) {
	localIndex := v.currFrame.Get(v.ip + 1)
	numSelectors := v.currFrame.Get(v.ip + 2)
	v.ip += 2
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
	if err := objects.IndexAssign(dst, val, selectors); err != nil {
		v.setError(err)
		return
	}
}

// OpGetLocal represents an operation to retrieve a local variable from the stack using its index.
type OpGetLocal struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocal creates a new OpGetLocal instance and initializes it with details for the OpGetLocal opcode.
func NewOpGetLocal() *OpGetLocal {
	return &OpGetLocal{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetLocal)}
}

// Execute retrieves a local variable from the current frame's base pointer and pushes it onto the stack.
func (op *OpGetLocal) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	val := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.stack.Push(val)
}

// OpGetBuiltin represents the operation code for retrieving a builtin function in the virtual machine.
type OpGetBuiltin struct {
	*bytecode.OpcodeDetails
}

// NewOpGetBuiltin creates a new instance of OpGetBuiltin with associated opcode details for OpGetBuiltin operations.
func NewOpGetBuiltin() *OpGetBuiltin {
	return &OpGetBuiltin{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetBuiltin)}
}

// Execute advances the instruction pointer, retrieves a builtin symbol, and pushes it onto the VM stack.
func (op *OpGetBuiltin) Execute(v *VM) {
	v.ip++
	builtinIndex := v.currFrame.Get(v.ip)
	symbol := v.loader.GetBuiltinFunction(builtinIndex)
	if symbol == nil {
		v.setError(fmt.Errorf("unkown builtin index: %d", builtinIndex))
		return
	}
	v.stack.Push(symbol)
}

// OpClosure represents a closure operation that creates a new closure in the virtual machine.
type OpClosure struct {
	*bytecode.OpcodeDetails
}

// NewOpClosure returns a new instance of OpClosure initialized with the details of the OpClosure opcode.
func NewOpClosure() *OpClosure {
	return &OpClosure{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpClosure)}
}

// Execute performs the operation associated with the OpClosure opcode, creating a closure and pushing it onto the stack.
func (op *OpClosure) Execute(v *VM) {
	v.ip += 3
	constIndex := v.currFrame.Pos(v.ip-1, v.ip-2)
	numFree := v.currFrame.Get(v.ip)
	glObj := v.globals.Get(uint(constIndex))
	fn, ok := glObj.(*objects.FunctionCompiled)
	if !ok {
		v.setError(fmt.Errorf("not a function: %s", fn.TypeName()))
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
			free[i] = objects.NewObjectPointer(&t)
		}
	}
	v.stack.DecrementCount(numFree)
	cl := objects.NewFunctionCompiled("closure", fn.Instructions().Data(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	v.stack.Push(cl)
}

// OpGetFreePtr represents the opcode for retrieving a free variable pointer in the virtual machine.
// This type embeds OpcodeDetails, which provides opcode metadata such as identifier, operands, and name.
type OpGetFreePtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFreePtr creates a new instance of OpGetFreePtr initialized with the corresponding OpcodeDetails.
func NewOpGetFreePtr() *OpGetFreePtr {
	return &OpGetFreePtr{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetFreePtr)}
}

// Execute executes the OpGetFreePtr operation, pushing a free variable onto the stack based on the current instruction pointer.
func (op *OpGetFreePtr) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := v.currFrame.FreeVarsIndex(freeIndex)
	v.stack.Push(val)
}

// OpGetFree represents an operation to retrieve a free variable in a closure during execution.
type OpGetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpGetFree creates and returns a new instance of OpGetFree, initializing its OpcodeDetails using bytecode metadata.
func NewOpGetFree() *OpGetFree {
	return &OpGetFree{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetFree)}
}

// Execute increments the instruction pointer, retrieves a value using free variable index, and pushes it onto the stack.
func (op *OpGetFree) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := *v.currFrame.FreeVarsIndex(freeIndex).Value()
	v.stack.Push(val)
}

// OpSetFree represents an operation to set the value of a free variable within a closure's environment.
type OpSetFree struct {
	*bytecode.OpcodeDetails
}

// NewOpSetFree creates and returns a new instance of OpSetFree initialized with its corresponding OpcodeDetails.
func NewOpSetFree() *OpSetFree {
	return &OpSetFree{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetFree)}
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpSetFree) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	o := v.stack.Pop()
	v.currFrame.FreeVarsIndex(freeIndex).SetValue(o)
}

// OpGetLocalPtr retrieves a local variable as a pointer using its index within the current frame.
type OpGetLocalPtr struct {
	*bytecode.OpcodeDetails
}

// NewOpGetLocalPtr creates and returns a new instance of OpGetLocalPtr, initializing its OpcodeDetails.
func NewOpGetLocalPtr() *OpGetLocalPtr {
	return &OpGetLocalPtr{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpGetLocalPtr)}
}

// Execute advances the instruction pointer, retrieves a local variable, and pushes an ObjectPointer to the stack.
func (op *OpGetLocalPtr) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	sp := v.currFrame.BasePointer() + localIndex
	val := v.stack.PeekAbsolute(sp)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		v.stack.Push(obj)
		return
	}
	freeVar := objects.NewObjectPointer(&val)
	v.stack.SetAbsolute(sp, freeVar)
	v.stack.Push(freeVar)
}

// OpSetSelFree represents an operation to set a free variable's value using selectors.
type OpSetSelFree struct {
	*bytecode.OpcodeDetails
}

// NewOpSetSelFree creates a new instance of OpSetSelFree with initialized OpcodeDetails referencing OpSetSelFree.
func NewOpSetSelFree() *OpSetSelFree {
	return &OpSetSelFree{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSetSelFree)}
}

// Execute updates the instruction pointer, retrieves operands, processes selectors, and performs indexed assignment in the VM.
func (op *OpSetSelFree) Execute(v *VM) {
	v.ip += 2
	freeIndex := v.currFrame.Get(v.ip - 1)
	numSelectors := v.currFrame.Get(v.ip)

	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(numSelectors + 1)
	fvi := v.currFrame.FreeVarsIndex(freeIndex)
	if err := objects.IndexAssign(*fvi.Value(), val, selectors); err != nil {
		v.setError(err)
		return
	}
}

// OpIteratorInit represents an operation that initializes an iterator over an iterable object.
// It embeds OpcodeDetails for additional opcode-specific metadata.
type OpIteratorInit struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorInit creates and returns a new instance of OpIteratorInit with associated opcode details.
func NewOpIteratorInit() *OpIteratorInit {
	return &OpIteratorInit{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpIteratorInit)}
}

// Execute initializes an iterator for an iterable object and stores it in the specified local slot in the current frame.
func (op *OpIteratorInit) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	iterable := v.stack.Pop()
	if !iterable.CanIterate() {
		v.setError(fmt.Errorf("not iterable: %s", iterable.TypeName()))
		return
	}
	iterator := iterable.Iterate()
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, iterator)
}

// OpIteratorNext represents an operation code for advancing an iterator to the next element in the virtual machine.
type OpIteratorNext struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorNext creates a new instance of OpIteratorNext with associated opcode details.
func NewOpIteratorNext() *OpIteratorNext {
	return &OpIteratorNext{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpIteratorNext)}
}

// Execute processes the next iterator state in the current frame, pushing a boolean to the stack indicating iteration status.
func (op *OpIteratorNext) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.setError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	if iterator.Next() {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// OpIteratorKey wraps bytecode.OpcodeDetails to represent the iterator key retrieval operation in a virtual machine.
type OpIteratorKey struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorKey creates a new instance of OpIteratorKey with associated opcode details.
func NewOpIteratorKey() *OpIteratorKey {
	return &OpIteratorKey{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpIteratorKey)}
}

// Execute processes the "iterator key" operation, retrieves the iterator key, and pushes it onto the VM stack.
func (op *OpIteratorKey) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.setError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.stack.Push(iterator.Key())
}

// OpIteratorValue retrieves the value from the current iterator position.
// It embeds OpcodeDetails, providing access to the opcode's metadata and operations.
type OpIteratorValue struct {
	*bytecode.OpcodeDetails
}

// NewOpIteratorValue creates and returns a new instance of OpIteratorValue with its associated OpcodeDetails initialized.
func NewOpIteratorValue() *OpIteratorValue {
	return &OpIteratorValue{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpIteratorValue)}
}

// Execute processes the next instruction to retrieve and push the current value of an iterator onto the stack.
func (op *OpIteratorValue) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	iteratorObj := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	iterator, ok := iteratorObj.(objects.IIterator)
	if !ok {
		v.setError(fmt.Errorf("not an iterator: %s", iteratorObj.TypeName()))
		return
	}
	v.stack.Push(iterator.Value())
}

// OpReferences extends OpcodeDetails to represent operations specifically related to reference handling in the bytecode.
type OpReferences struct {
	*bytecode.OpcodeDetails
}

// NewOpReferences initializes a new OpReferences instance with corresponding OpcodeDetails from the bytecode package.
func NewOpReferences() *OpReferences {
	return &OpReferences{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpReferences)}
}

// Execute processes the specified VM instruction, adjusts the instruction pointer, and pushes a reference onto the stack.
func (op *OpReferences) Execute(v *VM) {
	v.ip += 2
	nameIndex := v.currFrame.Pos(v.ip, v.ip-1)
	if nameIndex < 0 || nameIndex >= len(v.references) {
		v.setError(fmt.Errorf("invalid attribute index %d", nameIndex))
		return
	}
	symbol := v.references[nameIndex]
	v.stack.Push(symbol)
}

// OpSuspend represents an operation that suspends the execution of the virtual machine.
type OpSuspend struct {
	*bytecode.OpcodeDetails
}

// NewOpSuspend creates and returns a new OpSuspend instance with opcode details initialized for the suspend operation.
func NewOpSuspend() *OpSuspend {
	return &OpSuspend{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpSuspend)}
}

// Execute performs the suspend operation on the given virtual machine by setting its shutdown state to true.
func (op *OpSuspend) Execute(v *VM) {
	v.shutdown = true
}

// OpUnknown represents an unknown or unsupported operation in the bytecode execution context.
type OpUnknown struct {
	*bytecode.OpcodeDetails
}

// NewOpUnknown creates a new instance of OpUnknown with its corresponding OpcodeDetails configuration set.
func NewOpUnknown() *OpUnknown {
	return &OpUnknown{OpcodeDetails: bytecode.OpcodeToDetails(bytecode.OpUnknown)}
}

// Execute handles the execution of an unknown opcode, sets an error state, and stops the virtual machine.
func (op *OpUnknown) Execute(v *VM) {
	pos := v.currFrame.Get(v.ip)
	v.setError(fmt.Errorf("unknown opcode: %d", pos))
}
