package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// References represents a structure for managing object references in a container, utilizing a factory and error signaling.
type References struct {
	factory   objects.IGateKeeper
	container []objects.IObject
	errSignal func(err error)
}

// NewReferences creates and initializes a References instance with the provided IGateKeeper factory and error signaling function.
func NewReferences(factory objects.IGateKeeper, errSignal func(err error)) *References {
	return &References{
		factory:   factory,
		container: nil,
		errSignal: errSignal,
	}
}

// Setup replaces the current container with the provided list of references.
func (g *References) Setup(loader bytecode.ILoader, references []objects.IObject) error {
	var err error
	g.container, err = loader.Resolve(references)
	if err != nil {
		return err
	}
	return err
}

// Get retrieves the `IObject` at the specified index from the container or returns an undefined value if the index is invalid.
func (g *References) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid reference index: %d", index))
		return g.factory.UndefinedValue()
	}
	return g.container[index]
}
