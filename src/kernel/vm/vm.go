package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/compiler"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecodes"
	"github.com/markel1974/c64emu/src/kernel/vm/errors"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
	"github.com/markel1974/c64emu/src/kernel/vm/opcodes"
	"github.com/markel1974/c64emu/src/kernel/vm/stdlib"
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
	constants       []objects.IObject
	stack           []objects.IObject
	sp              int
	globals         []objects.IObject
	fileSet         *compiler.SourceFileSet
	frames          []*FunctionCallFrame
	framesIndex     int
	curFrame        *FunctionCallFrame
	curInstructions []byte
	ip              int
	abort           bool
	suspend         bool
	maxAllocations  int64
	allocations     int64
	err             error
	sequencer       []func()
}

// NewVM initializes and returns a new instance of the VM with provided sequencer, bytecode, globals, and max allocations.
func NewVM(sequencer ISequencer, bytecode *bytecodes.Bytecode, globals []objects.IObject, maxAllocations int64) *VM {
	if globals == nil {
		globals = make([]objects.IObject, GlobalsSize)
	}
	v := &VM{
		constants:      bytecode.Constants,
		sp:             0,
		globals:        globals,
		fileSet:        bytecode.FileSet,
		framesIndex:    1,
		ip:             -1,
		maxAllocations: maxAllocations,
		suspend:        false,
		stack:          make([]objects.IObject, StackSize),
		frames:         make([]*FunctionCallFrame, MaxFrames),
	}
	for i := range v.frames {
		v.frames[i] = NewFunctionCallFrame()
	}
	v.frames[0].SetCompiledFunction(bytecode.MainFunction)
	v.frames[0].ip = -1
	v.curFrame = v.frames[0]
	v.curInstructions = v.curFrame.Instructions()
	v.sequencer = sequencer.Create(v)
	return v
}

// Abort sets the VM's internal abort flag to true, signaling a termination or interruption of its current operation.
func (v *VM) Abort() {
	v.abort = true
}

// Run initializes the virtual machine's state and executes the current frame's instructions. Returns an error if execution fails.
func (v *VM) Run() error {
	v.sp = 0
	v.curFrame = v.frames[0]
	v.curInstructions = v.curFrame.Instructions()
	v.framesIndex = 1
	v.ip = -1
	v.allocations = v.maxAllocations + 1
	v.run()
	v.abort = false
	v.suspend = false
	if v.err != nil {
		filePos := v.fileSet.Position(v.curFrame.SourcePos(v.ip - 1))
		err := fmt.Errorf("runtime error %w at %s", v.err, filePos)
		for v.framesIndex > 1 {
			v.framesIndex--
			v.curFrame = v.frames[v.framesIndex-1]
			filePos = v.fileSet.Position(v.curFrame.SourcePos(v.curFrame.ip - 1))
			err = fmt.Errorf("%w at %s", err, filePos)
		}
		return err
	}
	return nil
}

// IsStackEmpty tests if the stack is empty or not.
//func (v *VM) IsStackEmpty() bool {
//	return v.sp == 0
//}

// run is the core execution loop of the virtual machine, iterating over and executing instructions until conditions are met.
func (v *VM) run() {
	for !v.abort || !v.suspend || v.err == nil {
		v.ip++
		v.sequencer[v.ip&sequenceMask]()
	}
}

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
	cIdx := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	v.stack[v.sp] = v.constants[cIdx]
	v.sp++
}

// doOpNull pushes an UndefinedValue onto the stack and increments the stack pointer.
func (v *VM) doOpNull() {
	v.stack[v.sp] = objects.UndefinedValue
	v.sp++
}

// doOpBinary handles the execution of a binary operation, updating the VM stack and evaluating the result.
// It increments the instruction pointer, performs a binary operation on the top two stack values, and handles errors.
// If an invalid operation is detected or the allocation limit is exceeded, an error is set in the VM.
func (v *VM) doOpBinary() {
	v.ip++
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	tok := objects.Operator(v.curInstructions[v.ip])
	res, e := left.BinaryOp(tok, right)
	if e != nil {
		v.sp -= 2
		if errors.Is(e, errors.ErrInvalidOperator) {
			v.err = fmt.Errorf("invalid operation: %s %d %s", left.TypeName(), tok, right.TypeName())
		}
		v.err = e
		return
	}
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp-2] = res
	v.sp--
}

