package sdk

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// IPackage represents a package interface that provides access to package-level functionality or resources.
// Name retrieves the name of the package.
// Get retrieves a specific object from the package by its name, returning the object and a boolean indicating success.
type IPackage interface {
	Name() string

	Get(name string) (objects.IObject, bool)
}

// IBuiltin represents an interface to define built-in functions or objects within the system.
// Container returns a collection of IObject instances associated with the built-in entity.
type IBuiltin interface {
	Container() []objects.IObject
}

// BuiltinEntry is a struct designed to wrap a Builtin instance and its associated IObject for additional functionality.
type BuiltinEntry struct {
	builtin *objects.Builtin
	object  objects.IObject
}

func NewBuiltinAdapter(builtin *objects.Builtin, object objects.IObject) *BuiltinEntry {
	return &BuiltinEntry{
		builtin: builtin,
		object:  object,
	}
}

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
