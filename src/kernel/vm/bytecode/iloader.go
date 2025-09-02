package bytecode

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ILoader defines an interface for loading and resolving objects, including built-in functions and imports.
// Id returns the loader's unique identifier.
// AddPackage adds a package to the loader's package map.
// BuiltinLen returns the number of built-in objects available.
// Builtin retrieves a built-in object by its index.
// Resolve attempts to resolve a slice of objects into their corresponding loaded versions or returns an error.
type ILoader interface {
	Id() string

	AddPackage(id string, attr map[string]objects.IObject)

	Resolve([]objects.IObject) ([]objects.IObject, error)
}
