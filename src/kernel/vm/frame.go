package vm

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// frame represents a function call frame.
type frame struct {
	fn          *objects.CompiledFunction
	freeVars    []*objects.ObjectPtr
	ip          int
	basePointer int
}
