package cartridges_c64

import (
	"github.com/markel1974/c64emu/src/references"
)

// Identifier returns a string that serves as a unique identifier for a specific component or entity.
func Identifier() string {
	return "cartridges_c64"
}

// Factory represents a construct responsible for creating and managing components in a hierarchical manner.
type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string representing the unique identifier of the Factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new IComponent instance using the provided parent, factory, and label parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewManager(parent, factory, label)
}
