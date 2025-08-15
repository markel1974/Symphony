package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Globals is a structure that manages global objects and handles error signaling through a callback function.
type Globals struct {
	globals   []objects.IObject
	errSignal func(err error)
}

// NewGlobals initializes and returns a new Globals instance with provided global objects and error signaling function.
func NewGlobals(globals []objects.IObject, errSignal func(err error)) *Globals {
	return &Globals{
		globals:   globals,
		errSignal: errSignal,
	}
}

// Get retrieves the object from the globals pool at the specified index. Returns UndefinedValue if the index is invalid.
func (g *Globals) Get(index uint) objects.IObject {
	if index > uint(len(g.globals)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return objects.UndefinedValue
	}
	return g.globals[index]
}

// Set updates the global variable at the specified index with the given value.
// Triggers an error signal if the index is invalid.
func (g *Globals) Set(index uint, value objects.IObject) {
	if index > uint(len(g.globals)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return
	}
	g.globals[index] = value
}
