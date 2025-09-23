package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
	"github.com/markel1974/c64emu/src/vm/opcodes"
)

// Frame represents a structure used for maintaining function call frames in a virtual machine execution context.
// It encapsulates the execution state, free variables, instruction pointer, and base pointer of a function call.
type Frame struct {
	gk                   objects.IGateKeeper
	id                   int
	compiledFunction     *objects.Func
	freeVars             []*objects.ObjectPointer
	savedIp              uint
	basePointer          int
	instructions         *opcodes.Instructions
	deferredCalls        []*objects.Func
	savedReturnValues    []objects.IObject
	hasSavedReturnValues bool
	errSignal            func(err error)
}

// NewFrame creates and returns a new Frame instance with its instruction pointer initialized to -1.
func NewFrame(gk objects.IGateKeeper, id int, errSignal func(err error)) *Frame {
	return &Frame{
		gk:                   gk,
		id:                   id,
		savedIp:              resetIp,
		errSignal:            errSignal,
		savedReturnValues:    []objects.IObject{},
		hasSavedReturnValues: false,
		deferredCalls:        []*objects.Func{},
	}
}

// Id returns the unique identifier of the frame.
func (f *Frame) Id() int {
	return f.id
}

// Bind initializes the frame with the given instruction pointer, compiled function, and base pointer values.
func (f *Frame) Bind(startIp uint, compiledFunction *objects.Func, basePointer int) {
	f.savedIp = startIp
	f.basePointer = basePointer
	f.compiledFunction = compiledFunction
	f.instructions = f.compiledFunction.Instructions()
	f.freeVars = f.compiledFunction.Free()
}

// Fetch retrieves the next operation's instruction pointer and opcode ID from the bytecode sequence.
func (f *Frame) Fetch(ip uint) (opcodes.OpcodeId, uint) {
	headerBytes := f.Get8Reverse(ip)
	opcode := f.Get16Reverse(ip + opcodes.HeaderOpcodeIdBytes)
	return opcodes.OpcodeId(opcode), uint(headerBytes)
}

// Get8Reverse retrieves an 8-bit unsigned integer from the instructions at the specified index in the current frame.
// It signals an error if the retrieval fails and returns 0.
func (f *Frame) Get8Reverse(low uint) uint8 {
	v, ok := f.instructions.Get8Reverse(low)
	if !ok {
		f.errSignal(fmt.Errorf("invalid instruction 8 pointer: %d", low))
		return 0
	}
	return v
}

// Get16Reverse retrieves a 16-bit unsigned integer from instructions at specified indices x and y in the current frame.
// If an error occurs during retrieval, it signals the error and returns 0.
func (f *Frame) Get16Reverse(low uint) uint16 {
	v, ok := f.instructions.Get16Reverse(low)
	if !ok {
		f.errSignal(fmt.Errorf("invalid instruction 16 pointer: %d", low))
		return 0
	}
	return v
}

// Get32Reverse retrieves a 32-bit unsigned integer from the instructions at the specified index in the current frame. It signals an error if the retrieval fails and returns 0.
func (f *Frame) Get32Reverse(low uint) uint32 {
	v, ok := f.instructions.Get32Reverse(low)
	if !ok {
		f.errSignal(fmt.Errorf("invalid instruction 32 pointer: %d", low))
		return 0
	}
	return v
}

// Get64Reverse retrieves a 64-bit unsigned integer from instructions at the specified index in the current frame.
// It signals an error if the retrieval fails and returns 0.
func (f *Frame) Get64Reverse(low uint) uint64 {
	v, ok := f.instructions.Get64Reverse(low)
	if !ok {
		f.errSignal(fmt.Errorf("invalid instruction 32 pointer: %d", low))
		return 0
	}
	return v
}

// FreeVarsIndex retrieves the free variable at the specified index from the frame's free variables.
func (f *Frame) FreeVarsIndex(idx uint) *objects.ObjectPointer {
	if idx >= uint(len(f.freeVars)) {
		return nil
	}
	return f.freeVars[idx]
}

// SavedIP retrieves the current instruction pointer (ip) of the frame, indicating the execution position in bytecode.
func (f *Frame) SavedIP() uint {
	return f.savedIp
}

// BasePointer retrieves the base pointer value of the current frame, indicating the starting position of its variables.
func (f *Frame) BasePointer() int {
	return f.basePointer
}

// SourcePos returns the source position mapped to the given instruction pointer (ip) in the current frame's context.
func (f *Frame) SourcePos(ip int) int {
	return f.compiledFunction.SourcePos(ip)
}

// SameFunction compares the given Func with the Frame's compiled function and returns true if they are the same.
func (f *Frame) SameFunction(callee *objects.Func) bool {
	return f.compiledFunction == callee
}

// NumLocals returns the number of local variables required by the current frame's compiled function.	'
func (f *Frame) NumLocals() int {
	return f.compiledFunction.NumLocals()
}

// NumParameters returns the total number of parameters required by the compiled function of the current frame.
func (f *Frame) NumParameters() int {
	return f.compiledFunction.NumParameters()
}

// DeferredAdd appends a deferred call object to the frame's deferred calls queue for later execution.
func (f *Frame) DeferredAdd(obj objects.IObject) {
	switch objT := obj.(type) {
	case *objects.Func:
		f.deferredCalls = append(f.deferredCalls, objT)
	default:
		f.errSignal(fmt.Errorf("invalid operation: defer %s", obj.TypeName()))
	}
}

// DeferredPop removes and returns the last deferred call object from the frame's deferred calls queue.
func (f *Frame) DeferredPop() *objects.Func {
	numDeferred := len(f.deferredCalls)
	if numDeferred == 0 {
		f.errSignal(fmt.Errorf("no deferred calls in frame %d", f.id))
		return nil
	}
	target := numDeferred - 1
	lastCall := f.deferredCalls[target]
	f.deferredCalls = f.deferredCalls[:target]
	return lastCall
}

// HasDeferredCalls returns true if there are any deferred calls queued in the frame; otherwise, it returns false.
func (f *Frame) HasDeferredCalls() bool {
	return len(f.deferredCalls) > 0
}

// SaveParentReturnValues saves the return values from the parent frame into the current frame and marks them as saved.
func (f *Frame) SaveParentReturnValues(returnValues []objects.IObject) {
	f.savedReturnValues = returnValues
	f.hasSavedReturnValues = true
}

// HasParentReturnValues checks if the frame currently holds saved return values from a parent frame and returns true or false.
func (f *Frame) HasParentReturnValues() bool {
	return f.hasSavedReturnValues
}

// PopParentReturnValues resets and retrieves the saved return values from the parent frame, returning a slice of IObject.
func (f *Frame) PopParentReturnValues() []objects.IObject {
	ret := f.savedReturnValues
	f.savedReturnValues = []objects.IObject{}
	f.hasSavedReturnValues = false
	return ret
}
