package vm

import (
	"fmt"

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
	globalsSize = 1024
	stackSize   = 2048
	maxFrames   = 1024
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
	sourceFiles      *bytecode.Files
	constants        []objects.IObject
	globals          []objects.IObject
	stack            *Stack
	frames           *Frames
	currFrame        *Frame
	ip               int
	currInstructions *objects.Instructions
	abort            bool
	suspend          bool
	maxAllocations   int64
	allocations      int64
	err              error
	sequencer        []func()
	references       []objects.IObject
	builtin          []*objects.FunctionBuiltin
}

// NewVM initializes and returns a new virtual machine instance configured with the provided components and settings.
func NewVM(loader ILoader, sequencer ISequencer, bc *bytecode.Bytecode, globals []objects.IObject, maxAllocations int64) (*VM, error) {
	if bc == nil {
		return nil, fmt.Errorf("bytecode is nil")
	}
	if maxAllocations < 1 {
		return nil, fmt.Errorf("max allocations must be greater than 0")
	}
	mainFn, err := bc.MainFunction()
	if err != nil {
		return nil, err
	}
	builtin, err := loader.ResolveBuiltinSymbols(bc.Constants())
	if err != nil {
		return nil, err
	}
	references, err := loader.ResolveSymbols(bc.References())
	if err != nil {
		return nil, err
	}
	if sequencer == nil {
		sequencer = NewSequencer()
	}
	if globals == nil {
		globals = make([]objects.IObject, globalsSize)
	}
	v := &VM{
		constants:      bc.Constants(),
		globals:        globals,
		sourceFiles:    bc.SourceFiles(),
		ip:             resetIp,
		maxAllocations: maxAllocations,
		suspend:        false,
		stack:          NewStack(stackSize),
		builtin:        builtin,
		references:     references,
		frames:         NewFrames(mainFn, maxFrames),
	}
	v.sequencer = sequencer.Create(v)
	return v, nil
}

// Abort sets the VM's internal abort flag to true, indicating that the execution should be halted.
func (v *VM) Abort() {
	v.abort = true
}

// Reset reinitializes the VM state, including stack, frames, instruction pointer, allocations, and error state.
func (v *VM) Reset() {
	v.stack.Reset()
	v.currFrame = v.frames.Head()
	v.currInstructions = v.currFrame.Instructions()
	v.frames.Clear()
	v.ip = resetIp
	v.allocations = v.maxAllocations + 1
	v.err = nil
	v.suspend = false
	v.abort = false
}

// Run executes the virtual machine's bytecode, managing the stack, frames, and instruction pointer state.
func (v *VM) Run() error {
	v.Reset()

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
		opcode, err := v.currInstructions.Get(v.ip)
		if err != nil {
			v.err = err
			return
		}
		opcode = opcode & sequenceMask
		fmt.Println("Executing instruction ", opcode, bytecode.OpcodeNames[opcode])
		v.sequencer[opcode]()
		if v.abort || v.suspend || v.err != nil {
			break
		}
	}
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
	cIdx, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	v.stack.Push(v.constants[cIdx])
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
	opcode, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	operator := objects.Operator(opcode)
	res, err := left.BinaryOp(operator, right)
	if err != nil {
		v.err = err
		return
	}
	v.stack.Push(res)
}

