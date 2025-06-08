package mechanic

type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Create initializes and returns an IDisk instance based on the provided image data or an empty disk if the image is nil.
func (f *Factory) Create(kind string) IMechanic {
	if kind == "async" {
		return NewAsync()
	}
	return NewSync()
}
