package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Internals represents a container for managing IObject instances and interacting through an IGateKeeper interface.
// It includes a mechanism for signaling errors and is primarily used for internal execution context management.
type Internals struct {
	gk        objects.IGateKeeper
	container []objects.IObject
	errSignal func(err error)
}

// NewInternals initializes and returns a new instance of Internals with the provided IGateKeeper and error signaling function.
func NewInternals(gk objects.IGateKeeper, errSignal func(err error)) *Internals {
	return &Internals{
		gk:        gk,
		container: nil,
		errSignal: errSignal,
	}
}

// Setup initializes the internals by creating a new set of function internals with a static execution frame.
func (g *Internals) Setup() error {
	g.container = g.gk.NewFuncInternals(objects.FrameStatic)
	return nil
}

// Get retrieves an object from the internal container at the specified index or returns an undefined value for invalid indices.
func (g *Internals) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid reference index: %d", index))
		return g.gk.UndefinedValue()
	}
	return g.container[index]
}
