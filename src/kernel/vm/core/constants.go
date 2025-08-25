package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Constants is a structure that manages global objects and handles error signaling through a callback function.
type Constants struct {
	factory   objects.IGateKeeper
	container []objects.IObject
	errSignal func(err error)
}

// NewConstants initializes and returns a new Constants instance with provided global objects and error signaling function.
func NewConstants(factory objects.IGateKeeper, errSignal func(err error)) *Constants {
	return &Constants{
		factory:   factory,
		container: nil,
		errSignal: errSignal,
	}
}

// Setup updates the constants pool with the provided values.
func (g *Constants) Setup(constants []objects.IObject) error {
	g.container = constants
	return nil
}

// Get retrieves the object from the constants pool at the specified index. Returns UndefinedValue if the index is invalid.
func (g *Constants) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return g.factory.UndefinedValue()
	}
	return g.container[index]
}
