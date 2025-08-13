package modules

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Modules represents a collection of named modules implementing the IImportable interface.
type Modules struct {
	m map[string]IImportable
}

// NewModules creates and returns a new instance of Modules with an initialized internal map.
func NewModules() *Modules {
	return &Modules{
		m: make(map[string]IImportable),
	}
}

// Add registers a new module in the Modules map with the specified name and module instance.
func (m *Modules) Add(name string, module IImportable) {
	m.m[name] = module
}

// AddBuiltinModule registers a built-in module with the specified name and attributes in the modules map.
func (m *Modules) AddBuiltinModule(name string, attrs map[string]objects.IObject) {
	m.m[name] = &ModuleBuiltin{Attrs: attrs}
}

// AddSourceModule adds a new source-based module to the module map with the specified name and source code.
func (m *Modules) AddSourceModule(name string, src []byte) {
	m.m[name] = NewSourceModule(src)
}

// Remove deletes the module with the specified name from the Modules map.
func (m *Modules) Remove(name string) {
	delete(m.m, name)
}

// Get retrieves the module associated with the specified name and returns it as an IImportable instance.
func (m *Modules) Get(name string) IImportable {
	return m.m[name]
}

// GetBuiltinModule retrieves a built-in module by its name and returns it as a *ModuleBuiltin or nil if not found.
func (m *Modules) GetBuiltinModule(name string) *ModuleBuiltin {
	mod, _ := m.m[name].(*ModuleBuiltin)
	return mod
}

// GetSourceModule retrieves a source module by name, returning nil if the module does not exist or is not a ModuleSource.
func (m *Modules) GetSourceModule(name string) *ModuleSource {
	mod, _ := m.m[name].(*ModuleSource)
	return mod
}

// Copy creates and returns a deep copy of the Modules instance, duplicating its internal map of modules.
func (m *Modules) Copy() *Modules {
	c := &Modules{
		m: make(map[string]IImportable),
	}
	for name, mod := range m.m {
		c.m[name] = mod
	}
	return c
}

// Len returns the number of entries in the Modules map.
func (m *Modules) Len() int {
	return len(m.m)
}

// AddMap merges all modules from the provided Modules instance into the current instance, overwriting existing keys if necessary.
func (m *Modules) AddMap(o *Modules) {
	for name, mod := range o.m {
		m.m[name] = mod
	}
}
