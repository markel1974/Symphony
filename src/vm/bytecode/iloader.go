package bytecode

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

// IPackage represents a package interface that provides access to package-level functionality or resources.
// Name retrieves the name of the package.
// Get retrieves a specific object from the package by its name, returning the object and a boolean indicating success.
type IPackage interface {
	Name() string

	Get(name string) (objects.IObject, bool)
}

// ILoader is an interface for managing and resolving packages and objects within a runtime or compilation environment.
// AddPackage adds a package, represented by IPackage, to the loader and returns an error if the operation fails.
// Resolve takes a slice of IObject references and resolves them, returning the resolved objects or an error.
type ILoader interface {
	AddPackage(pkg IPackage) error

	AddPackages(packages []IPackage) error

	Resolve([]objects.IObject) ([]objects.IObject, error)
}
