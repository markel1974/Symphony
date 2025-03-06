package board

import "fmt"

// Dumper represents a utility for serializing and deserializing data using defined property metadata.
type Dumper struct {
	container map[string]interface{}
	props     map[string]*PropertyInfo
}

// NewDumper initializes a Dumper instance with a provided map of properties for managing data storage and retrieval.
func NewDumper(props map[string]*PropertyInfo) *Dumper {
	return &Dumper{
		container: make(map[string]interface{}),
		props:     props,
	}
}

// Dump retrieves the property associated with the given key and populates the provided values with its contents.
// Returns an error if the property is not found or if the retrieval operation fails.
func (d *Dumper) Dump(key string, values ...interface{}) error {
	prop, ok := d.props[key]
	if !ok {
		return fmt.Errorf("property '%s' not found", key)
	}
	return prop.get(d.container, key, values)
}

// Restore sets the value of a specified property in the container using the provided key and values. Returns an error if the property is not found or the operation fails.
func (d *Dumper) Restore(key string, value ...interface{}) error {
	prop, ok := d.props[key]
	if !ok {
		return fmt.Errorf("property '%s' not found", key)
	}
	return prop.set(d.container, key, value)
}

// GetContainer returns the internal container map that stores key-value pairs for the Dumper instance.
func (d *Dumper) GetContainer() map[string]interface{} {
	return d.container
}
