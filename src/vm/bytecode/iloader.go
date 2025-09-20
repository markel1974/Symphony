package bytecode

import (
	"github.com/markel1974/c64emu/src/vm/objects"
)

// ILoader defines an interface for loading and resolving objects, including built-in functions and imports.
// AddPackage adds a package to the loader's package map.
// Resolve attempts to resolve a slice of objects into their corresponding loaded versions or returns an error.
type ILoader interface {
	AddPackage(id string, functions []objects.IObject, constants map[string]objects.IObject)

	Resolve([]objects.IObject) ([]objects.IObject, error)
}
