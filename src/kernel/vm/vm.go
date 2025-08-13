package vm

import (
	"fmt"
	"strings"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/stdlib"
)

const (
	globalsSize = 1024
	stackSize   = 2048
	maxFrames   = 1024
)

// sequenceLen defines the length of a sequence using a bitwise shift for efficient computation.
// sequenceMask is a bit mask derived from sequenceLen to efficiently limit values within the sequence range.
const (
	sequenceLen  = 1 << 8
	sequenceMask = sequenceLen - 1
)

// ISequencer defines an interface for creating a sequence of functions for a virtual machine.
type ISequencer interface {
	Create(vm *VM) []func()
}

// VM represents a virtual machine responsible for executing bytecode instructions.
type VM struct {
	sourceFiles      *bytecode.Files
	constants        []objects.IObject
	globals          []objects.IObject
	stack            *Stack
	ip               int
	frames           *Frames
	currFrame        *Frame
	currInstructions *objects.Instructions
	abort            bool
	suspend          bool
	maxAllocations   int64
	allocations      int64
	err              error
	sequencer        []func()
	linker           modules.IModuleGetter
}

// NewVM initializes and returns a new instance of the VM with provided sequencer, bytecode, globals, and max allocations.
func NewVM(linker modules.IModuleGetter, sequencer ISequencer, bc *bytecode.Bytecode, globals []objects.IObject, maxAllocations int64) (*VM, error) {
	if bc == nil {
		return nil, fmt.Errorf("bytecode is nil")
	}
	if sequencer == nil {
		sequencer = NewSequencer()
	}
	if globals == nil {
		globals = make([]objects.IObject, globalsSize)
	}
	if maxAllocations < 1 {
		return nil, fmt.Errorf("max allocations must be greater than 0")
	}
	v := &VM{
		constants:      bc.Constants(),
		globals:        globals,
		sourceFiles:    bc.SourceFiles(),
		ip:             -1,
		maxAllocations: maxAllocations,
		suspend:        false,
		stack:          NewStack(stackSize),
		linker:         linker,
	}
	main, err := bc.MainFunction()
	if err != nil {
		return nil, err
	}
	v.frames = NewFrames(main, maxFrames)
	v.currFrame = v.frames.Head()
	v.currInstructions = v.currFrame.Instructions()
	v.sequencer = sequencer.Create(v)
	return v, nil
}

// Abort sets the VM's internal abort flag to true, signaling a termination or interruption of its current operation.
func (v *VM) Abort() {
	v.abort = true
}

