package core

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Frame represents a structure used for maintaining function call frames in a virtual machine execution context.
// It encapsulates the execution state, free variables, instruction pointer, and base pointer of a function call.
type Frame struct {
	gk               objects.IGateKeeper
	id               int
	compiledFunction *objects.FuncCompiled
	freeVars         []*objects.ObjectPointer
	savedIp          int
	basePointer      int
	instructions     *objects.Instructions
	errSignal        func(err error)
}

// NewFrame creates and returns a new Frame instance with its instruction pointer initialized to -1.
func NewFrame(gk objects.IGateKeeper, id int, errSignal func(err error)) *Frame {
	return &Frame{
		gk:        gk,
		id:        id,
		savedIp:   resetIp,
		errSignal: errSignal,
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

// Get8 retrieves an 8-bit unsigned integer from the instructions at the specified index in the current frame.
// It signals an error if the retrieval fails and returns 0.
func (f *Frame) Get8(index int) uint8 {
	v, err := f.instructions.Get8(index)
	if err != nil {
		f.errSignal(err)
		return 0
	}
	return v
}

// Get16 retrieves a 16-bit unsigned integer from instructions at specified indices x and y in the current frame.
// If an error occurs during retrieval, it signals the error and returns 0.
func (f *Frame) Get16(low int) uint16 {
	v, err := f.instructions.Get16(low)
	if err != nil {
		f.errSignal(err)
		return 0
	}
	return v
}

// FreeVarsIndex retrieves the free variable at the specified index from the frame's free variables.
func (f *Frame) FreeVarsIndex(idx int) *objects.ObjectPointer {
	if idx < 0 || idx >= len(f.freeVars) {
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
