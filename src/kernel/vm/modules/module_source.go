package modules

// ModuleSource represents a source-based module with its source code stored as a byte slice.
type ModuleSource struct {
	src []byte
}

// NewSourceModule creates and returns a new ModuleSource object initialized with the provided source code.
func NewSourceModule(src []byte) *ModuleSource {
	return &ModuleSource{src: src}
}

// Import loads the content from the ModuleSource and returns it as an interface, along with any potential error.
func (m *ModuleSource) Import(_ string) (interface{}, error) {
	return m.src, nil
}