// Run initializes the virtual machine's state and executes the current frame's instructions. Returns an error if execution fails.
func (v *VM) Run() error {
	v.stack.Clear()
	v.currFrame = v.frames.Head()
	v.currInstructions = v.currFrame.Instructions()
	v.frames.Clear()
	v.ip = -1
	v.allocations = v.maxAllocations + 1
	v.run()
	v.abort = false
	v.suspend = false
	if v.err != nil {
		filePos, _ := v.sourceFiles.Position(v.currFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for _, frame := range v.frames.Unroll() {
			filePos, _ = v.sourceFiles.Position(frame.SourcePos(frame.IP() - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
}

// run is the core execution loop of the virtual machine, iterating over and executing instructions until conditions are met.
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

// checkBounds calculates and validates lower and upper bounds for slicing operations, ensuring they fit within valid indices.
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

// doOpConstant fetches a constant from the constants pool and pushes it onto the stack.
func (v *VM) doOpConstant() {
	v.ip += 2
	//cIdx := int(v.currInstructions.Get(v.ip)) | int(v.currInstructions.Get(v.ip-1))<<8
	cIdx, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	v.stack.Push(v.constants[cIdx])
}

// doOpNull pushes an UndefinedValue onto the stack and increments the stack pointer.
func (v *VM) doOpNull() {
	v.stack.Push(objects.UndefinedValue)
}

// in vm.go
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

// doOpEqual compares the top two values on the stack for equality and pushes TrueValue or FalseValue based on the result.
func (v *VM) doOpEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	if left.Equals(right) {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// doOpNotEqual compares the top two values on the stack for inequality and pushes the result (true or false) onto the stack.
func (v *VM) doOpNotEqual() {
	right := v.stack.Pop()
	left := v.stack.Pop()
	if left.Equals(right) {
		v.stack.Push(objects.FalseValue)
	} else {
		v.stack.Push(objects.TrueValue)
	}
}

// doOpPop decreases the stack pointer by one during execution in the virtual machine. This effectively pops the top stack value.
func (v *VM) doOpPop() {
	v.stack.Decrement()
	//v.sp--
}

// doOpTrue pushes the TrueValue object onto the VM stack and increments the stack pointer.
func (v *VM) doOpTrue() {
	v.stack.Push(objects.TrueValue)
}

// doOpFalse pushes the predefined false value onto the stack and increments the stack pointer.
func (v *VM) doOpFalse() {
	v.stack.Push(objects.FalseValue)
}

// doOpLNot performs a logical NOT operation on the top value of the stack, replacing it with the corresponding boolean value.
// If the operand is falsy, `objects.TrueValue` is pushed, otherwise `objects.FalseValue` is pushed.
func (v *VM) doOpLNot() {
	operand := v.stack.Pop()
	if operand.Boolean() {
		v.stack.Push(objects.TrueValue)
	} else {
		v.stack.Push(objects.FalseValue)
	}
}

// doOpBComplement performs a bitwise complement operation on the top stack element, expecting it to be of type *objects.Int.
// It handles errors such as invalid operand type and allocation limit breaches.
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

// doOpMinus negates the top operand on the stack if it is an Int or Float, updates the stack, and handles allocation limits.
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

// doOpJumpFalsy performs a conditional jump based on the falsy value of the top stack item, adjusting the instruction pointer.
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

// doOpAndJump adjusts the instruction pointer and stack pointer based on the falsiness of the top stack value.
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

// doOpOrJump updates the instruction pointer and stack pointer based on the falsy state of the top stack value.
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

// doOpJump adjusts the instruction pointer to the position specified by the next two bytes in the instruction sequence.
func (v *VM) doOpJump() {
	pos, err := v.currInstructions.Pos(v.ip+2, v.ip+1)
	if err != nil {
		v.err = err
		return
	}
	v.ip = pos - 1
}

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

// doOpSetSelGlobal handles the assignment of a value to a global object with nested selectors.
// It updates the instruction pointer, extracts operands, modifies the stack, and performs indexed assignment.
// Errors during assignment, such as invalid selectors or non-assignable objects, are captured and set in the VM.
func (v *VM) doOpSetSelGlobal() {
	v.ip += 3
	//globalIndex := int(v.currInstructions.Get(v.ip-1)) | int(v.currInstructions.Get(v.ip-2))<<8
	globalIndex, err := v.currInstructions.Pos(v.ip-1, v.ip-2)
	if err != nil {
		v.err = err
		return
	}
	numSelectors, err := v.currInstructions.Get(v.ip)
	if err != nil {
		return
	}
	// selectors and RHS value
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack.PeekOffset(-numSelectors + i)
	}
	val := v.stack.PeekOffset(-numSelectors - 1)
	v.stack.DecrementCount(int(numSelectors) + 1)
	//v.sp -= int(numSelectors) + 1
	if e := objects.IndexAssign(v.globals[globalIndex], val, selectors); e != nil {
		v.err = e
		return
	}
}

// doOpGetGlobal retrieves a global variable by its index and pushes its value onto the stack.
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

// doOpArray handles the creation of an array object by allocating elements from the stack, ensuring allocation limits.
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

// doOpMap creates a new map object from key-value pairs on the stack and places the map object back onto the stack.
// It also checks for object allocation limits and updates the instruction pointer and stack pointer accordingly.
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

// doOpError handles the creation of an `Error` object by wrapping the top stack item and replacing it with the error object.
// It decrements the allocation counter and sets an allocation limit error if the counter reaches zero.
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

// doOpImmutable converts a mutable array or map at the top of the stack to its immutable counterpart if possible.
// Reduces the allocation counter, setting an error if the allocation limit is exceeded.
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

// doOpIndex handles the indexing operation on the stack by retrieving and validating indexed values or setting an error.
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

// doOpSliceIndex performs slicing operation on arrays, strings, or bytes based on indices from the stack and updates the stack.
// It validates index types and bounds, processes allocations, and handles errors for invalid operations.
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

// doOpCall handles the execution of a call operation, validating the callable object and managing arguments.
// Handles variadic calls, checks for recursion, and updates the call stack or returns any runtime errors encountered.
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
			// recursive call
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
				v.ip = -1 // reset IP to beginning of the frame
				return
			}
		}
		v.currFrame.SetIP(v.ip)
		v.currFrame = v.frames.Get()
		v.currFrame.SetCompiledFunction(callee)
		v.currFrame.SetFreeVars(callee.Free())
		v.currFrame.SetBasePointer(v.stack.StackPointer() - numArgs)
		v.currInstructions = callee.Instructions()
		v.ip = -1
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

// doOpReturn handles the execution of the return opcode, updating the VM's state and stack with the returned value.
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
		v.ip = v.currFrame.IP()
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

// doOpDefineLocal defines a local variable in the virtual machine's current frame and stack.
// It retrieves the local index, assigns the top value of the stack to the computed slot, and adjusts the stack pointer.
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

// doOpSetLocal updates a local variable on the stack using its index.
// If the existing local variable is an ObjectPointer (a free variable captured by a closure),
// it updates the value inside the pointer instead of replacing the pointer itself.
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

// doOpSetSelLocal handles the opcode for setting the value of a local variable using selectors.
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

// doOpGetLocal retrieves a local variable from the stack and places it at the top of the stack, incrementing the stack pointer.
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

// doOpGetBuiltin retrieves a built-in function by index and pushes it onto the stack, then increments the stack pointer.
func (v *VM) doOpGetBuiltin() {
	v.ip++
	builtinIndex, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	v.stack.Push(stdlib.GetBuiltin(builtinIndex))
}

// doOpClosure handles the creation of a closure object by capturing free variables from the stack and setting it up.
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

// doOpGetFreePtr executes the opcode to retrieve a free variable pointer and push its value onto the stack.
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

// doOpGetFree retrieves the value of a free variable by its index and pushes it onto the stack.
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

// doOpSetFree updates the value of a free variable in the current frame using the top value on the stack, then decrements the stack pointer.
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

// doOpGetLocalPtr retrieves a local pointer from the stack or creates one if it doesn't exist, updating the stack state.
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

// doOpSetSelFree performs a selective update on free variables using provided selectors and a right-hand side value.
// It adjusts the stack pointer and increments the instruction pointer after processing.
// If an error occurs during the update, it sets the error in the virtual machine.
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

// doOpIteratorInit initializes an iterator for the object at the top of the stack and handles errors if not iterable.
// If the object is iterable, it pushes the iterator onto the stack.
// Decrements the allocation counter and checks for allocation limits.
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

// doOpIteratorNext executes the Next operation on an iterator and pushes the result (true/false) onto the stack.
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

// doOpIteratorKey retrieves the key from the top stack iterator and pushes it onto the stack.
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

// doOpIteratorValue retrieves the current value from the iterator on the stack and updates the stack pointer accordingly.
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

func (v *VM) doOpGetAttr() {
	if v.linker == nil {
		v.err = fmt.Errorf("no module loaded")
		return
	}
	v.ip += 2
	nameIndex, err := v.currInstructions.Pos(v.ip, v.ip-1)
	if err != nil {
		v.err = err
		return
	}
	attrName, ok := v.constants[nameIndex].(*objects.String)
	if !ok {
		v.err = fmt.Errorf("invalid attribute name constant")
		return
	}
	values := strings.Split(attrName.Value(), ".")
	if len(values) != 2 {
		v.err = fmt.Errorf("invalid attribute name")
		return
	}
	packageName := values[0]
	symbolName := values[1]
	container := v.linker.Get(packageName)
	if container == nil {
		v.err = fmt.Errorf("module '%s' not found", packageName)
		return
	}
	symbol, ok := container.Symbol(symbolName)
	if !ok {
		v.err = fmt.Errorf("symbol '%s' not found", symbolName)
		return
	}
	v.stack.Push(symbol)
}

// doOpSuspend sets the suspend state of the VM to true, pausing its operation until resumed.
func (v *VM) doOpSuspend() {
	v.suspend = true
}

// doOpUnknown handles situations where the opcode is unrecognized by setting an appropriate error on the VM.
func (v *VM) doOpUnknown() {
	pos, err := v.currInstructions.Get(v.ip)
	if err != nil {
		v.err = err
		return
	}
	v.err = fmt.Errorf("unknown opcode: %d", pos)
}
