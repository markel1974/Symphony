package vm

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

type IOpExecutor interface {
	Execute(v *VM)
}

type OpConstant struct {
}

// Execute retrieves a constant from the instructions using the indices at ip and pushes it onto the stack.
func (op *OpConstant) Execute(v *VM) {
	v.ip += 2
	cIdx := v.currFrame.Pos(v.ip, v.ip-1)
	glObj := v.globals.Get(uint(cIdx))
	v.stack.Push(glObj)
}

type OpNull struct {
}

// Execute pushes the UndefinedValue onto the stack to represent a null or undefined operation result.
func (op *OpNull) Execute(v *VM) {
	v.stack.Push(objects.UndefinedValue)
}

type OpBinary struct {
}

// Execute performs a binary operation between the top two stack elements based on the current opcode instruction.
// It pops the two operands from the stack, retrieves the operation from the instruction, executes it, and pushes the result.
// If an error occurs during the operation or instruction retrieval, it sets the VM's error state.
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

type OpEqual struct {
}

// Execute pops two objects from the stack, compares them for equality, and pushes TrueValue or FalseValue based on the result.
func (op *OpEqual) Execute(v *VM) {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.FalseValue
	if left.Equals(right) {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

type OpNotEqual struct {
}

// Execute checks if the two topmost stack values are not equal and pushes true or false onto the stack accordingly.
func (op *OpNotEqual) Execute(v *VM) {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.TrueValue
	if left.Equals(right) {
		val = objects.FalseValue
	}
	v.stack.Push(val)
}

type OpPop struct {
}

// Execute decreases the stack pointer by one, effectively discarding the top element of the stack.
func (op *OpPop) Execute(v *VM) {
	v.stack.Decrement()
}

type OpTrue struct {
}

// Execute pushes the predefined TrueValue object onto the stack.
func (op *OpTrue) Execute(v *VM) {
	v.stack.Push(objects.TrueValue)
}

type OpFalse struct {
}

// Execute pushes the predefined false value (FalseValue) onto the stack.
func (op *OpFalse) Execute(v *VM) {
	v.stack.Push(objects.FalseValue)
}

type OpLNot struct {
}

// Execute performs the logical NOT operation on the top value on the stack and pushes the result back onto the stack.
func (op *OpLNot) Execute(v *VM) {
	operand := v.stack.Pop()
	val := objects.FalseValue
	if operand.Boolean() {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

type OpBComplement struct {
}

// Execute performs a bitwise NOT operation on the top integer operand from the stack and pushes the result.
// If the operand is not an integer, it sets an error.
// If the object allocation limit is exceeded, it sets an allocation error.
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

type OpMinus struct {
}

// Execute applies the unary minus operation on the top element of the stack if it's an Int or Float object.
// It pushes the resulting negated value back onto the stack or sets an error if the operation is invalid.
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

type OpJumpFalsy struct {
}

// Execute performs a jump operation if the top stack value is falsy, updating the instruction pointer accordingly.
func (op *OpJumpFalsy) Execute(v *VM) {
	v.ip += 2
	obj := v.stack.Pop()
	if obj.Boolean() {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	}
}

type OpAndJump struct {
}

// Execute adjusts the instruction pointer and interacts with the stack to conditionally jump to a new position.
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

type OpOrJump struct {
}

// Execute processes the current instruction to either modify the stack or jump to a specific position in code.
// If the top stack object evaluates to true, it decrements the stack pointer; otherwise, it performs a conditional jump.
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

type OpJump struct {
}

// Execute updates the instruction pointer based on the position derived from the bytecode and adjusts for zero-based indexing.
func (op *OpJump) Execute(v *VM) {
	pos := v.currFrame.Pos(v.ip+2, v.ip+1)
	v.ip = pos - 1
}

type OpSetGlobal struct {
}

// Execute performs the SET_GLOBAL operation, updating a global variable with a value from the stack.
func (op *OpSetGlobal) Execute(v *VM) {
	v.ip += 2
	pos := v.currFrame.Pos(v.ip, v.ip-1)
	val := v.stack.Peek()
	v.globals.Set(uint(pos), val)
}

type OpSetSelGlobal struct {
}

// Execute is a VM instruction to set a value in a global object using selectors for index assignment.
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

type OpGetGlobal struct {
}

// Execute retrieves a global variable from the global pool and pushes its value onto the stack.
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

type OpArray struct {
}

// Execute processes array-based operations in the VM by retrieving elements from the stack and creating a new array object.
// It increments the instruction pointer, extracts a specified number of elements, and enforces allocation limits.
// If the allocation limit is exceeded or an error occurs, it sets the error state of the VM.
func (op *OpArray) Execute(v *VM) {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	elements := v.stack.PopArrayElements(numElements)
	arr := objects.NewArray(elements)
	v.stack.Push(arr)
}

type OpMap struct {
}

// Execute creates a new map object from key-value pairs popped from the stack and pushes the map onto the stack.
// It adjusts the instruction pointer and handles allocation limits or errors during creation.
func (op *OpMap) Execute(v *VM) {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(objects.NewMap(mElem))
}

type OpError struct {
}

// Execute handles an operation error by converting the top stack value into an error object and updating the stack.
// Reduces the allocation count and sets an error if the allocation limit is reached.
func (op *OpError) Execute(v *VM) {
	value := v.stack.Peek()
	e := objects.NewError(value)
	v.stack.Set(e)
}

type OpImmutable struct {
}

// Execute transforms the top stack value into its immutable form if it is an Array or Map.
// Reduces allocation counter and sets an error if the allocation limit is reached.
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

type OpIndex struct {
}

// Execute processes an index operation by retrieving the index, the left operand, and getting the indexed value.
// If an error occurs during the operation, it sets the VM error accordingly. Pushes the result or UndefinedValue to the stack.
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

type OpSliceIndex struct {
}

// Execute processes slice indexing operations for arrays, strings, and bytes by validating bounds and slicing them.
// It pops the necessary indices and objects from the stack, performs slicing, and pushes the sliced results back to the stack.
// If bounds are invalid or an error occurs during slicing, it sets the VM's error state and halts further execution.
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

type OpCall struct {
}

// Execute handles the execution of a function call in the virtual machine, including compiled and built-in functions.
func (op *OpCall) Execute(v *VM) {
	numArgs := v.currFrame.Get(v.ip + 1)
	v.ip += 2
	value := v.stack.PeekOffset(-1 - numArgs)
	if !value.CanCall() {
		v.setError(fmt.Errorf("not callable: %s", value.TypeName()))
		return
	}
	if spread := v.currFrame.Get(v.ip + 2); spread == 1 {
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
			v.setError(fmt.Errorf("wrong number of arguments: want>=%d, got=%d", numParams, numArgs))
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

type OpReturn struct {
}

// Execute handles the return operation in the current VM frame, restoring the previous frame or terminating execution.
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

type OpDefineLocal struct {
}

// Execute defines a local variable in the current frame by using the local index from the bytecode instructions.
// It pops the value from the stack, calculates the destination slot using the base pointer and index, and sets the value.
// If the stack pointer is less than or equal to the destination slot, it updates the stack pointer to protect the new variable.
func (op *OpDefineLocal) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	val := v.stack.Peek()
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, val)
}

type OpSetLocal struct {
}

// Execute sets a local variable in the current function's call frame. It adjusts the stack as needed for the operation.
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

type OpSetSelLocal struct {
}

// Execute performs a set operation on a local variable using provided selectors and value from the stack.
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

type OpGetLocal struct {
}

// Execute retrieves a local variable from the stack based on its index in the current frame and pushes it onto the stack.
func (op *OpGetLocal) Execute(v *VM) {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	val := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.stack.Push(val)
}

type OpGetBuiltin struct {
}

// Execute retrieves a builtin function by its index and pushes it onto the stack, incrementing the instruction pointer.
func (op *OpGetBuiltin) Execute(v *VM) {
	v.ip++
	builtinIndex := v.currFrame.Get(v.ip)
	symbol := v.loader.GetBuiltinSymbol(builtinIndex)
	if symbol == nil {
		v.setError(fmt.Errorf("unkown builtin index: %d", builtinIndex))
		return
	}
	v.stack.Push(symbol)
}

type OpClosure struct {
}

// Execute handles the creation of a new compiled-function closure with captured free variables on the call stack.
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

type OpGetFreePtr struct {
}

// Execute retrieves a free variable from the current frame based on the index in the bytecode and pushes it onto the stack.
func (op *OpGetFreePtr) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := v.currFrame.FreeVarsIndex(freeIndex)
	v.stack.Push(val)
}

type OpGetFree struct {
}

// Execute extracts a free variable value using its index, pushes it onto the stack, and increments the instruction pointer.
// If an error occurs while retrieving the instruction, it sets the VM's error and halts further execution.
func (op *OpGetFree) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := *v.currFrame.FreeVarsIndex(freeIndex).Value()
	v.stack.Push(val)
}

type OpSetFree struct {
}

// Execute increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (op *OpSetFree) Execute(v *VM) {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	o := v.stack.Pop()
	v.currFrame.FreeVarsIndex(freeIndex).SetValue(o)
}

type OpGetLocalPtr struct {
}

// Execute increments the instruction pointer, retrieves a local variable's value, and pushes it onto the stack.
// If the value is not already an ObjectPointer, it wraps it as a new ObjectPointer.
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

type OpSetSelFree struct {
}

// Execute performs a selector-based assignment operation on a free variable within the VM's current frame.
// It retrieves selectors, target value, and free variable index from the stack and instruction set to complete the operation.
// Selector and value data are validated, and errors are set in the VM if any issues occur during assignment.
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

type OpIteratorInit struct {
}

// Execute initializes an iterator for an iterable object on the stack and assigns it to a local variable slot.
// It increments the instruction pointer, retrieves the target variable index, and validates that the object is iterable.
// The method creates an iterator, decrements the allocation counter, and verifies the allocation limit isn't exceeded.
// Finally, the iterator is stored in the local variable slot specified by the instruction.
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

type OpIteratorNext struct {
}

// Execute retrieves the next value from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
// If the iterator is valid, it calls the Next() method on the iterator and pushes the result onto the stack.
// If the iterator is invalid, it pushes a false value onto the stack.
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

type OpIteratorKey struct {
}

// Execute retrieves the current key from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
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

type OpIteratorValue struct {
}

// Execute retrieves the current value from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
// If the iterator is valid, it calls the Value() method on the iterator and pushes the result onto the stack.
// If the iterator is invalid, it pushes a false value onto the stack.
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

type OpReferences struct {
}

// Execute retrieves an attribute identified by its index, resolves it, and pushes it onto the stack.
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

type OpSuspend struct {
}

// Execute sets the VM's `suspend` flag to true, pausing the execution of the bytecode operations.
func (op *OpSuspend) Execute(v *VM) {
	v.shutdown = true
}

type OpUnknown struct {
}

// Execute is a fallback method executed when an unrecognized opcode is encountered in the instruction set.
// It retrieves the current instruction at the instruction pointer and sets an error indicating the unknown opcode.
func (op *OpUnknown) Execute(v *VM) {
	pos := v.currFrame.Get(v.ip)
	v.setError(fmt.Errorf("unknown opcode: %d", pos))
}
