package handler

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Constants represents a structure that encapsulates gatekeeping, a container of objects, and functions for initialization.
type Constants struct {
	gk           objects.IGateKeeper
	container    []objects.IObject
	preInitFuncs []*objects.Func
	init         []*objects.Func
}

// NewConstants initializes and returns a new Constants object with the provided gatekeeper and empty configurations.
func NewConstants(gk objects.IGateKeeper) *Constants {
	return &Constants{
		gk:           gk,
		container:    nil,
		preInitFuncs: []*objects.Func{},
		init:         []*objects.Func{},
	}
}

// Setup initializes the Constants object with given constants and identifies preInit, init, and other entry points.
// It returns a map of entry point names and their indices, or an error if initialization fails.
func (g *Constants) Setup(constants []objects.IObject, preInit string, init string) (map[string]uint, error) {
	g.container = constants
	entryPoints := make(map[string]uint)
	for idx, global := range g.container {
		switch c := global.(type) {
		case *objects.Func:
			if c.Name() == preInit {
				g.preInitFuncs = append(g.preInitFuncs, c)
			} else if c.Name() == init {
				g.init = append(g.init, c)
			} else {
				entryPoints[c.Name()] = uint(idx)
			}
		}
	}
	return entryPoints, nil
}

// PreInitFuncs returns the list of functions marked as pre-initialization functions for the Constants instance.
func (g *Constants) PreInitFuncs() []*objects.Func {
	return g.preInitFuncs
}

// InitFuncs returns a slice of *objects.Func that were initialized during the setup phase.
func (g *Constants) InitFuncs() []*objects.Func {
	return g.init
}

// Retrieve fetches the object at the specified index in the constants container or returns an error if the index is invalid.
func (g *Constants) Retrieve(index uint) (objects.IObject, error) {
	if index >= uint(len(g.container)) {
		return g.gk.UndefinedValue(), fmt.Errorf("invalid constant index: %d", index)
	}
	return g.container[index], nil
}

// RetrieveFunc retrieves a function from the constants container at the specified index and returns an error if not found.
func (g *Constants) RetrieveFunc(index uint) (*objects.Func, error) {
	obj, err := g.Retrieve(index)
	if err != nil {
		return nil, err
	}
	entryFn, ok := obj.(*objects.Func)
	if !ok {
		return nil, fmt.Errorf("entry point not found: %d", index)
	}
	return entryFn, nil
}

// Get retrieves a constant at the specified index, creates a copy with the given frameId and predefined depth, and returns it.
func (g *Constants) Get(frameId int, index uint) (objects.IObject, error) {
	if index >= uint(len(g.container)) {
		return g.gk.UndefinedValue(), fmt.Errorf("invalid constant index: %d", index)
	}
	obj := g.container[index]
	return obj.Copy(frameId, objects.MaxDepth), nil
}