// doOpEqual pops two objects from the stack, compares them for equality, and pushes TrueValue or FalseValue based on the result.
func (v *VM) doOpEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	if left.Equals(right) {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// doOpNotEqual checks if the two topmost stack values are not equal and pushes true or false onto the stack accordingly.
func (v *VM) doOpNotEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	if left.Equals(right) {
		v.stack.Push(objects.FalseValue)
	} else {
		v.stack.Push(objects.TrueValue)
	}
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
	if operand.Boolean() {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// doOpBComplement performs a bitwise NOT operation on the top integer operand from the stack and pushes the result.
// If the operand is not an integer, it sets an error.
// If the object allocation limit is exceeded, it sets an allocation error.
func (v *VM) doOpBComplement() {
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		var res objects.IObject = objects.NewInt(^x.Value())
		v.allocations--
		if v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Push(res)
	default:
		v.err = fmt.Errorf("invalid operation: ^%s", operand.TypeName())
		return
	}
}

// doOpMinus applies the unary minus operation on the top element of the stack if it's an Int or Float object.
// It pushes the resulting negated value back onto the stack or sets an error if the operation is invalid.
func (v *VM) doOpMinus() {
	operand := v.stack.Pop()
	switch x := operand.(type) {
	case *objects.Int:
		var res objects.IObject = objects.NewInt(-x.Value())
		if v.allocations--; v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Push(res)
	case *objects.Float:
		var res objects.IObject = objects.NewFloat(-x.Value())
		if v.allocations--; v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Push(res)
	default:
		v.err = fmt.Errorf("invalid operation: -%s", operand.TypeName())
		return
	}
}

// doOpJumpFalsy performs a jump operation if the top stack value is falsy, updating the instruction pointer accordingly.
func (v *VM) doOpJumpFalsy() {
	v.ip += 2
	obj := v.stack.Pop()
	if obj.Boolean() {
		pos, err := v.currInstructions.Pos(v.ip, v.ip-1)
		if err != nil {
			v.err = err
			return
		}
		v.ip = pos - 1
	}
}

// doOpAndJump adjusts the instruction pointer and interacts with the stack to conditionally jump to a new position.
func (v *VM) doOpAndJump() {
	v.ip += 2
	obj := v.stack.Peek()
	if obj.Boolean() {
		pos, err := v.currInstructions.Pos(v.ip, v.ip-1)
		if err != nil {
			v.err = err
			return
		}
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
		pos, err := v.currInstructions.Pos(v.ip, v.ip-1)
		if err != nil {
			v.err = err
			return
		}
		v.ip = pos - 1
	}
}

// doOpJump updates the instruction pointer based on the position derived from the bytecode and adjusts for zero-based indexing.
func (v *VM) doOpJump() {
	pos, err := v.currInstructions.Pos(v.ip+2, v.ip+1)
	if err != nil {
		v.err = err
		return
	}
	v.ip = pos - 1
}

// doOpSetGlobal performs the SET_GLOBAL operation, updating a global variable with a value from the stack.
func (v *VM) doOpSetGlobal() {
	v.ip += 2
	pos, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	val := v.stack.Pop()
	v.globals[pos] = val
}

// doOpSetSelGlobal is a VM instruction to set a value in a global object using selectors for index assignment.
func (v *VM) doOpSetSelGlobal() {
	v.ip += 3
	globalIndex, err := v.currInstructions.Pos(v.ip-1, v.ip-2)
	if err != nil {
		v.err = err
		return
	}
	numSelectors, err := v.currInstructions.Get(v.ip)
	if err != nil {
		return
	}
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(int(numSelectors) + 1)
	if e := objects.IndexAssign(v.globals[globalIndex], val, selectors); e != nil {
		v.err = e
		return
	}
}

// doOpGetGlobal retrieves a global variable from the global pool and pushes its value onto the stack.
func (v *VM) doOpGetGlobal() {
	v.ip += 2
	globalIndex, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	val := v.globals[globalIndex]
	if val == nil {
		//v.err = fmt.Errorf("undefined global: %d", globalIndex)
		//return
	}
	v.stack.Push(val)
}

// doOpArray processes array-based operations in the VM by retrieving elements from the stack and creating a new array object.
// It increments the instruction pointer, extracts a specified number of elements, and enforces allocation limits.
// If the allocation limit is exceeded or an error occurs, it sets the error state of the VM.
func (v *VM) doOpArray() {
	v.ip += 2
	numElements, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	elements := v.stack.PopArrayElements(numElements)
	arr := objects.NewArray(elements)
	v.allocations--
	if v.allocations == 0 {
		v.err = objects.ErrObjectAllocLimit
		return
	}
	v.stack.Push(arr)
}

// doOpMap creates a new map object from key-value pairs popped from the stack and pushes the map onto the stack.
// It adjusts the instruction pointer and handles allocation limits or errors during creation.
func (v *VM) doOpMap() {
	v.ip += 2
	numElements, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	kv := v.stack.PopMapElements(numElements)
	m := objects.NewMap(kv)
	v.allocations--
	if v.allocations == 0 {
		v.err = objects.ErrObjectAllocLimit
		return
	}
	v.stack.Push(m)
}

// doOpError handles an operation error by converting the top stack value into an error object and updating the stack.
// Reduces the allocation count and sets an error if the allocation limit is reached.
func (v *VM) doOpError() {
	value := v.stack.Peek()
	var e objects.IObject = objects.NewError(value)
	v.allocations--
	if v.allocations == 0 {
		v.err = objects.ErrObjectAllocLimit
		return
	}
	v.stack.Set(e)
}

// doOpImmutable transforms the top stack value into its immutable form if it is an Array or Map.
// Reduces allocation counter and sets an error if the allocation limit is reached.
func (v *VM) doOpImmutable() {
	val := v.stack.Peek()
	switch value := val.(type) {
	case *objects.Array:
		immutableArray := objects.NewArrayImmutable(value.Values())
		v.allocations--
		if v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Set(immutableArray)
	case *objects.Map:
		immutableMap := objects.NewMapImmutable(value.Values())
		v.allocations--
		if v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Set(immutableMap)
	}
}

// doOpIndex processes an index operation by retrieving the index, the left operand, and obtaining the indexed value.
// If an error occurs during the operation, it sets the VM error accordingly. Pushes the result or UndefinedValue to the stack.
func (v *VM) doOpIndex() {
	index := v.stack.Pop()
	left := v.stack.Pop()
	val, err := left.IndexGet(index)
	if err != nil {
		if objects.Is(err, objects.ErrNotIndexable) {
			v.err = fmt.Errorf("not indexable: %s", index.TypeName())
			return
		}
		if objects.Is(err, objects.ErrInvalidIndexType) {
			v.err = fmt.Errorf("invalid index type: %s", index.TypeName())
			return
		}
		v.err = err
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
	var val objects.IObject = nil

	switch left := leftStack.(type) {
	case *objects.Array:
		if lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(left.Length())); err != nil {
			v.err = err
			return
		} else {
			val = objects.NewArray(left.Values()[lowIdx:highIdx])
		}
	case *objects.ArrayImmutable:
		if lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(left.Length())); err != nil {
			v.err = err
			return
		} else {
			val = objects.NewArray(left.Values()[lowIdx:highIdx])
		}
	case *objects.String:
		if lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(left.Length())); err != nil {
			v.err = err
			return
		} else {
			if val, err = objects.NewString(left.Value()[lowIdx:highIdx]); err != nil {
				v.err = err
				return
			}
		}
	case *objects.Bytes:
		if lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(left.Length())); err != nil {
			v.err = err
			return
		} else {
			val = objects.NewBytes(left.Value()[lowIdx:highIdx])
		}
	}

	if val != nil {
		v.allocations--
		if v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Push(val)
	}
}

