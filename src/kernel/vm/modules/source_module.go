package modules

// SourceModule represents a source-based module with its source code stored as a byte slice.
type SourceModule struct {
	src []byte
}

// NewSourceModule creates and returns a new SourceModule object initialized with the provided source code.
func NewSourceModule(src []byte) *SourceModule {
	return &SourceModule{src: src}
}

// Import loads the content from the SourceModule and returns it as an interface, along with any potential error.
func (m *SourceModule) Import(_ string) (interface{}, error) {
	return m.src, nil
}
