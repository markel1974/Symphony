package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Frame represents a structure used for maintaining function call frames in a virtual machine execution context.
// It encapsulates the execution state, free variables, instruction pointer, and base pointer of a function call.
type Frame struct {
	gk                   objects.IGateKeeper
	id                   int
	compiledFunction     *objects.FuncCompiled
	freeVars             []*objects.ObjectPointer
	savedIp              int
	basePointer          int
	instructions         *objects.Instructions
	deferredCalls        []*objects.FuncCompiled
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
		deferredCalls:        []*objects.FuncCompiled{},
	}
}

// Id returns the unique identifier of the frame.
func (f *Frame) Id() int {
	return f.id
}

// Bind initializes the frame with the given instruction pointer, compiled function, and base pointer values.
func (f *Frame) Bind(startIp int, compiledFunction *objects.FuncCompiled, basePointer int) {
	f.savedIp = startIp
	f.basePointer = basePointer
	f.compiledFunction = compiledFunction
	f.instructions = f.compiledFunction.Instructions()
	f.freeVars = f.compiledFunction.Free()
}

// Fetch retrieves the next instruction pointer and its corresponding OpcodeId from the frame's instructions.
func (f *Frame) Fetch(ip int) (int, bytecode.OpcodeId) {
	offset := uint(ip + bytecode.HeaderSizeBytes)
	headerBytes := f.Get8Reverse(offset)
	offset += bytecode.HeaderOpcodeIdBytes
	opcode := f.Get16Reverse(offset)
	ip += int(headerBytes)
	return ip, bytecode.OpcodeId(opcode)
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
func (f *Frame) SavedIP() int {
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

// SameFunction compares the given FuncCompiled with the Frame's compiled function and returns true if they are the same.
func (f *Frame) SameFunction(callee *objects.FuncCompiled) bool {
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
func (f *Frame) DeferredAdd(call *objects.FuncCompiled) {
	f.deferredCalls = append(f.deferredCalls, call)
}

// DeferredPop removes and returns the last deferred call object from the frame's deferred calls queue.
func (f *Frame) DeferredPop() *objects.FuncCompiled {
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

func (f *Frame) SaveParentReturnValues(returnValues []objects.IObject) {
	f.savedReturnValues = returnValues
	f.hasSavedReturnValues = true
}

func (f *Frame) HasParentReturnValues() bool {
	return f.hasSavedReturnValues
}

func (f *Frame) PopParentReturnValues() []objects.IObject {
	ret := f.savedReturnValues
	f.savedReturnValues = []objects.IObject{}
	f.hasSavedReturnValues = false
	return ret
}
