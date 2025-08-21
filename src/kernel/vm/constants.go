package vm

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Constants is a structure that manages global objects and handles error signaling through a callback function.
type Constants struct {
	factory   *objects.GateKeeper
	constants []objects.IObject
	errSignal func(err error)
}

// NewConstants initializes and returns a new Constants instance with provided global objects and error signaling function.
func NewConstants(factory *objects.GateKeeper, errSignal func(err error)) *Constants {
	return &Constants{
		factory:   factory,
		constants: nil,
		errSignal: errSignal,
	}
}

// SetConstants updates the constants pool with the provided values.
func (g *Constants) SetConstants(constants []objects.IObject) {
	g.constants = constants
}

// Get retrieves the object from the constants pool at the specified index. Returns UndefinedValue if the index is invalid.
func (g *Constants) Get(index uint) objects.IObject {
	if index >= uint(len(g.constants)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return g.factory.UndefinedValue()
	}
	return g.constants[index]
}

// Set updates the global variable at the specified index with the given value.
// Triggers an error signal if the index is invalid.
func (g *Constants) Set(index uint, value objects.IObject) {
	if index > uint(len(g.constants)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return
	}
	g.constants[index] = value
}
