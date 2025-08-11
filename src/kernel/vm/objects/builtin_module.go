package objects

// BuiltinModule represents a module with predefined attributes that can be imported or accessed at runtime.
// Attrs stores a map of string keys to IObject values, representing the module's predefined attributes.
type BuiltinModule struct {
	Attrs map[string]IObject
}

// Import loads a built-in module by its name and returns its immutable map representation or an error if it fails.
func (m *BuiltinModule) Import(moduleName string) (interface{}, error) {
	return m.AsImmutableMap(moduleName), nil
}

// AsImmutableMap transforms the BuiltinModule's attributes into an immutable map, embedding the given module name.
func (m *BuiltinModule) AsImmutableMap(moduleName string) *ImmutableMap {
	attrs := make(map[string]IObject, len(m.Attrs))
	for k, v := range m.Attrs {
		attrs[k] = v.Copy()
	}
	attrs["__module_name__"] = NewString(moduleName)
	return &ImmutableMap{Value: attrs}
}
