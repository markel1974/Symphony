package modules

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// IImportable defines an interface with a single method for importing a module by name, returning the module or an error.
type IImportable interface {
	Import(moduleName string) (interface{}, error)

	Symbol(name string) (objects.IObject, bool)
}

type IModuleGetter interface {
	Get(module string) IImportable

	GetSymbol(module string, symbol string) (objects.IObject, bool)

	GetSymbolFromDefinition(definition objects.IObject) (objects.IObject, bool)
}
