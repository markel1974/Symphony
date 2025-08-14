package modules

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// IModule defines an interface for modules that can be imported and queried for named symbols.
// Import loads a module by its name and returns a generic interface or an error if loading fails.
// Symbol retrieves a named symbol from the module and indicates if it exists.
type IModule interface {
	Import(moduleName string) (interface{}, error)

	Symbol(name string) (objects.IObject, bool)
}
