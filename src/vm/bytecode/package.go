package bytecode

import "github.com/markel1974/symphony/src/vm/objects"

// Package represents a container for storing and managing named objects implementing the IObject interface.
type Package struct {
	name      string
	container map[string]objects.IObject
}

// NewPackage creates and returns a new instance of Package with the specified name and an initialized container map.
func NewPackage(name string) *Package {
	pkg := &Package{
		name:      name,
		container: make(map[string]objects.IObject),
	}
	return pkg
}

// Add associates the given id with the provided IObject in the package's container.
func (p *Package) Add(id string, obj objects.IObject) {
	p.container[id] = obj
}

// Name returns the name of the Package as a string.
func (p *Package) Name() string {
	return p.name
}

// Get retrieves an object by its ID from the package's container. It returns the object and a boolean indicating existence.
func (p *Package) Get(id string) (objects.IObject, bool) {
	v, ok := p.container[id]
	return v, ok
}