// doOpEqual compares the top two values on the stack for equality and pushes TrueValue or FalseValue based on the result.
func (v *VM) doOpEqual() {
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2
	if left.Equals(right) {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

// doOpNotEqual compares the top two values on the stack for inequality and pushes the result (true or false) onto the stack.
func (v *VM) doOpNotEqual() {
	right := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2
	if left.Equals(right) {
		v.stack[v.sp] = objects.FalseValue
	} else {
		v.stack[v.sp] = objects.TrueValue
	}
	v.sp++
}

// doOpPop decreases the stack pointer by one during execution in the virtual machine. This effectively pops the top stack value.
func (v *VM) doOpPop() {
	v.sp--
}

// doOpTrue pushes the TrueValue object onto the VM stack and increments the stack pointer.
func (v *VM) doOpTrue() {
	v.stack[v.sp] = objects.TrueValue
	v.sp++
}

// doOpFalse pushes the predefined false value onto the stack and increments the stack pointer.
func (v *VM) doOpFalse() {
	v.stack[v.sp] = objects.FalseValue
	v.sp++
}

// doOpLNot performs a logical NOT operation on the top value of the stack, replacing it with the corresponding boolean value.
// If the operand is falsy, `objects.TrueValue` is pushed, otherwise `objects.FalseValue` is pushed.
func (v *VM) doOpLNot() {
	operand := v.stack[v.sp-1]
	v.sp--
	if operand.Falsy() {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

// doOpBComplement performs a bitwise complement operation on the top stack element, expecting it to be of type *objects.Int.
// It handles errors such as invalid operand type and allocation limit breaches.
func (v *VM) doOpBComplement() {
	operand := v.stack[v.sp-1]
	v.sp--
	switch x := operand.(type) {
	case *objects.Int:
		var res objects.IObject = objects.NewInt(^x.Value())
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	default:
		v.err = fmt.Errorf("invalid operation: ^%s", operand.TypeName())
		return
	}
}

// doOpMinus negates the top operand on the stack if it is an Int or Float, updates the stack, and handles allocation limits.
func (v *VM) doOpMinus() {
	operand := v.stack[v.sp-1]
	v.sp--

	switch x := operand.(type) {
	case *objects.Int:
		var res objects.IObject = objects.NewInt(-x.Value())
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	case *objects.Float:
		var res objects.IObject = objects.NewFloat(-x.Value())
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = res
		v.sp++
	default:
		v.err = fmt.Errorf("invalid operation: -%s", operand.TypeName())
		return
	}
}

// doOpJumpFalsy performs a conditional jump based on the falsy value of the top stack item, adjusting the instruction pointer.
func (v *VM) doOpJumpFalsy() {
	v.ip += 2
	v.sp--
	if v.stack[v.sp].Falsy() {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	}
}

// doOpAndJump adjusts the instruction pointer and stack pointer based on the falsiness of the top stack value.
func (v *VM) doOpAndJump() {
	v.ip += 2
	if v.stack[v.sp-1].Falsy() {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	} else {
		v.sp--
	}
}

// doOpOrJump updates the instruction pointer and stack pointer based on the falsy state of the top stack value.
func (v *VM) doOpOrJump() {
	v.ip += 2
	if v.stack[v.sp-1].Falsy() {
		v.sp--
	} else {
		pos := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
		v.ip = pos - 1
	}
}

// doOpJump adjusts the instruction pointer to the position specified by the next two bytes in the instruction sequence.
func (v *VM) doOpJump() {
	pos := int(v.curInstructions[v.ip+2]) | int(v.curInstructions[v.ip+1])<<8
	v.ip = pos - 1
}

// doOpSetGlobal updates a global variable by setting its value from the stack, using the global index derived from instructions.
func (v *VM) doOpSetGlobal() {
	v.ip += 2
	v.sp--
	globalIndex := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	v.globals[globalIndex] = v.stack[v.sp]
}

// doOpSetSelGlobal handles the assignment of a value to a global object with nested selectors.
// It updates the instruction pointer, extracts operands, modifies the stack, and performs indexed assignment.
// Errors during assignment, such as invalid selectors or non-assignable objects, are captured and set in the VM.
func (v *VM) doOpSetSelGlobal() {
	v.ip += 3
	globalIndex := int(v.curInstructions[v.ip-1]) | int(v.curInstructions[v.ip-2])<<8
	numSelectors := int(v.curInstructions[v.ip])

	// selectors and RHS value
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	if e := objects.IndexAssign(v.globals[globalIndex], val, selectors); e != nil {
		v.err = e
		return
	}
}

// doOpGetGlobal retrieves a global variable by its index and pushes its value onto the stack.
func (v *VM) doOpGetGlobal() {
	v.ip += 2
	globalIndex := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	val := v.globals[globalIndex]
	v.stack[v.sp] = val
	v.sp++
}

// doOpArray handles the creation of an array object by allocating elements from the stack, ensuring allocation limits.
func (v *VM) doOpArray() {
	v.ip += 2
	numElements := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	var elements []objects.IObject
	for i := v.sp - numElements; i < v.sp; i++ {
		elements = append(elements, v.stack[i])
	}
	v.sp -= numElements
	arr := objects.NewArray(elements)
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = arr
	v.sp++
}

// doOpMap creates a new map object from key-value pairs on the stack and places the map object back onto the stack.
// It also checks for object allocation limits and updates the instruction pointer and stack pointer accordingly.
func (v *VM) doOpMap() {
	v.ip += 2
	numElements := int(v.curInstructions[v.ip]) | int(v.curInstructions[v.ip-1])<<8
	kv := make(map[string]objects.IObject, numElements)
	for i := v.sp - numElements; i < v.sp; i += 2 {
		key := v.stack[i]
		value := v.stack[i+1]
		kv[key.(*objects.String).Value()] = value
	}
	v.sp -= numElements
	m := objects.NewMap(kv)
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = m
	v.sp++
}

// doOpError handles the creation of an `Error` object by wrapping the top stack item and replacing it with the error object.
// It decrements the allocation counter and sets an allocation limit error if the counter reaches zero.
func (v *VM) doOpError() {
	value := v.stack[v.sp-1]
	var e objects.IObject = objects.NewError(value)
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp-1] = e
}

// doOpImmutable converts a mutable array or map at the top of the stack to its immutable counterpart if possible.
// Reduces the allocation counter, setting an error if the allocation limit is exceeded.
func (v *VM) doOpImmutable() {
	val := v.stack[v.sp-1]
	switch value := val.(type) {
	case *objects.Array:
		immutableArray := objects.NewImmutableArray(value.Values())
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp-1] = immutableArray
	case *objects.Map:
		immutableMap := objects.NewImmutableMap(value.Values())
		v.allocations--
		if v.allocations == 0 {
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp-1] = immutableMap
	}
}

// doOpIndex handles the indexing operation on the stack by retrieving and validating indexed values or setting an error.
func (v *VM) doOpIndex() {
	index := v.stack[v.sp-1]
	left := v.stack[v.sp-2]
	v.sp -= 2
	val, err := left.IndexGet(index)
	if err != nil {
		if errors.Is(err, errors.ErrNotIndexable) {
			v.err = fmt.Errorf("not indexable: %s", index.TypeName())
			return
		}
		if errors.Is(err, errors.ErrInvalidIndexType) {
			v.err = fmt.Errorf("invalid index type: %s", index.TypeName())
			return
		}
		v.err = err
		return
	}
	if val == nil {
		val = objects.UndefinedValue
	}
	v.stack[v.sp] = val
	v.sp++
}

// doOpSliceIndex performs slicing operation on arrays, strings, or bytes based on indices from the stack and updates the stack.
// It validates index types and bounds, processes allocations, and handles errors for invalid operations.
func (v *VM) doOpSliceIndex() {
	highStack := v.stack[v.sp-1]
	lowStack := v.stack[v.sp-2]
	leftStack := v.stack[v.sp-3]
	v.sp -= 3

	var val objects.IObject = nil

	switch left := leftStack.(type) {
	case *objects.Array:
		if lowIdx, highIdx, err := v.checkBounds(lowStack, highStack, int64(left.Length())); err != nil {
			v.err = err
			return
		} else {
			val = objects.NewArray(left.Values()[lowIdx:highIdx])
		}
	case *objects.ImmutableArray:
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
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = val
		v.sp++
	}
}

// doOpCall handles the execution of a call operation, validating the callable object and managing arguments.
// Handles variadic calls, checks for recursion, and updates the call stack or returns any runtime errors encountered.
func (v *VM) doOpCall() {
	numArgs := int(v.curInstructions[v.ip+1])
	v.ip += 2
	value := v.stack[v.sp-1-numArgs]
	if !value.CanCall() {
		v.err = fmt.Errorf("not callable: %s", value.TypeName())
		return
	}
	if spread := int(v.curInstructions[v.ip+2]); spread == 1 {
		v.sp--
		switch arr := v.stack[v.sp].(type) {
		case *objects.Array:
			for _, item := range arr.Values() {
				v.stack[v.sp] = item
				v.sp++
			}
			numArgs += arr.Length() - 1
		case *objects.ImmutableArray:
			for _, item := range arr.Values() {
				v.stack[v.sp] = item
				v.sp++
			}
			numArgs += arr.Length() - 1
		default:
			v.err = fmt.Errorf("not an array: %s", arr.TypeName())
			return
		}
	}
	if callee, ok := value.(*objects.CompiledFunction); ok {
		if callee.VarArgs() {
			realArgs := callee.NumParameters() - 1
			if varArgs := numArgs - realArgs; varArgs >= 0 {
				numArgs = realArgs + 1
				args := make([]objects.IObject, varArgs)
				spStart := v.sp - varArgs
				for i := spStart; i < v.sp; i++ {
					args[i-spStart] = v.stack[i]
				}
				v.stack[spStart] = objects.NewArray(args)
				v.sp = spStart + 1
			}
		}
		if numArgs != callee.NumParameters() {
			numParams := callee.NumParameters()
			if callee.VarArgs() {
				numParams = callee.NumParameters() - 1
			}
			v.err = fmt.Errorf("wrong number of arguments: want>=%d, got=%d", numParams, numArgs)
			return
		}

		if v.curFrame.SameFunction(callee) {
			// recursive call
			nextOp := v.curInstructions[v.ip+1]
			if nextOp == opcodes.OpReturn || (nextOp == opcodes.OpPop && v.curInstructions[v.ip+2] == opcodes.OpReturn) {
				for p := 0; p < numArgs; p++ {
					v.stack[v.curFrame.basePointer+p] =
						v.stack[v.sp-numArgs+p]
				}
				v.sp -= numArgs + 1
				v.ip = -1 // reset IP to beginning of the frame
				return
			}
		}
		if v.framesIndex >= MaxFrames {
			v.err = errors.ErrStackOverflow
			return
		}
		v.curFrame.ip = v.ip
		v.curFrame = v.frames[v.framesIndex]
		v.curFrame.SetCompiledFunction(callee)
		v.curFrame.freeVars = callee.Free()
		v.curFrame.basePointer = v.sp - numArgs
		v.curInstructions = callee.Instructions()
		v.ip = -1
		v.framesIndex++
		v.sp = v.sp - numArgs + callee.NumLocals()
	} else {
		var args []objects.IObject
		args = append(args, v.stack[v.sp-numArgs:v.sp]...)
		ret, err := value.Call(args...)
		v.sp -= numArgs + 1
		if err != nil {
			if errors.Is(err, errors.ErrWrongNumArguments) {
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
			v.err = errors.ErrObjectAllocLimit
			return
		}
		v.stack[v.sp] = ret
		v.sp++
	}
}

// doOpReturn handles the execution of the return opcode, updating the VM's state and stack with the returned value.
func (v *VM) doOpReturn() {
	v.ip++
	var retVal objects.IObject
	if int(v.curInstructions[v.ip]) == 1 {
		retVal = v.stack[v.sp-1]
	} else {
		retVal = objects.UndefinedValue
	}
	v.framesIndex--
	v.curFrame = v.frames[v.framesIndex-1]
	v.curInstructions = v.curFrame.Instructions()
	v.ip = v.curFrame.ip
	v.sp = v.frames[v.framesIndex].basePointer
	v.stack[v.sp-1] = retVal
}

// doOpDefineLocal handles the definition of a local variable by storing a value from the stack into a calculated stack position.
func (v *VM) doOpDefineLocal() {
	v.ip++
	localIndex := int(v.curInstructions[v.ip])
	sp := v.curFrame.basePointer + localIndex
	val := v.stack[v.sp-1]
	v.sp--
	v.stack[sp] = val
}

// doOpSetLocal updates a local variable on the stack using its index and manages references for free variables if necessary.
func (v *VM) doOpSetLocal() {
	localIndex := int(v.curInstructions[v.ip+1])
	v.ip++
	sp := v.curFrame.basePointer + localIndex
	val := v.stack[v.sp-1]
	v.sp--
	if obj, ok := v.stack[sp].(*objects.ObjectPtr); ok {
		obj.SetValue(val)
		val = obj
	}
	v.stack[sp] = val
}

// doOpSetSelLocal handles the opcode for setting the value of a local variable using selectors.
func (v *VM) doOpSetSelLocal() {
	localIndex := int(v.curInstructions[v.ip+1])
	numSelectors := int(v.curInstructions[v.ip+2])
	v.ip += 2
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	dst := v.stack[v.curFrame.basePointer+localIndex]
	if obj, ok := dst.(*objects.ObjectPtr); ok {
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
	localIndex := int(v.curInstructions[v.ip])
	val := v.stack[v.curFrame.basePointer+localIndex]
	if obj, ok := val.(*objects.ObjectPtr); ok {
		val = *obj.Value()
	}
	v.stack[v.sp] = val
	v.sp++
}

// doOpGetBuiltin retrieves a built-in function by index and pushes it onto the stack, then increments the stack pointer.
func (v *VM) doOpGetBuiltin() {
	v.ip++
	builtinIndex := int(v.curInstructions[v.ip])
	v.stack[v.sp] = stdlib.GetBuiltin(builtinIndex)
	v.sp++
}

// doOpClosure handles the creation of a closure object by capturing free variables from the stack and setting it up.
func (v *VM) doOpClosure() {
	v.ip += 3
	constIndex := int(v.curInstructions[v.ip-1]) | int(v.curInstructions[v.ip-2])<<8
	numFree := int(v.curInstructions[v.ip])
	fn, ok := v.constants[constIndex].(*objects.CompiledFunction)
	if !ok {
		v.err = fmt.Errorf("not function: %s", fn.TypeName())
		return
	}
	free := make([]*objects.ObjectPtr, numFree)
	for i := 0; i < numFree; i++ {
		switch freeVar := (v.stack[v.sp-numFree+i]).(type) {
		case *objects.ObjectPtr:
			free[i] = freeVar
		default:
			free[i] = objects.NewObjectPtr(&v.stack[v.sp-numFree+i])
		}
	}
	v.sp -= numFree
	cl := objects.NewCompiledFunction(fn.Instructions(), fn.NumLocals(), fn.NumParameters(), fn.VarArgs(), nil, free)
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = cl
	v.sp++
}

// doOpGetFreePtr executes the opcode to retrieve a free variable pointer and push its value onto the stack.
func (v *VM) doOpGetFreePtr() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	val := v.curFrame.freeVars[freeIndex]
	v.stack[v.sp] = val
	v.sp++
}

// doOpGetFree retrieves the value of a free variable by its index and pushes it onto the stack.
func (v *VM) doOpGetFree() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	val := *v.curFrame.freeVars[freeIndex].Value()
	v.stack[v.sp] = val
	v.sp++
}

// doOpSetFree updates the value of a free variable in the current frame using the top value on the stack, then decrements the stack pointer.
func (v *VM) doOpSetFree() {
	v.ip++
	freeIndex := int(v.curInstructions[v.ip])
	v.curFrame.freeVars[freeIndex].SetValue(v.stack[v.sp-1])
	v.sp--
}

// doOpGetLocalPtr updates the instruction pointer, retrieves a local variable, and creates or updates an ObjectPtr for it.
func (v *VM) doOpGetLocalPtr() {
	v.ip++
	localIndex := int(v.curInstructions[v.ip])
	sp := v.curFrame.basePointer + localIndex
	val := v.stack[sp]
	var freeVar *objects.ObjectPtr
	if obj, ok := val.(*objects.ObjectPtr); ok {
		freeVar = obj
	} else {
		freeVar = objects.NewObjectPtr(&val)
		v.stack[sp] = freeVar
	}
	v.stack[v.sp] = freeVar
	v.sp++
}

// doOpSetSelFree performs a selective update on free variables using provided selectors and a right-hand side value.
// It adjusts the stack pointer and increments the instruction pointer after processing.
// If an error occurs during the update, it sets the error in the virtual machine.
func (v *VM) doOpSetSelFree() {
	v.ip += 2
	freeIndex := int(v.curInstructions[v.ip-1])
	numSelectors := int(v.curInstructions[v.ip])
	selectors := make([]objects.IObject, numSelectors)
	for i := 0; i < numSelectors; i++ {
		selectors[i] = v.stack[v.sp-numSelectors+i]
	}
	val := v.stack[v.sp-numSelectors-1]
	v.sp -= numSelectors + 1
	if err := objects.IndexAssign(*v.curFrame.freeVars[freeIndex].Value(), val, selectors); err != nil {
		v.err = err
		return
	}
}

// doOpIteratorInit initializes an iterator for the object at the top of the stack and handles errors if not iterable.
// If the object is iterable, it pushes the iterator onto the stack.
// Decrements the allocation counter and checks for allocation limits.
func (v *VM) doOpIteratorInit() {
	var iterator objects.IObject
	dst := v.stack[v.sp-1]
	v.sp--
	if !dst.CanIterate() {
		v.err = fmt.Errorf("not iterable: %s", dst.TypeName())
		return
	}
	iterator = dst.Iterate()
	v.allocations--
	if v.allocations == 0 {
		v.err = errors.ErrObjectAllocLimit
		return
	}
	v.stack[v.sp] = iterator
	v.sp++
}

// doOpIteratorNext executes the Next operation on an iterator and pushes the result (true/false) onto the stack.
func (v *VM) doOpIteratorNext() {
	iterator := v.stack[v.sp-1]
	v.sp--
	hasMore := iterator.(objects.IIterator).Next()
	if hasMore {
		v.stack[v.sp] = objects.TrueValue
	} else {
		v.stack[v.sp] = objects.FalseValue
	}
	v.sp++
}

// doOpIteratorKey retrieves the key from the top stack iterator and pushes it onto the stack.
func (v *VM) doOpIteratorKey() {
	iterator := v.stack[v.sp-1]
	v.sp--
	val := iterator.(objects.IIterator).Key()
	v.stack[v.sp] = val
	v.sp++
}

// doOpIteratorValue retrieves the current value from the iterator on the stack and updates the stack pointer accordingly.
func (v *VM) doOpIteratorValue() {
	iterator := v.stack[v.sp-1]
	v.sp--
	val := iterator.(objects.IIterator).Value()
	v.stack[v.sp] = val
	v.sp++
}

// doOpSuspend sets the suspend state of the VM to true, pausing its operation until resumed.
func (v *VM) doOpSuspend() {
	v.suspend = true
}

// doOpUnknown handles situations where the opcode is unrecognized by setting an appropriate error on the VM.
func (v *VM) doOpUnknown() {
	v.err = fmt.Errorf("unknown opcode: %d", v.curInstructions[v.ip])
}
