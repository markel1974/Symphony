package modules

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Modules represents a set of named modules. Use NewModuleMap to create a
// new module map.
type Modules struct {
	m map[string]Importable
}

// NewModuleMap creates a new module map.
func NewModuleMap() *Modules {
	return &Modules{
		m: make(map[string]Importable),
	}
}

// Add adds an import module.
func (m *Modules) Add(name string, module Importable) {
	m.m[name] = module
}

// AddBuiltinModule adds a builtin module.
func (m *Modules) AddBuiltinModule(name string, attrs map[string]objects.IObject) {
	m.m[name] = &BuiltinModule{Attrs: attrs}
}

// AddSourceModule adds a source module.
func (m *Modules) AddSourceModule(name string, src []byte) {
	m.m[name] = NewSourceModule(src)
}

// Remove removes a named module.
func (m *Modules) Remove(name string) {
	delete(m.m, name)
}

// Get returns an import module identified by name. It returns if the name is
// not found.
func (m *Modules) Get(name string) Importable {
	return m.m[name]
}

// GetBuiltinModule returns a builtin module identified by name. It returns
// if the name is not found or the module is not a builtin module.
func (m *Modules) GetBuiltinModule(name string) *BuiltinModule {
	mod, _ := m.m[name].(*BuiltinModule)
	return mod
}

// GetSourceModule returns a source module identified by name. It returns if
// the name is not found or the module is not a source module.
func (m *Modules) GetSourceModule(name string) *SourceModule {
	mod, _ := m.m[name].(*SourceModule)
	return mod
}

// Copy creates a copy of the module map.
func (m *Modules) Copy() *Modules {
	c := &Modules{
		m: make(map[string]Importable),
	}
	for name, mod := range m.m {
		c.m[name] = mod
	}
	return c
}

// Len returns the number of named modules.
func (m *Modules) Len() int {
	return len(m.m)
}

// AddMap adds named modules from another module map.
func (m *Modules) AddMap(o *Modules) {
	for name, mod := range o.m {
		m.m[name] = mod
	}
}
