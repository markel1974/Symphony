package vm

import (
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// resetIp is the instruction pointer value used to reset the VM's instruction pointer to the beginning of the main function.'
const (
	resetIp = -1
)

// globalsSize defines the maximum size of the global variables space.
// stackSize specifies the size limit of the stack for function execution.
// maxFrames indicates the maximum number of call frames allowed.
const (
	stackSize = 2048
	maxFrames = 1024
)

// sequenceLen defines the length of a sequence using a bitwise left shift operation.
// sequenceMask is a bitmask derived from sequenceLen to restrict sequence values within the specified range.
const (
	sequenceLen  = 1 << 8
	sequenceMask = sequenceLen - 1
)

// ISequencer defines an interface to generate a sequence of functions for a given Virtual Machine instance.
type ISequencer interface {
	Create(vm *VM) []func()
}

// VM represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type VM struct {
	sourceFiles *bytecode.Files
	stack       *Stack
	frames      *Frames
	currFrame   *Frame
	ip          int
	shutdown    bool
	err         error
	sequencer   []func()
	references  []objects.IObject
	loader      ILoader
	globals     *Globals
	functions   map[string]*objects.FunctionCompiled
}

// NewVM initializes and returns a new virtual machine instance configured with the provided components and settings.
func NewVM(loader ILoader, sequencer ISequencer, bc *bytecode.Bytecode, maxAllocations int64) (*VM, error) {
	if bc == nil {
		return nil, fmt.Errorf("bytecode is nil")
	}
	if maxAllocations < 1 {
		return nil, fmt.Errorf("max allocations must be greater than 0")
	}
	references, err := loader.ResolveSymbols(bc.References())
	if err != nil {
		return nil, err
	}
	if sequencer == nil {
		sequencer = NewSequencer()
	}
	functions := make(map[string]*objects.FunctionCompiled)
	globals := make([]objects.IObject, len(bc.Constants()))
	for idx, constant := range bc.Constants() {
		switch v := constant.(type) {
		case *objects.FunctionCompiled:
			functions[v.Name()] = v
		}
		globals[idx] = constant
	}
	v := &VM{
		sourceFiles: bc.SourceFiles(),
		ip:          resetIp,
		loader:      loader,
		references:  references,
		functions:   functions,
	}
	v.stack = NewStack(stackSize, maxAllocations, v.setError)
	v.frames = NewFrames(maxFrames, v.setError)
	v.globals = NewGlobals(globals, v.setError)
	v.sequencer = sequencer.Create(v)
	v.Reset()
	return v, nil
}

// Shutdown gracefully shuts down the virtual machine by setting its internal state to signify termination.
func (v *VM) Shutdown() {
	v.shutdown = true
}