// doOpCall executes a function call operation, handling arguments, callee type, and stack modifications accordingly.
func (v *VM) doOpCall() {
	numArgs, err := v.currInstructions.Get(v.ip + 1)
	if err != nil {
		v.err = err
		return
	}
	v.ip += 2
	value := v.stack.PeekOffset(-1 - numArgs)
	if !value.CanCall() {
		v.err = fmt.Errorf("not callable: %s", value.TypeName())
		return
	}
	spread, err := v.currInstructions.Get(v.ip + 2)
	if err != nil {
		v.err = err
		return
	}
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
			v.err = fmt.Errorf("not an array: %s", arrObj.TypeName())
			return
		}
	}
	if callee, ok := value.(*objects.FunctionCompiled); ok {
		if callee.VarArgs() {
			realArgs := callee.NumParameters() - 1
			v.stack.PushVarArgs(numArgs, realArgs)
		}
		if numArgs != callee.NumParameters() {
			numParams := callee.NumParameters()
			if callee.VarArgs() {
				numParams = callee.NumParameters() - 1
			}
			v.err = fmt.Errorf("wrong number of arguments: want>=%d, got=%d", numParams, numArgs)
			return
		}
		if v.currFrame.SameFunction(callee) {
			// recursive
			nextOp, err := v.currInstructions.Get(v.ip + 1)
			if err != nil {
				v.err = err
				return
			}
			nextNextOp, err := v.currInstructions.Get(v.ip + 2)
			if err != nil {
				v.err = err
				return
			}
			if byte(nextOp) == bytecode.OpReturn || (byte(nextOp) == bytecode.OpPop && byte(nextNextOp) == bytecode.OpReturn) {
				for p := 0; p < numArgs; p++ {
					o := v.stack.PeekOffset(-numArgs + p)
					v.stack.SetAbsolute(v.currFrame.BasePointer()+p, o)
				}
				v.stack.DecrementCount(numArgs + 1)
				v.ip = resetIp
				return
			}
		}
		v.currFrame.SetStartIP(v.ip)
		v.currFrame = v.frames.Get()
		v.currFrame.SetCompiledFunction(callee)
		v.currFrame.SetFreeVars(callee.Free())
		v.currFrame.SetBasePointer(v.stack.StackPointer() - numArgs)
		v.currInstructions = callee.Instructions()
		v.ip = resetIp
		if err = v.frames.Next(); err != nil {
			v.err = err
			return
		}
		v.stack.DecrementCount(numArgs + callee.NumLocals())
	} else {
		var args []objects.IObject
		args = append(args, v.stack.PeekArrayObject(numArgs)...)
		ret, err := value.Call(args...)
		v.stack.DecrementCount(numArgs + 1)
		if err != nil {
			if objects.Is(err, objects.ErrWrongNumArguments) {
				v.err = fmt.Errorf("wrong number of arguments in call to '%s'", value.TypeName())
				return
			}
			v.err = err
			return
		}
		if ret == nil {
			ret = objects.UndefinedValue
		}
		v.allocations--
		if v.allocations == 0 {
			v.err = objects.ErrObjectAllocLimit
			return
		}
		v.stack.Push(ret)
	}
}

