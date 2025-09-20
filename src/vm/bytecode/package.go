package bytecode

import "github.com/markel1974/c64emu/src/vm/objects"

// Package represents a module with a name and a collection of named objects.
type Package struct {
	name      string
	container map[string]objects.IObject
}

// NewPackage creates and returns a new Package instance with the specified name and attribute mapping.
func NewPackage(name string, functions []objects.IObject, constants map[string]objects.IObject) *Package {
	pkg := &Package{
		name:      name,
		container: make(map[string]objects.IObject),
	}
	for _, obj := range functions {
		switch fn := obj.(type) {
		case *objects.Func:
			pkg.add(fn.Name(), fn)
		case *objects.FuncImport:
			pkg.add(fn.Name(), fn)
		case *objects.FuncJit:
			pkg.add(fn.Name(), fn)
		}
	}
	for id, c := range constants {
		pkg.add(id, c)
	}
	return pkg
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

func (p *Package) add(id string, obj objects.IObject) {
	p.container[id] = obj
}
