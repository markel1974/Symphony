package core

import (
	"fmt"

	"github.com/markel1974/c64emu/src/vm/objects"
)

// Constants is a structure that manages global objects and handles error signaling through a callback function.
type Constants struct {
	gk           objects.IGateKeeper
	container    []objects.IObject
	preInitFuncs []*objects.Func
	init         []*objects.Func
	errSignal    func(err error)
}

// NewConstants initializes and returns a new Constants instance with provided global objects and error signaling function.
func NewConstants(gk objects.IGateKeeper, errSignal func(err error)) *Constants {
	return &Constants{
		gk:           gk,
		container:    nil,
		preInitFuncs: []*objects.Func{},
		init:         []*objects.Func{},
		errSignal:    errSignal,
	}
}

// Setup updates the constant pool with the provided values.
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

// PreInitFuncs returns a slice of pre-initialization compiled functions associated with the Globals instance.
func (g *Constants) PreInitFuncs() []*objects.Func {
	return g.preInitFuncs
}

// InitFuncs returns the slice of compiled functions designated to run during the initialization phase.
func (g *Constants) InitFuncs() []*objects.Func {
	return g.init
}

func (g *Constants) Retrieve(index uint) (objects.IObject, error) {
	if index >= uint(len(g.container)) {
		return g.gk.UndefinedValue(), fmt.Errorf("invalid constant index: %d", index)
	}
	return g.container[index], nil
}

// Get retrieves a constant object by its index and returns a copy adjusted for the specified frame and maximum depth.
func (g *Constants) Get(frameId int, index uint) objects.IObject {
	if index >= uint(len(g.container)) {
		g.errSignal(fmt.Errorf("invalid constant index: %d", index))
		return g.gk.UndefinedValue()
	}
	obj := g.container[index]
	return obj.Copy(frameId, objects.MaxDepth)
}
