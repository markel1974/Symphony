package modules

import (
	"github.com/markel1974/c64emu/src/kernel/vm/objects"
)

// Modules represents a collection of modules identified by their names, allowing management and access to importable modules.
type Modules struct {
	m map[string]IModule
}

// NewModules creates and returns a new instance of Modules with an initialized map for storing module data.
func NewModules() *Modules {
	return &Modules{
		m: make(map[string]IModule),
	}
}

// Add inserts a module into the Modules map with the specified name as the key.
func (m *Modules) Add(name string, module IModule) {
	m.m[name] = module
}

// AddBuiltinModule registers a built-in module with the given name and its attributes in the modules map.
func (m *Modules) AddBuiltinModule(name string, attrs map[string]objects.IObject) {
	m.m[name] = &ModuleBuiltin{Attrs: attrs}
}

// AddSourceModule adds a new source module to the collection with the given name and source code.
func (m *Modules) AddSourceModule(name string, src []byte) {
	m.m[name] = NewSourceModule(src)
}

// Remove removes the module with the specified name from the Modules map. If the name does not exist, no action is taken.
func (m *Modules) Remove(name string) {
	delete(m.m, name)
}

// Get retrieves a module by its name from the Modules map and returns it as an IModule instance.
func (m *Modules) Get(name string) IModule {
	return m.m[name]
}

// GetSymbol returns the symbol with the given name from the specified module. It also indicates if the symbol was found.
func (m *Modules) GetSymbol(module string, symbol string) (objects.IObject, bool) {
	mod, ok := m.m[module]
	if !ok {
		return nil, false
	}
	sym, ok := mod.Symbol(symbol)
	return sym, ok
}

// GetSymbolFromDefinition retrieves a symbol from a module based on the definition array containing package and symbol names.
// Returns the found symbol and a boolean indicating success or failure.
func (m *Modules) GetSymbolFromDefinition(in objects.IObject) (objects.IObject, bool) {
	definition, ok := in.(*objects.Array)
	if !ok {
		return nil, false
	}
	pName, err := definition.Index(0)
	if err != nil {
		return nil, false
	}
	sName, err := definition.Index(1)
	if err != nil {
		return nil, false
	}
	packageName, ok := pName.(*objects.String)
	if !ok {
		return nil, false
	}
	symbolName, ok := sName.(*objects.String)
	if !ok {
		return nil, false
	}
	mod, ok := m.m[packageName.Value()]
	if !ok {
		return nil, false
	}
	sym, ok := mod.Symbol(symbolName.Value())
	return sym, ok
}

// GetBuiltinModule retrieves a built-in module by its name from the module collection and returns it as *ModuleBuiltin.
func (m *Modules) GetBuiltinModule(name string) *ModuleBuiltin {
	mod, _ := m.m[name].(*ModuleBuiltin)
	return mod
}

// GetSourceModule retrieves a ModuleSource instance by its name from the Modules map or returns nil if not found.
func (m *Modules) GetSourceModule(name string) *ModuleSource {
	mod, _ := m.m[name].(*ModuleSource)
	return mod
}

// Copy creates and returns a new Modules instance with a deep copy of the current module map.
func (m *Modules) Copy() *Modules {
	c := &Modules{
		m: make(map[string]IModule),
	}
	for name, mod := range m.m {
		c.m[name] = mod
	}
	return c
}

// Len returns the number of modules currently stored in the Modules instance.
func (m *Modules) Len() int {
	return len(m.m)
}

// AddMap merges the modules from another Modules instance into the current one, overriding any existing modules with the same name.
func (m *Modules) AddMap(o *Modules) {
	for name, mod := range o.m {
		m.m[name] = mod
	}
}
