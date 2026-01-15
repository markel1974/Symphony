package handler

import (
	"fmt"

	"github.com/markel1974/symphony/src/vm/objects"
)

// Globals encapsulates shared system-wide data like objects, kind, and error signals for centralized management.
type Globals struct {
	gk        objects.IGateKeeper
	container []objects.IObject
}

// NewGlobals initializes and returns a new Globals instance with the provided factory, kind, and error signal handler.
func NewGlobals(gk objects.IGateKeeper) *Globals {
	return &Globals{
		gk:        gk,
		container: []objects.IObject{},
	}
}

// Setup initializes the Globals instance with the provided constants list and returns a mapping of function names to indices.
func (g *Globals) Setup(constants []objects.IObject) error {
	g.container = constants
	return nil
}

// Get retrieves the object at the specified index from the container.
// If the index is out of bounds, it triggers an error signal and returns an undefined value.
func (g *Globals) Get(index uint) (objects.IObject, error) {
	if index >= uint(len(g.container)) {
		return g.gk.UndefinedValue(), fmt.Errorf("invalid global index: %d", index)
	}
	return g.container[index], nil
}

// Set updates the value at the specified index in the container. Emits an error signal if the index is invalid.
func (g *Globals) Set(index uint, value objects.IObject) error {
	if index >= uint(len(g.container)) {
		return fmt.Errorf("invalid global index: %d", index)
	}
	value.SetStatic()
	g.container[index] = value
	return nil
}
