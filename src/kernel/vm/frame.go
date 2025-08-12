package vm

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Frame represents a frame in the execution call stack for a compiled function.
// It contains information about the current function, its free variables, instruction pointer, and base pointer.
type Frame struct {
	compiledFunction *objects.FunctionCompiled
	freeVars         []*objects.ObjectPointer
	ip               int
	basePointer      int
}

// NewFunctionCallFrame creates and returns a new instance of Frame with default values.
func NewFunctionCallFrame() *Frame {
	return &Frame{}
}

// SetCompiledFunction sets the compiled function for the current Frame.
func (f *Frame) SetCompiledFunction(compiledFunction *objects.FunctionCompiled) {
	f.compiledFunction = compiledFunction
}

// Instructions returns the bytecode instructions of the currently compiled function in the call frame.
func (f *Frame) Instructions() *objects.Instructions {
	return f.compiledFunction.Instructions()
}

// SourcePos returns the source position of the instruction at the given instruction pointer (ip) in the call frame.
func (f *Frame) SourcePos(ip int) int {
	return f.compiledFunction.SourcePos(ip)
}

// SameFunction checks if the given compiled function matches the compiled function in the current call frame.
func (f *Frame) SameFunction(callee *objects.FunctionCompiled) bool {
	return callee == f.compiledFunction
}
