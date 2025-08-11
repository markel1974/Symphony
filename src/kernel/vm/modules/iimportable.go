package modules

// Importable represents an interface for importing a module by its name, returning the module instance or an error.
type Importable interface {
	Import(moduleName string) (interface{}, error)
}