// doOpReturn handles the return operation in the current VM frame, restoring the previous frame or terminating execution.
func (v *VM) doOpReturn() {
	v.ip++
	var retVal objects.IObject
	opcode, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	if opcode == 1 {
		retVal = v.stack.Peek()
	} else {
		retVal = objects.UndefinedValue
	}
	if v.frames.Index() > 1 {
		if err = v.frames.Previous(); err != nil {
			v.err = err
			return
		}
		v.currFrame = v.frames.GetPrev()
		v.currInstructions = v.currFrame.Instructions()
		v.ip = v.currFrame.StartIP()
		v.stack.SetStackPointer(v.frames.Get().BasePointer())
		v.stack.Set(retVal)
	} else {
		//log.Printf("returning from the root frame")
		fmt.Println("returning from the root frame")
		fmt.Println("stack:")
		v.stack.Print()
		v.abort = true
	}
}

// doOpDefineLocal defines a local variable in the current frame by using the local index from the bytecode instructions.
// It pops the value from the stack, calculates the destination slot using the base pointer and index, and sets the value.
// If the stack pointer is less than or equal to the destination slot, it updates the stack pointer to protect the new variable.
func (v *VM) doOpDefineLocal() {
	v.ip++
	localIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}

	val := v.stack.Pop()
	destSlot := v.currFrame.BasePointer() + localIndex
	v.stack.SetAbsolute(destSlot, val)

	// Assicura che lo stack pointer avanzi per "proteggere" la nuova variabile.
	if v.stack.StackPointer() <= destSlot {
		v.stack.SetStackPointer(destSlot + 1)
	}
}

// doOpSetLocal sets a local variable in the current function's call frame. It adjusts the stack as needed for the operation.
func (v *VM) doOpSetLocal() {
	localIndex, err := v.currInstructions.Get(v.ip + 1)
	if err != nil {
		v.err = err
		return
	}
	v.ip++
	val := v.stack.Pop()
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
	localIndex, err := v.currInstructions.Get(v.ip + 1)
	if err != nil {
		v.err = err
		return
	}
	numSelectors, err := v.currInstructions.Get(v.ip + 2)
	if err != nil {
		v.err = err
		return
	}
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
	if e := objects.IndexAssign(dst, val, selectors); e != nil {
		v.err = e
		return
	}
}

// doOpGetLocal retrieves a local variable from the stack based on its index in the current frame and pushes it onto the stack.
func (v *VM) doOpGetLocal() {
	v.ip++
	localIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	val := v.stack.PeekAbsolute(v.currFrame.BasePointer() + localIndex)
	if obj, ok := val.(*objects.ObjectPointer); ok {
		val = *obj.Value()
	}
	v.stack.Push(val)
}

// doOpGetBuiltin retrieves a builtin function by its index and pushes it onto the stack, incrementing the instruction pointer.
func (v *VM) doOpGetBuiltin() {
	v.ip++
	builtinIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	if builtinIndex < 0 || builtinIndex >= len(v.builtin) {
		v.err = fmt.Errorf("builtin index out of range: %d", builtinIndex)
		return
	}
	symbol := v.builtin[builtinIndex]
	v.stack.Push(symbol)
}

