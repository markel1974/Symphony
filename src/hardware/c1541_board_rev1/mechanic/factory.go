package mechanic

import "github.com/markel1974/symphony/src/references"

type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Create initializes and returns an IDisk instance based on the provided image data or an empty disk if the image is nil.
func (f *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int, kind string) IMechanic {
	if kind == "async" {
		return NewAsync(parent, factory, label, instance)
	}
	return NewSync(parent, factory, label, instance)
}
