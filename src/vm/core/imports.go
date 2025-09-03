package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Imports represents a structure for managing object imports in a container, utilizing a factory and error signaling.
type Imports struct {
	gk        objects.IGateKeeper
	container []objects.IObject
	errSignal func(err error)
}

// NewImports creates and initializes a Imports instance with the provided IGateKeeper factory and error signaling function.
func NewImports(gk objects.IGateKeeper, errSignal func(err error)) *Imports {
	return &Imports{
		gk:        gk,
		container: nil,
		errSignal: errSignal,
	}
}

// Setup replaces the current container with the provided list of imports.
func (g *Imports) Setup(loader bytecode.ILoader, references []objects.IObject) error {
	var err error
	g.container, err = loader.Resolve(references)
	if err != nil {
		return err
	}
	return err
}

// Get retrieves the `IObject` at the specified index from the container or returns an undefined value if the index is invalid.
func (g *Imports) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid reference index: %d", index))
		return g.gk.UndefinedValue()
	}
	return g.container[index]
}
