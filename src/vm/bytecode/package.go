package bytecode

import "github.com/markel1974/c64emu/src/vm/objects"

// RegisterPackageFn is a function type defining a method that registers a package using IGateKeeper and returns an IPackage.
type RegisterPackageFn func(f objects.IGateKeeper) IPackage

// Package represents a module with a name and a collection of named objects.
type Package struct {
	name      string
	functions []objects.IObject
	constants map[string]objects.IObject
	container map[string]objects.IObject
}

// NewPackage creates and returns a new Package instance with the specified name and attribute mapping.
func NewPackage(name string, functions []objects.IObject, constants map[string]objects.IObject) *Package {
	container := make(map[string]objects.IObject)
	for _, obj := range functions {
		switch fn := obj.(type) {
		case *objects.Func:
			container[fn.Name()] = fn
		case *objects.FuncImport:
			container[fn.Name()] = fn
		case *objects.FuncJit:
			container[fn.Name()] = fn
		}
	}
	for id, c := range constants {
		container[id] = c
	}
	return &Package{
		name:      name,
		functions: functions,
		constants: constants,
		container: container,
	}
}

// Name returns the name of the Package as a string.
func (p *Package) Name() string {
	return p.name
}

// Get retrieves an object from the package's container using the given id and returns the object and a boolean indicating success.
func (p *Package) Get(id string) (objects.IObject, bool) {
	v, ok := p.container[id]
	return v, ok
}
