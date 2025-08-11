package modules

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// BuiltinModule represents a module with predefined attributes that can be imported or accessed at runtime.
// Attrs stores a map of string keys to IObject values, representing the module's predefined attributes.
type BuiltinModule struct {
	Attrs map[string]objects.IObject
}

// Import loads a built-in module by its name and returns its immutable map representation or an error if it fails.
func (m *BuiltinModule) Import(moduleName string) (interface{}, error) {
	return m.AsImmutableMap(moduleName), nil
}

// AsImmutableMap transforms the BuiltinModule's attributes into an immutable map, embedding the given module name.
func (m *BuiltinModule) AsImmutableMap(moduleName string) *objects.ImmutableMap {
	attrs := make(map[string]objects.IObject, len(m.Attrs))
	for k, v := range m.Attrs {
		attrs[k] = v.Copy()
	}
	attrs["__module_name__"] = objects.NewString(moduleName)
	return &objects.ImmutableMap{Value: attrs}
}
