package bytecode

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ILoader defines an interface for loading and resolving objects, including built-in functions and references.
// BuiltinLen returns the number of built-in objects available.
// Builtin retrieves a built-in object by its index.
// Resolve attempts to resolve a slice of objects into their corresponding loaded versions or returns an error.
type ILoader interface {
	AddPackage(id string, attr map[string]objects.IObject)

	BuiltinLen() int

	Builtin(idx int) *objects.Builtin

	Resolve([]objects.IObject) ([]objects.IObject, error)
}
