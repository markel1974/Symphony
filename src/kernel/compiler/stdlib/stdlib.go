package stdlib

import (
	"github.com/markel1974/c64emu/src/kernel/compiler/modules"
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Loader provides access to and management of modules, including built-in and source modules, for symbol resolution and imports.
type Loader struct {
	mod *modules.Modules
}

// NewLoader initializes and returns a new Loader instance with built-in modules preloaded.
func NewLoader() *Loader {
	mods := modules.NewModules()
	for name, mod := range BuiltinModules {
		mods.AddBuiltinModule(name, mod)
	}
	return &Loader{
		mod: mods,
	}
}

// GetBuiltin retrieves a FunctionBuiltin instance from the predefined list of built-in functions based on the given index.
func (l *Loader) GetBuiltin(idx int) *objects.FunctionBuiltin {
	return GetBuiltin(idx)
}

// GetSymbolFromDefinition retrieves a symbol from a module using the provided definition as an identifier.
// Returns the located symbol and a boolean indicating whether the symbol was successfully found.
func (l *Loader) GetSymbolFromDefinition(definition objects.IObject) (objects.IObject, bool) {
	return l.mod.GetSymbolFromDefinition(definition)
}
