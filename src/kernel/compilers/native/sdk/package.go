package sdk

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// Package represents a module with a name and a collection of named objects.
type Package struct {
	name      string
	container map[string]objects.IObject
}

// NewExternalPackage creates and returns a new Package instance with the specified name and attribute mapping.
func NewExternalPackage(name string, attr map[string]objects.IObject) *Package {
	return &Package{
		name:      name,
		container: attr,
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