// Reset reinitializes the virtual machine's state, clears the stack and frames, and resets execution-related variables.
func (v *VM) Reset() {
	v.ip = resetIp
	v.stack.Reset()
	v.frames.Reset()
	v.err = nil
	v.shutdown = false
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) Run(main string) error {
	mainFn, _ := v.functions[main]
	if mainFn == nil {
		return fmt.Errorf("main function not found")
	}
	v.currFrame = v.frames.Head()
	v.currFrame.Bind(v.ip, mainFn, 0)
	v.stack.SetStackPointer(v.currFrame.NumLocals())

	v.run()

	if v.err != nil {
		filePos, _ := v.sourceFiles.Position(v.currFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.sourceFiles.Position(frame.SourcePos(frame.StartIP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
}

// run executes the main instruction loop for the virtual machine, updating the instruction pointer and processing opcodes.
func (v *VM) run() {
	for {
		v.ip++
		opcode := v.currFrame.Get(v.ip)
		opcode = opcode & sequenceMask
		log.Println("Executing instruction ", opcode, bytecode.OpcodeNames[opcode])
		v.sequencer[opcode]()
		if v.shutdown {
			break
		}
	}
}

// setError sets the internal error state of the VM and marks it for shutdown.
func (v *VM) setError(err error) {
	v.err = err
	v.shutdown = true
}

// checkBounds validates and adjusts slice bounds using provided low and high indices, ensuring they are within valid range.
func (v *VM) checkBounds(lowStack objects.IObject, highStack objects.IObject, numElements int64) (int64, int64, error) {
	var lowIdx int64
	if lowStack != objects.UndefinedValue {
		if low, ok := lowStack.(*objects.Int); ok {
			lowIdx = low.Value()
		} else {
			return 0, 0, fmt.Errorf("invalid slice index type: %s", low.TypeName())
		}
	}
	var highIdx int64
	if highStack == objects.UndefinedValue {
		highIdx = numElements
	} else if high, ok := highStack.(*objects.Int); ok {
		highIdx = high.Value()
	} else {
		return 0, 0, fmt.Errorf("invalid slice index type: %s", high.TypeName())
	}
	if lowIdx > highIdx {
		return 0, 0, fmt.Errorf("invalid slice index: %d > %d", lowIdx, highIdx)
	}
	if lowIdx < 0 {
		lowIdx = 0
	} else if lowIdx > numElements {
		lowIdx = numElements
	}
	if highIdx < 0 {
		highIdx = 0
	} else if highIdx > numElements {
		highIdx = numElements
	}
	return lowIdx, highIdx, nil
}

// doOpConstant retrieves a constant from the instructions using the indices at ip and pushes it onto the stack.
func (v *VM) doOpConstant() {
	v.ip += 2
	cIdx := v.currFrame.Pos(v.ip, v.ip-1)
	glObj := v.globals.Get(uint(cIdx))
	v.stack.Push(glObj)
}

// doOpNull pushes the UndefinedValue onto the stack to represent a null or undefined operation result.
func (v *VM) doOpNull() {
	v.stack.Push(objects.UndefinedValue)
}

// doOpBinary performs a binary operation between the top two stack elements based on the current opcode instruction.
// It pops the two operands from the stack, retrieves the operation from the instruction, executes it, and pushes the result.
// If an error occurs during the operation or instruction retrieval, it sets the VM's error state.
func (v *VM) doOpBinary() {
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

// doOpEqual pops two objects from the stack, compares them for equality, and pushes TrueValue or FalseValue based on the result.
func (v *VM) doOpEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.FalseValue
	if left.Equals(right) {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

// doOpNotEqual checks if the two topmost stack values are not equal and pushes true or false onto the stack accordingly.
func (v *VM) doOpNotEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	val := objects.TrueValue
	if left.Equals(right) {
		val = objects.FalseValue
	}
	v.stack.Push(val)
}

// doOpPop decreases the stack pointer by one, effectively discarding the top element of the stack.
func (v *VM) doOpPop() {
	v.stack.Decrement()
}

// doOpTrue pushes the predefined TrueValue object onto the stack.
func (v *VM) doOpTrue() {
	v.stack.Push(objects.TrueValue)
}

// doOpFalse pushes the predefined false value (FalseValue) onto the stack.
func (v *VM) doOpFalse() {
	v.stack.Push(objects.FalseValue)
}

// doOpLNot performs the logical NOT operation on the top value on the stack and pushes the result back onto the stack.
func (v *VM) doOpLNot() {
	operand := v.stack.Pop()
	val := objects.FalseValue
	if operand.Boolean() {
		val = objects.TrueValue
	}
	v.stack.Push(val)
}

// doOpBComplement performs a bitwise NOT operation on the top integer operand from the stack and pushes the result.
// If the operand is not an integer, it sets an error.
// If the object allocation limit is exceeded, it sets an allocation error.
func (v *VM) doOpBComplement() {
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

// doOpMinus applies the unary minus operation on the top element of the stack if it's an Int or Float object.
// It pushes the resulting negated value back onto the stack or sets an error if the operation is invalid.
func (v *VM) doOpMinus() {
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

// doOpJumpFalsy performs a jump operation if the top stack value is falsy, updating the instruction pointer accordingly.
func (v *VM) doOpJumpFalsy() {
	v.ip += 2
	obj := v.stack.Pop()
	if obj.Boolean() {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	}
}

// doOpAndJump adjusts the instruction pointer and interacts with the stack to conditionally jump to a new position.
func (v *VM) doOpAndJump() {
	v.ip += 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	} else {
		v.stack.Decrement()
	}
}

// doOpOrJump processes the current instruction to either modify the stack or jump to a specific position in code.
// If the top stack object evaluates to true, it decrements the stack pointer; otherwise, it performs a conditional jump.
func (v *VM) doOpOrJump() {
	v.ip += 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		v.stack.Decrement()
	} else {
		pos := v.currFrame.Pos(v.ip, v.ip-1)
		v.ip = pos - 1
	}
}

// doOpJump updates the instruction pointer based on the position derived from the bytecode and adjusts for zero-based indexing.
func (v *VM) doOpJump() {
	pos := v.currFrame.Pos(v.ip+2, v.ip+1)
	v.ip = pos - 1
}

// doOpSetGlobal performs the SET_GLOBAL operation, updating a global variable with a value from the stack.
func (v *VM) doOpSetGlobal() {
	v.ip += 2
	pos := v.currFrame.Pos(v.ip, v.ip-1)
	val := v.stack.Peek()
	v.globals.Set(uint(pos), val)
}

// doOpSetSelGlobal is a VM instruction to set a value in a global object using selectors for index assignment.
func (v *VM) doOpSetSelGlobal() {
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

// doOpGetGlobal retrieves a global variable from the global pool and pushes its value onto the stack.
func (v *VM) doOpGetGlobal() {
	v.ip += 2
	glIndex := v.currFrame.Pos(v.ip, v.ip-1)
	glObj := v.globals.Get(uint(glIndex))
	if glObj == nil {
		//v.setError(fmt.Errorf("undefined global: %d", globalIndex))
		//return
	}
	v.stack.Push(glObj)
}

// doOpArray processes array-based operations in the VM by retrieving elements from the stack and creating a new array object.
// It increments the instruction pointer, extracts a specified number of elements, and enforces allocation limits.
// If the allocation limit is exceeded or an error occurs, it sets the error state of the VM.
func (v *VM) doOpArray() {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	elements := v.stack.PopArrayElements(numElements)
	arr := objects.NewArray(elements)
	v.stack.Push(arr)
}

// doOpMap creates a new map object from key-value pairs popped from the stack and pushes the map onto the stack.
// It adjusts the instruction pointer and handles allocation limits or errors during creation.
func (v *VM) doOpMap() {
	v.ip += 2
	numElements := v.currFrame.Pos(v.ip, v.ip-1)
	mElem := v.stack.PopMapElements(numElements)
	v.stack.Push(objects.NewMap(mElem))
}

// doOpError handles an operation error by converting the top stack value into an error object and updating the stack.
// Reduces the allocation count and sets an error if the allocation limit is reached.
func (v *VM) doOpError() {
	value := v.stack.Peek()
	e := objects.NewError(value)
	v.stack.Set(e)
}

// doOpImmutable transforms the top stack value into its immutable form if it is an Array or Map.
// Reduces allocation counter and sets an error if the allocation limit is reached.
func (v *VM) doOpImmutable() {
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

// doOpIndex processes an index operation by retrieving the index, the left operand, and getting the indexed value.
// If an error occurs during the operation, it sets the VM error accordingly. Pushes the result or UndefinedValue to the stack.
func (v *VM) doOpIndex() {
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

// doOpSliceIndex processes slice indexing operations for arrays, strings, and bytes by validating bounds and slicing them.
// It pops the necessary indices and objects from the stack, performs slicing, and pushes the sliced results back to the stack.
// If bounds are invalid or an error occurs during slicing, it sets the VM's error state and halts further execution.
func (v *VM) doOpSliceIndex() {
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

// doOpCall handles the execution of a function call in the virtual machine, including compiled and built-in functions.
func (v *VM) doOpCall() {
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

// doOpReturn handles the return operation in the current VM frame, restoring the previous frame or terminating execution.
func (v *VM) doOpReturn() {
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

// doOpDefineLocal defines a local variable in the current frame by using the local index from the bytecode instructions.
// It pops the value from the stack, calculates the destination slot using the base pointer and index, and sets the value.
// If the stack pointer is less than or equal to the destination slot, it updates the stack pointer to protect the new variable.
func (v *VM) doOpDefineLocal() {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	val := v.stack.Peek()
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, val)
}

// doOpSetLocal sets a local variable in the current function's call frame. It adjusts the stack as needed for the operation.
func (v *VM) doOpSetLocal() {
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

// doOpSetSelLocal performs a set operation on a local variable using provided selectors and value from the stack.
func (v *VM) doOpSetSelLocal() {
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

// doOpGetLocal retrieves a local variable from the stack based on its index in the current frame and pushes it onto the stack.
func (v *VM) doOpGetLocal() {
	v.ip++
	localIndex := v.currFrame.Get(v.ip)
	val := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.stack.Push(val)
}

// doOpGetBuiltin retrieves a builtin function by its index and pushes it onto the stack, incrementing the instruction pointer.
func (v *VM) doOpGetBuiltin() {
	v.ip++
	builtinIndex := v.currFrame.Get(v.ip)
	symbol := v.loader.GetBuiltinSymbol(builtinIndex)
	if symbol == nil {
		v.setError(fmt.Errorf("unkown builtin index: %d", builtinIndex))
		return
	}
	v.stack.Push(symbol)
}

// doOpClosure handles the creation of a new compiled-function closure with captured free variables on the call stack.
func (v *VM) doOpClosure() {
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

// doOpGetFreePtr retrieves a free variable from the current frame based on the index in the bytecode and pushes it onto the stack.
func (v *VM) doOpGetFreePtr() {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := v.currFrame.FreeVarsIndex(freeIndex)
	v.stack.Push(val)
}

// doOpGetFree extracts a free variable value using its index, pushes it onto the stack, and increments the instruction pointer.
// If an error occurs while retrieving the instruction, it sets the VM's error and halts further execution.
func (v *VM) doOpGetFree() {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	val := *v.currFrame.FreeVarsIndex(freeIndex).Value()
	v.stack.Push(val)
}

// doOpSetFree increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (v *VM) doOpSetFree() {
	v.ip++
	freeIndex := v.currFrame.Get(v.ip)
	o := v.stack.Pop()
	v.currFrame.FreeVarsIndex(freeIndex).SetValue(o)
}

// doOpGetLocalPtr increments the instruction pointer, retrieves a local variable's value, and pushes it onto the stack.
// If the value is not already an ObjectPointer, it wraps it as a new ObjectPointer.
func (v *VM) doOpGetLocalPtr() {
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

// doOpSetSelFree performs a selector-based assignment operation on a free variable within the VM's current frame.
// It retrieves selectors, target value, and free variable index from the stack and instruction set to complete the operation.
// Selector and value data are validated, and errors are set in the VM if any issues occur during assignment.
func (v *VM) doOpSetSelFree() {
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

// doOpIteratorInit initializes an iterator for an iterable object on the stack and assigns it to a local variable slot.
// It increments the instruction pointer, retrieves the target variable index, and validates that the object is iterable.
// The method creates an iterator, decrements the allocation counter, and verifies the allocation limit isn't exceeded.
// Finally, the iterator is stored in the local variable slot specified by the instruction.
func (v *VM) doOpIteratorInit() {
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

// doOpIteratorNext retrieves the next value from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
// If the iterator is valid, it calls the Next() method on the iterator and pushes the result onto the stack.
// If the iterator is invalid, it pushes a false value onto the stack.
func (v *VM) doOpIteratorNext() {
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

// doOpIteratorKey retrieves the current key from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
func (v *VM) doOpIteratorKey() {
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

// doOpIteratorValue retrieves the current value from an iterator and pushes it onto the stack.
// It increments the instruction pointer, retrieves the iterator index, and validates that the iterator is valid.
// If the iterator is valid, it calls the Value() method on the iterator and pushes the result onto the stack.
// If the iterator is invalid, it pushes a false value onto the stack.
func (v *VM) doOpIteratorValue() {
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

// doOpReferences retrieves an attribute identified by its index, resolves it, and pushes it onto the stack.
func (v *VM) doOpReferences() {
	v.ip += 2
	nameIndex := v.currFrame.Pos(v.ip, v.ip-1)
	if nameIndex < 0 || nameIndex >= len(v.references) {
		v.setError(fmt.Errorf("invalid attribute index %d", nameIndex))
		return
	}
	symbol := v.references[nameIndex]
	v.stack.Push(symbol)
}

// doOpSuspend sets the VM's `suspend` flag to true, pausing the execution of the bytecode operations.
func (v *VM) doOpSuspend() {
	v.shutdown = true
}

// doOpUnknown is a fallback method executed when an unrecognized opcode is encountered in the instruction set.
// It retrieves the current instruction at the instruction pointer and sets an error indicating the unknown opcode.
func (v *VM) doOpUnknown() {
	pos := v.currFrame.Get(v.ip)
	v.setError(fmt.Errorf("unknown opcode: %d", pos))
}
