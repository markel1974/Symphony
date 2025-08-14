package modules

import "github.com/markel1974/c64emu/src/kernel/vm/objects"

// ModuleSource represents a source-based module with its source code stored as a byte slice.
type ModuleSource struct {
	src []byte
}

func (m *ModuleSource) Symbol(name string) (objects.IObject, bool) {
	return nil, false
}

// NewSourceModule creates and returns a new ModuleSource object initialized with the provided source code.
func NewSourceModule(src []byte) *ModuleSource {
	return &ModuleSource{src: src}
}

// Import loads the content from the ModuleSource and returns it as an interface, along with any potential error.
func (m *ModuleSource) Import(_ string) (interface{}, error) {
	return m.src, nil
}
