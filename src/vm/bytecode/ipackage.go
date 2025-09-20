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
