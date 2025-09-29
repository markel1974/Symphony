package handler

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/bytecode"
	"github.com/markel1974/c64emu/src/vm/objects"
)

// Imports represent a structure for managing object imports in a container, utilizing a factory and error signaling.
type Imports struct {
	gk             objects.IGateKeeper
	container      []objects.IObject
	shutdownSignal func(err error)
}

// NewImports creates and initializes a Imports instance with the provided IGateKeeper factory and error signaling function.
func NewImports(gk objects.IGateKeeper, shutdownSignal func(err error)) *Imports {
	return &Imports{
		gk:             gk,
		container:      nil,
		shutdownSignal: shutdownSignal,
	}
}

// Setup replaces the current container with the provided list of imports.
func (g *Imports) Setup(loader bytecode.ILoader, references []objects.IObject) error {
	var err error
	g.container, err = loader.Resolve(references)
	if err != nil {
		return err
	}
	if len(g.container) != len(references) {
		return fmt.Errorf("invalid number of imports: %d", len(references))
	}
	//internals
	for idx, val := range g.container {
		switch val.(type) {
		case *objects.Undefined:
			internal := references[idx].AsString()
			callId, ok := objects.CallIdFromString(references[idx].AsString())
			if !ok {
				return fmt.Errorf("invalid internal reference: %s", internal)
			}
			g.container[idx] = g.gk.NewFuncInternal(objects.FrameStatic, callId)
		}
	}
	return err
}

// Get retrieves an object from the container at the specified index, validates the index, and returns a copied object.
func (g *Imports) Get(frame int, index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.shutdownSignal(fmt.Errorf("invalid reference index: %d", index))
		return g.gk.UndefinedValue()
	}
	obj := g.container[index]
	return obj.Copy(frame, objects.MaxDepth)
}
