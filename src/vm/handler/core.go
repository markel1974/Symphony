package handler

import (
	"fmt"
	"io"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// resetIp is the instruction pointer value used to reset the Core's instruction pointer to the beginning of the main function.'
const (
	resetIp = 0 //
)

// Core represents a virtual machine that executes bytecode instructions, handles stack, and manages execution frames.
type Core struct {
	gk               objects.IGateKeeper
	id               uint
	shutdownSignal   func(id uint, err error)
	createCoreSignal func(id uint, callee *objects.Func, args []objects.IObject)
	stack            *Stack
	frames           *Frames
	currFrame        *Frame
	ip               uint
	sequencer        []*Decoder
	sequencerMask    int
	imports          *Imports
	constants        *Constants
	globals          *Globals
	runningIndex     int
}

// NewCore initializes and returns a new virtual machine instance configured with the provided components and settings.
func NewCore(gk objects.IGateKeeper, maxFrames int, stackSize int, id uint, shutdownSignal func(id uint, err error), createCoreSignal func(id uint, callee *objects.Func, args []objects.IObject)) *Core {
	v := &Core{
		gk:               gk,
		id:               id,
		shutdownSignal:   shutdownSignal,
		createCoreSignal: createCoreSignal,
		ip:               resetIp,
		runningIndex:     -1,
	}
	v.stack = NewStack(gk, stackSize, v.Shutdown)
	startInterval := id * uint(maxFrames)
	v.frames = NewFrames(gk, maxFrames, startInterval, v.Shutdown)
	return v
}

// Version returns the current version number of the Core instance as an integer.
func (v *Core) Version() int {
	return 1
}

// Id returns the unique identifier of the Core instance.
func (v *Core) Id() uint {
	return v.id
}

// Setup initializes the virtual machine with the provided bytecode and loader components.
func (v *Core) Setup(imports *Imports, constants *Constants, globals *Globals, seq ISequencer) error {
	executors, mask := seq.Executors()
	v.sequencerMask = mask
	v.sequencer = make([]*Decoder, len(executors))
	for i, s := range executors {
		err := s.Bind(v)
		if err != nil {
			return err
		}
		if v.sequencer[i], err = NewDecoder(s); err != nil {
			return err
		}
	}
	v.imports = imports
	v.constants = constants
	v.globals = globals
	return nil
}

// Shutdown gracefully terminates the execution of the Core, signaling an error if provided.
func (v *Core) Shutdown(err error) {
	v.shutdownSignal(v.id, err)
}

// CreateObjectPointer creates an object pointer within the current frame and returns it, or an undefined value if an error occurs.
func (v *Core) CreateObjectPointer(obj objects.IObject) objects.IObject {
	freeObjPtr, err := v.Factory().CreateObjectPointer(v.currFrame.Id(), obj)
	if err != nil {
		v.Shutdown(err)
		return v.gk.UndefinedValue()
	}
	return freeObjPtr
}

// CreateClosure creates and returns a new closure object from the given function object and its free variables.
func (v *Core) CreateClosure(fn *objects.Func, objRequired []objects.IObject) objects.IObject {
	cl, ok := fn.Copy(v.currFrame.Id(), 0).(*objects.Func)
	if !ok {
		v.Shutdown(fmt.Errorf("not a function: %s", fn.TypeName()))
		return v.gk.UndefinedValue()
	}
	if err := cl.FreeSet(v.currFrame.Id(), objRequired); err != nil {
		v.Shutdown(err)
		return v.gk.UndefinedValue()
	}
	return cl
}

// CreateError generates a new error object using the provided IObject as a source and assigns it to the current frame.
func (v *Core) CreateError(src objects.IObject) objects.IObject {
	errObj := v.Factory().NewError(v.currFrame.Id(), src.AsString())
	return errObj
}

// CreateSlice creates a slice using the provided high, low, and target objects and returns the resulting IObject.
// If an error occurs during slice creation, it sets the error on the Core and returns an undefined value.
func (v *Core) CreateSlice(highIdx int, lowIdx int, targetObj objects.IObject) objects.IObject {
	ret, err := v.Factory().CreateSlice(v.currFrame.Id(), highIdx, lowIdx, targetObj)
	if err != nil {
		v.Shutdown(err)
		return v.gk.UndefinedValue()
	}
	return ret
}

// StackPeek returns the object currently at the top of the stack without removing it.
func (v *Core) StackPeek() objects.IObject {
	return v.stack.Peek()
}

// StackPop removes and returns the top element from the Core's execution stack. It delegates the operation to the stack's Pop method.
func (v *Core) StackPop() objects.IObject {
	return v.stack.Pop()
}

// StackPush pushes the given IObject value onto the Core's stack.
func (v *Core) StackPush(value objects.IObject) {
	v.stack.Push(value)
}

// StackSet sets a value in the virtual machine's stack to the provided object.
func (v *Core) StackSet(value objects.IObject) {
	v.stack.Set(value)
}

// StackPeekBP retrieves an object from the stack at the given offset relative to the base pointer of the current frame.
func (v *Core) StackPeekBP(offset uint) objects.IObject {
	return v.stack.PeekAbsolute(uint(v.currFrame.BasePointer()) + offset)
}

// StackSetBP sets a value in the stack at the specified offset from the current frame's base pointer.
func (v *Core) StackSetBP(offset uint, value objects.IObject) {
	v.stack.SetAbsolute(uint(v.currFrame.BasePointer())+offset, value)
}

// StackSetSP sets the value of a stack element at the specified offset from the stack pointer (SP).
func (v *Core) StackSetSP(offset uint, value objects.IObject) {
	v.stack.SetOffset(offset, value)
}

// StackPeekSP retrieves the item at the specified offset from the stack, relative to the stack pointer.
func (v *Core) StackPeekSP(offset uint) objects.IObject {
	return v.stack.PeekOffset(offset)
}

// StackPopArray pops a specified number of elements from the stack and returns them as a slice of IObject.
func (v *Core) StackPopArray(numElements uint) objects.IObject {
	a := v.stack.PopArray(numElements)
	arrObj := v.Factory().NewArray(v.currFrame.Id(), a)
	return arrObj
}

// StackPopMap pops a specified number of key-value pairs from the stack and returns them as a map.
func (v *Core) StackPopMap(numElements uint) objects.IObject {
	m := v.stack.PopMap(numElements)
	mObj := v.Factory().NewMap(v.currFrame.Id(), m)
	return mObj
}

// StackPopStruct pops a specified number of key-value pairs from the stack and returns them as a map.
func (v *Core) StackPopStruct(numElements uint) objects.IObject {
	name, s := v.stack.PopStruct(numElements)
	sObj := v.Factory().NewStruct(v.currFrame.Id(), name, s)
	return sObj
}

// StackPopInterface pops `numElements` interfaces from the stack and wraps them into a new `objects.IObject` instance.
func (v *Core) StackPopInterface(numElements int) objects.IObject {
	concrete, iTable := v.stack.PopInterface(numElements)
	iObj := v.Factory().NewInterface(v.currFrame.Id(), concrete, iTable)
	return iObj
}

// StackDecrementCount reduces the count of items on the stack by the specified decrement amount.
func (v *Core) StackDecrementCount(decrement uint) {
	v.stack.DecrementCount(decrement)
}

// StackDecrement decreases the size of the stack by calling the Decrement method on the stack instance.
func (v *Core) StackDecrement() {
	v.stack.Decrement()
}

// ImportsGet fetches an imported object from the current frame using the specified index and returns it as an IObject.
func (v *Core) ImportsGet(idx uint) (objects.IObject, error) {
	return v.imports.Get(v.currFrame.Id(), idx)
}

// ConstantsGet fetches an imported object from the current frame using the specified index and returns it as an IObject.
func (v *Core) ConstantsGet(idx uint) (objects.IObject, error) {
	return v.constants.Get(v.currFrame.Id(), idx)
}

// GlobalsGet retrieves a global object by index from the Core's global container and returns it as an IObject.
func (v *Core) GlobalsGet(idx uint) (objects.IObject, error) {
	return v.globals.Get(idx)
}

// GlobalsSet assigns an object to the global store at the specified index.
func (v *Core) GlobalsSet(idx uint, obj objects.IObject) error {
	return v.globals.Set(idx, obj)
}

// FrameId returns the identifier of the current frame in the virtual machine.
func (v *Core) FrameId() int {
	return v.currFrame.Id()
}

// FrameDeferredAdd appends the given compiled function to the deferred call stack of the current frame.
func (v *Core) FrameDeferredAdd(obj objects.IObject) {
	v.currFrame.DeferredAdd(obj)
}

// FrameFreeVarsIndex retrieves the object pointer for the specified index from the current frame's free variables.
func (v *Core) FrameFreeVarsIndex(index uint) *objects.ObjectPointer {
	return v.currFrame.FreeVarsIndex(index)
}

// Factory returns the IGateKeeper instance associated with the Core.
func (v *Core) Factory() objects.IGateKeeper {
	return v.gk
}

// SetIp sets the virtual machine's instruction pointer to the specified value.
func (v *Core) SetIp(ip uint) {
	v.ip = ip
}

// GetIp retrieves the current instruction pointer value from the virtual machine.
func (v *Core) GetIp() uint {
	return v.ip
}

// Call executes a function or method with the specified number of arguments and handles variadic functions if applicable.
func (v *Core) Call(value objects.IObject, async bool, spread bool, numArgs int) {
	if numArgs < 0 {
		v.Shutdown(fmt.Errorf("invalid number of arguments: %d", numArgs))
		return
	}

	if spread {
		var args []objects.IObject
		obj := v.stack.Pop()
		switch z := obj.(type) {
		case *objects.Array:
			args = z.Values()
		default:
			v.Shutdown(fmt.Errorf("unexpected type (array required): %s", obj.TypeName()))
			return
		}
		if argsLen := len(args); argsLen > 0 {
			for _, item := range args {
				v.stack.Push(item)
			}
			numArgs += argsLen - 1
		}
	}

	switch ce := value.(type) {
	case *objects.Func:
		if ce.VarArgs() {
			np := ce.NumParameters()
			if np > 0 {
				np--
			}
			if varArgs := numArgs - np; varArgs > 0 {
				elements := v.stack.PopArray(uint(varArgs))
				v.stack.Push(v.gk.NewArray(v.currFrame.Id(), elements))
			}
		} else {
			if np := ce.NumParameters(); numArgs != np {
				v.Shutdown(fmt.Errorf("%s invalid arguments [vargs: %v] : want>=%d, got=%d", ce.Name(), ce.VarArgs(), np, numArgs))
				return
			}
		}
		if async {
			args := v.stack.PeekArray(uint(numArgs))
			v.createCoreSignal(v.id, ce, args)
		} else {
			v.prepareForCall(ce, numArgs)
		}
	default:
		args := v.stack.PeekArray(uint(numArgs))
		v.CallObject(value, numArgs, args...)
	}
}

// CallObject invokes a callable object with the given arguments and handles stack cleanup and error management.
func (v *Core) CallObject(value objects.IObject, numArgs int, args ...objects.IObject) {
	retCount, ret, err := value.Call(v.currFrame.Id(), args...)
	v.stack.DecrementCount(uint(numArgs) + 1)
	if err != nil {
		v.Shutdown(objects.ComputeCallError(err, value.TypeName()))
		return
	}
	switch retCount {
	case 0:
		v.stack.Push(v.gk.UndefinedValue())
		return
	case 1:
		if ret != nil {
			v.stack.Push(ret)
		} else {
			v.stack.Push(v.gk.UndefinedValue())
		}
		return
	default:
		container, ok := ret.(*objects.Array)
		if !ok {
			v.Shutdown(fmt.Errorf("invalid return count: %d", retCount))
			return
		}
		for _, item := range container.Values() {
			v.stack.Push(item)
		}
		return
	}
}

// Return concludes the execution of the current frame and handles return values, including handling deferred calls.
// If deferred calls are present, prepares the frame for execution of the first deferred call without immediate execution.
// Saves return values in the parent frame during 'defer' execution chains and recursively processes parent returns.
// Defaults to standard frame return if no deferred calls or chaining is present.
func (v *Core) Return(returnValues []objects.IObject) {
	// CASE 1: We are in a frame that is ending AND has pending 'defer' calls.
	if v.currFrame.HasDeferredCalls() {
		if deferredCall := v.currFrame.DeferredPop(); deferredCall != nil {
			// Prepare a frame for 'defer' call but don't execute it yet.
			deferredFrame := v.prepareForCall(deferredCall, 0)
			// Save return values of the current frame (parent) in a new 'defer' call frame. This creates the chain link.
			deferredFrame.SaveParentReturnValues(returnValues)
			// Break function. Core will execute the 'defer' call in the next cycle.
			return
		}
	}

	// CASE 2: We are in a frame that just finished (e.g. 'defer' closure) and must restore its parent's return flow.
	if v.currFrame.HasParentReturnValues() {
		savedReturnValues := v.currFrame.PopParentReturnValues()
		// A. Finalize the return of the current frame ('defer' frame).
		v.returnApply(returnValues)
		// B. Now that we're back in the parent frame, continue its return process.
		//	This creates recursion: calling 'Return' with parent values.
		v.Return(savedReturnValues)
	} else {
		// CASE 3: Standard return. No, pending 'defer' calls, and we're not inside a 'defer' chain.
		v.returnApply(returnValues)
	}
}

func (v *Core) returnApply(returnValues []objects.IObject) {
	shutdown := false
	prevIp := v.currFrame.SavedIP()
	leavingFrameBasePointer := v.currFrame.BasePointer()
	leavingFrameId := v.currFrame.Id()

	// Function arguments belong to the caller and should not be released when the function ends.
	numArgs := v.currFrame.NumParameters()
	totalArgs := numArgs + len(returnValues)
	var objectsToPreserve []objects.IObject

	if totalArgs > 0 {
		// Create a combined list, pre-allocating capacity for efficiency
		objectsToPreserve = make([]objects.IObject, 0, totalArgs)
		for _, obj := range returnValues {
			if obj.Frame() != objects.FrameStatic {
				objectsToPreserve = append(objectsToPreserve, obj)
			}
		}
		// add arguments which are located at the start of the current frame
		for i := 0; i < numArgs; i++ {
			obj := v.stack.PeekAbsolute(uint(leavingFrameBasePointer + i))
			if obj.Frame() != objects.FrameStatic {
				objectsToPreserve = append(objectsToPreserve, obj)
			}
		}
	}

	v.stack.ReleaseObjects(leavingFrameId, leavingFrameBasePointer, int(v.stack.StackPointer()), objectsToPreserve)

	if v.frames.CanMovePrevious() {
		v.frames.MovePrevious()
		v.currFrame = v.frames.Previous()
		v.ip = prevIp
	} else {
		shutdown = true
	}

	v.stack.SetStackPointer(uint(leavingFrameBasePointer))

	if lRet := len(returnValues); lRet > 0 {
		for i := lRet - 1; i >= 0; i-- {
			v.stack.Push(returnValues[i])
		}
	} else {
		v.stack.Push(v.gk.UndefinedValue())
	}

	if shutdown {
		v.Shutdown(nil)
	}
}

// ReseIp resets the instruction pointer of the virtual machine to its initial reset state defined by `resetIp`.
func (v *Core) ReseIp() {
	v.ip = resetIp
}

// GetReturnValue returns the value from the top of the stack as an interface value.
func (v *Core) GetReturnValue(idx int) interface{} {
	obj := v.stack.PeekAbsolute(uint(idx))
	if obj == nil {
		return nil
	}
	return obj.AsInterface()
}

// GetReturnValues returns the values from the top of the stack as an array of interface values.
func (v *Core) GetReturnValues() []interface{} {
	values := v.stack.StackPointer()
	if values == 0 {
		return nil
	}
	out := make([]interface{}, values)
	for x := 0; x < int(values); x++ {
		out[x] = v.GetReturnValue(x)
	}
	return out
}

// Rewrite updates the global function mapping with the given JIT-compiled function for the specified identifier.
func (v *Core) Rewrite(id uint, jit *objects.FuncJit) error {
	return v.globals.Set(id, jit)
}

// Initialize sets up the initial state of the virtual machine and prepares it to execute the given main function.
func (v *Core) Initialize(runningIndex int, mainFn *objects.Func, args []objects.IObject) error {
	v.runningIndex = runningIndex
	v.ip = resetIp
	v.gk.Reset()
	v.stack.Reset()
	v.frames.Reset()
	v.currFrame = v.frames.Head()
	v.currFrame.Bind(v.ip, mainFn, 0)
	v.stack.SetStackPointer(uint(v.currFrame.NumLocals()))
	if v.currFrame.NumParameters() != len(args) {
		return fmt.Errorf("[%s] wrong number of arguments provided: want=%d, got=%d", mainFn.Name(), v.currFrame.NumParameters(), len(args))
	}
	for idx, arg := range args {
		v.stack.SetAbsolute(uint(idx), arg)
	}
	return nil
}

// Update sets the running index of the Core instance to the provided value.
func (v *Core) Update(runningIndex int) {
	v.runningIndex = runningIndex
}

// Finalize cleans up resources by releasing all stack items and resetting the running index to its initial state.
func (v *Core) Finalize() int {
	v.stack.ReleaseAll()
	runningIndex := v.runningIndex
	v.runningIndex = -1
	return runningIndex
}

// Execute processes the current instruction in the Core, updating the instruction pointer and executing the decoded operation.
func (v *Core) Execute() {
	opcode, headerSize := v.currFrame.Fetch(v.ip)
	decoder := v.sequencer[opcode&v.sequencerMask]
	v.ip += headerSize + decoder.OperandsSize() - 1 //zero-based index
	decoder.DecodeReverse(v.currFrame, v.ip)
	v.ip++ //next instruction
	//log.Printf("Executing instruction opcode: %d name: %s ip: %d decoded: %v", opcode, decoder.Name(), v.ip, decoder.decodedOperands[:decoder.operandsSize])
	decoder.Execute()
}

// callCompiled sets up a new execution frame for a compiled function and manages stack allocation for local variables.
// Callee specifies the compiled function to be executed, and numArgs determines the number of arguments passed.
// It reserves stack space for all local variables and adjusts the instruction pointer accordingly.
func (v *Core) prepareForCall(callee *objects.Func, numArgs int) *Frame {
	// 1. Calculate the new basePointer safely, anchoring it to the caller's frame.
	//	The new "floor" begins exactly where the caller's local variable space ends.
	bp := v.currFrame.BasePointer() + v.currFrame.NumLocals()

	v.stack.CopyOffset(uint(bp), uint(numArgs))

	// 3. Advance to the next frame
	v.currFrame = v.frames.Current()
	v.frames.MoveNext()

	// 4. Bind the new frame to function and correct basePointer
	v.currFrame.Bind(v.GetIp(), callee, bp)

	// 5. Set stack pointer to include arguments and space for new locals
	v.stack.SetStackPointer(uint(bp + numArgs + callee.NumLocals()))
	v.ReseIp()

	return v.currFrame
}

// SourcePos returns the source code position corresponding to the current instruction pointer in the execution frame.
func (v *Core) SourcePos() int {
	return v.currFrame.SourcePos(int(v.ip) - 1)
}

// FramesUnroll retrieves and returns all the frames in the current execution context as a slice of Frame pointers.
func (v *Core) FramesUnroll() []*Frame {
	return v.frames.Unroll()
}

// FramesMax returns the maximum frame index accessed during execution as a uint64.
func (v *Core) FramesMax() uint64 {
	return uint64(v.frames.Max())
}

// Print writes the representation of the stack to the provided io.Writer.
func (v *Core) Print(writer io.Writer) {
	v.stack.Print(writer)
}
