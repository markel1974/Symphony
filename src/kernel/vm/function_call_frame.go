package vm

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// FunctionCallFrame represents a function call frame.
type FunctionCallFrame struct {
	compiledFunction *objects.CompiledFunction
	freeVars         []*objects.ObjectPtr
	ip               int
	basePointer      int
}

func NewFunctionCallFrame() *FunctionCallFrame {
	return &FunctionCallFrame{}
}

func (f *FunctionCallFrame) SetCompiledFunction(compiledFunction *objects.CompiledFunction) {
	f.compiledFunction = compiledFunction
}

func (f *FunctionCallFrame) Instructions() []byte {
	return f.compiledFunction.Instructions
}

func (f *FunctionCallFrame) SourcePos(ip int) objects.Pos {
	return f.compiledFunction.SourcePos(ip)
}

func (f *FunctionCallFrame) SameFunction(callee *objects.CompiledFunction) bool {
	return callee == f.compiledFunction
}
