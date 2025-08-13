package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Frame represents a structure used for maintaining function call frames in a virtual machine execution context.
// It encapsulates the execution state, free variables, instruction pointer, and base pointer of a function call.
type Frame struct {
	compiledFunction *objects.FunctionCompiled
	freeVars2        []*objects.ObjectPointer
	ip               int
	basePointer      int
}

// NewFunctionCallFrame creates and returns a new Frame instance with its instruction pointer initialized to -1.
func NewFunctionCallFrame() *Frame {
	return &Frame{
		ip: -1,
	}
}

// FreeVarsIndex retrieves the free variable at the specified index from the frame's free variables.
func (f *Frame) FreeVarsIndex(idx int) *objects.ObjectPointer {
	return f.freeVars2[idx]
}

// SetFreeVars sets the free variables of the frame to the provided slice of ObjectPointer instances.
func (f *Frame) SetFreeVars(freeVars []*objects.ObjectPointer) {
	f.freeVars2 = freeVars
}

// IP returns the current instruction pointer stored in the frame.
func (f *Frame) IP() int {
	return f.ip
}

// SetIP updates the instruction pointer (IP) of the Frame to the provided value.
func (f *Frame) SetIP(ip int) {
	f.ip = ip
}

// BasePointer retrieves the base pointer value of the current frame, indicating the starting position of its variables.
func (f *Frame) BasePointer() int {
	return f.basePointer
}

// SetBasePointer sets the base pointer index for the frame, which is used to manage the call stack during execution.
func (f *Frame) SetBasePointer(basePointer int) {
	f.basePointer = basePointer
}

// SetCompiledFunction sets the compiled function for the current frame.
func (f *Frame) SetCompiledFunction(compiledFunction *objects.FunctionCompiled) {
	f.compiledFunction = compiledFunction
}

// Instructions retrieve the bytecode instructions associated with the current frame's compiled function.
func (f *Frame) Instructions() *objects.Instructions {
	return f.compiledFunction.Instructions()
}

// SourcePos returns the source position mapped to the given instruction pointer (ip) in the current frame's context.
func (f *Frame) SourcePos(ip int) int {
	return f.compiledFunction.SourcePos(ip)
}

// SameFunction compares the given FunctionCompiled with the Frame's compiled function and returns true if they are the same.
func (f *Frame) SameFunction(callee *objects.FunctionCompiled) bool {
	return callee == f.compiledFunction
}
