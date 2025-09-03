package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Globals encapsulates shared system-wide data like objects, kind, and error signals for centralized management.
type Globals struct {
	gk           objects.IGateKeeper
	container    []objects.IObject
	preInitFuncs []*objects.FuncCompiled
	init         []*objects.FuncCompiled
	errSignal    func(err error)
}

// NewGlobals initializes and returns a new Globals instance with the provided factory, kind, and error signal handler.
func NewGlobals(gk objects.IGateKeeper, errSignal func(err error)) *Globals {
	return &Globals{
		gk:           gk,
		container:    []objects.IObject{},
		preInitFuncs: []*objects.FuncCompiled{},
		init:         []*objects.FuncCompiled{},
		errSignal:    errSignal,
	}
}

// Setup initializes the Globals instance with the provided constants list and returns a mapping of function names to indices.
func (g *Globals) Setup(constants []objects.IObject, preInit string, init string) (map[string]uint, error) {
	g.container = constants
	entryPoints := make(map[string]uint)
	for idx, global := range g.container {
		switch c := global.(type) {
		case *objects.FuncCompiled:
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

// PreInitFuncs returns a slice of pre-initialization compiled functions associated with the Globals instance.
func (g *Globals) PreInitFuncs() []*objects.FuncCompiled {
	return g.preInitFuncs
}

// InitFuncs returns the slice of compiled functions designated to run during the initialization phase.
func (g *Globals) InitFuncs() []*objects.FuncCompiled {
	return g.init
}

// Get retrieves the object at the specified index from the container.
// If the index is out of bounds, it triggers an error signal and returns an undefined value.
func (g *Globals) Get(index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid global index: %d", index))
		return g.gk.UndefinedValue()
	}
	return g.container[index]
}

// Set updates the value at the specified index in the container. Emits an error signal if the index is invalid.
func (g *Globals) Set(index uint, value objects.IObject) {
	if index > uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return
	}
	value.SetStatic()
	g.container[index] = value
}