// doOpClosure handles the creation of a new compiled-function closure with captured free variables on the call stack.
func (v *VM) doOpClosure() {
	v.ip += 3
	constIndex, err := v.currInstructions.Pos(v.ip-1, v.ip-2)
	if err != nil {
		v.err = err
		return
	}
	numFree, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	fn, ok := v.constants[constIndex].(*objects.FunctionCompiled)
	if !ok {
		v.err = fmt.Errorf("not function: %s", fn.TypeName())
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
	v.allocations--
	if v.allocations == 0 {
		v.err = objects.ErrObjectAllocLimit
		return
	}
	v.stack.Push(cl)
}

// doOpGetFreePtr retrieves a free variable from the current frame based on the index in the bytecode and pushes it onto the stack.
func (v *VM) doOpGetFreePtr() {
	v.ip++
	freeIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	val := v.currFrame.FreeVarsIndex(freeIndex)
	v.stack.Push(val)
}

// doOpGetFree extracts a free variable value using its index, pushes it onto the stack, and increments the instruction pointer.
// If an error occurs while retrieving the instruction, it sets the VM's error and halts further execution.
func (v *VM) doOpGetFree() {
	v.ip++
	freeIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	val := *v.currFrame.FreeVarsIndex(freeIndex).Value()
	v.stack.Push(val)
}

// doOpSetFree increments the instruction pointer, retrieves a free variable index, and sets its value from the stack.
func (v *VM) doOpSetFree() {
	v.ip++
	freeIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	o := v.stack.Pop()
	v.currFrame.FreeVarsIndex(freeIndex).SetValue(o)
}

// doOpGetLocalPtr increments the instruction pointer, retrieves a local variable's value, and pushes it onto the stack.
// If the value is not already an ObjectPointer, it wraps it as a new ObjectPointer.
func (v *VM) doOpGetLocalPtr() {
	v.ip++
	localIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
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
	freeIndex, err := v.currInstructions.Get(v.ip - 1)
	if err != nil {
		v.err = err
		return
	}
	numSelectors, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(numSelectors + 1)
	fvi := v.currFrame.FreeVarsIndex(freeIndex)
	if err = objects.IndexAssign(*fvi.Value(), val, selectors); err != nil {
		v.err = err
		return
	}
}

// doOpIteratorInit initializes an iterator for the top stack object if it is iterable; otherwise, sets an error.
func (v *VM) doOpIteratorInit() {
	dst := v.stack.Pop()
	if !dst.CanIterate() {
		v.err = fmt.Errorf("not iterable: %s", dst.TypeName())
		return
	}
	iterator := dst.Iterate()
	v.allocations--
	if v.allocations == 0 {
		v.err = objects.ErrObjectAllocLimit
		return
	}
	v.stack.Push(iterator)
}

// doOpIteratorNext retrieves an iterator from the stack, advances it, and pushes a boolean indicating if more elements exist.
func (v *VM) doOpIteratorNext() {
	it := v.stack.Pop()
	iterator, ok := it.(objects.IIterator)
	if !ok {
		v.err = fmt.Errorf("not an iterator: %s", it.TypeName())
		return
	}
	hasMore := iterator.Next()
	if hasMore {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// doOpIteratorKey retrieves the current key from an iterator on the stack and pushes it back to the stack.
// If the top stack element is not an iterator, it sets an error in the VM.
func (v *VM) doOpIteratorKey() {
	it := v.stack.Pop()
	iterator, ok := it.(objects.IIterator)
	if !ok {
		v.err = fmt.Errorf("not an iterator: %s", it.TypeName())
		return
	}
	val := iterator.Key()
	v.stack.Push(val)
}

// doOpIteratorValue retrieves the current value from the top-most iterator on the stack and pushes it back onto the stack.
func (v *VM) doOpIteratorValue() {
	it := v.stack.Pop()
	iterator, ok := it.(objects.IIterator)
	if !ok {
		v.err = fmt.Errorf("not an iterator: %s", it.TypeName())
		return
	}
	val := iterator.Value()
	v.stack.Push(val)
}

// doOpReferences retrieves an attribute identified by its index, resolves it, and pushes it onto the stack.
func (v *VM) doOpReferences() {
	v.ip += 2
	nameIndex, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	if nameIndex < 0 || nameIndex >= len(v.references) {
		v.err = fmt.Errorf("invalid attribute index %d", nameIndex)
		return
	}
	symbol := v.references[nameIndex]
	v.stack.Push(symbol)
	return
}

// doOpSuspend sets the VM's `suspend` flag to true, pausing the execution of the bytecode operations.
func (v *VM) doOpSuspend() {
	v.suspend = true
}

// doOpUnknown is a fallback method executed when an unrecognized opcode is encountered in the instruction set.
// It retrieves the current instruction at the instruction pointer and sets an error indicating the unknown opcode.
func (v *VM) doOpUnknown() {
	pos, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	v.err = fmt.Errorf("unknown opcode: %d", pos)
}
