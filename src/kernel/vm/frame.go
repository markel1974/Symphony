package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Frame represents a structure used for maintaining function call frames in a virtual machine execution context.
// It encapsulates the execution state, free variables, instruction pointer, and base pointer of a function call.
type Frame struct {
	id               int
	compiledFunction *objects.FuncCompiled
	freeVars         []*objects.ObjectPointer
	ip               int
	basePointer      int
	instructions     *objects.Instructions
	errSignal        func(err error)
}

// NewFunctionCallFrame creates and returns a new Frame instance with its instruction pointer initialized to -1.
func NewFunctionCallFrame(id int, errSignal func(err error)) *Frame {
	return &Frame{
		id:        id,
		ip:        resetIp,
		errSignal: errSignal,
	}
}

func (f *Frame) ID() int {
	return f.id
}

// Bind initializes the frame with the given instruction pointer, compiled function, and base pointer values.
func (f *Frame) Bind(startIp int, compiledFunction *objects.FuncCompiled, basePointer int) {
	f.ip = startIp
	f.basePointer = basePointer
	f.compiledFunction = compiledFunction
	f.instructions = f.compiledFunction.Instructions()
	f.freeVars = f.compiledFunction.Free()
}

// Get retrieves the integer value at the specified index from the instructions in the current frame.
func (f *Frame) Get(index int) int {
	v, err := f.instructions.Get(index)
	if err != nil {
		f.errSignal(err)
		return 0
	}
	return v
}

// Pos calculates and returns an instruction position based on the given indices in the frame's instruction set.
func (f *Frame) Pos(x int, y int) int {
	v, err := f.instructions.Pos(x, y)
	if err != nil {
		f.errSignal(err)
		return 0
	}
	return v
}

// FreeVarsIndex retrieves the free variable at the specified index from the frame's free variables.
func (f *Frame) FreeVarsIndex(idx int) *objects.ObjectPointer {
	return f.freeVars[idx]
}

// StartIP retrieves the current instruction pointer (ip) of the frame, indicating the execution position in bytecode.
func (f *Frame) StartIP() int {
	return f.ip
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
