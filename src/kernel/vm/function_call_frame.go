package vm

import (
	"github.com/markel1974/c64emu/src/kernel/compiler"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// FunctionCallFrame represents a frame in the execution call stack for a compiled function.
// It contains information about the current function, its free variables, instruction pointer, and base pointer.
type FunctionCallFrame struct {
	compiledFunction *objects.CompiledFunction
	freeVars         []*objects.ObjectPtr
	ip               int
	basePointer      int
}

// NewFunctionCallFrame creates and returns a new instance of FunctionCallFrame with default values.
func NewFunctionCallFrame() *FunctionCallFrame {
	return &FunctionCallFrame{}
}

// SetCompiledFunction sets the compiled function for the current FunctionCallFrame.
func (f *FunctionCallFrame) SetCompiledFunction(compiledFunction *objects.CompiledFunction) {
	f.compiledFunction = compiledFunction
}

// Instructions returns the bytecode instructions of the currently compiled function in the call frame.
func (f *FunctionCallFrame) Instructions() []byte {
	return f.compiledFunction.Instructions()
}

// SourcePos returns the source position of the instruction at the given instruction pointer (ip) in the call frame.
func (f *FunctionCallFrame) SourcePos(ip int) compiler.Pos {
	return f.compiledFunction.SourcePos(ip)
}

// SameFunction checks if the given compiled function matches the compiled function in the current call frame.
func (f *FunctionCallFrame) SameFunction(callee *objects.CompiledFunction) bool {
	return callee == f.compiledFunction
}
