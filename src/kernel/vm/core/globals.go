package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Globals encapsulates shared system-wide data like objects, kind, and error signals for centralized management.
type Globals struct {
	factory   objects.IGateKeeper
	container []objects.IObject
	errSignal func(err error)
}

// NewGlobals initializes and returns a new Globals instance with the provided factory, kind, and error signal handler.
func NewGlobals(factory objects.IGateKeeper, errSignal func(err error)) *Globals {
	return &Globals{
		factory:   factory,
		container: nil,
		errSignal: errSignal,
	}
}

// SetContainer replaces the current container with the provided list of constants.
func (g *Globals) SetContainer(constants []objects.IObject) {
	g.container = constants
}

// Get retrieves the object at the specified index from the container.
// If the index is out of bounds, it triggers an error signal and returns an undefined value.
func (g *Globals) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid global index: %d", index))
		return g.factory.UndefinedValue()
	}
	return g.container[index]
}

// Set updates the value at the specified index in the container. Emits an error signal if the index is invalid.
func (g *Globals) Set(index uint, value objects.IObject) {
	if index > uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return
	}
	g.container[index] = value
}
