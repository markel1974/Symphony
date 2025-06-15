package media_drive_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the string identifier for a specific component or factory, in this case, "media".
func Identifier() string {
	return "media"
}

// Factory provides a structure for creating and managing components, enabling the creation of hierarchical systems.
type Factory struct {
}

// NewFactory creates and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string representing the unique identifier for the factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new IComponent instance with the given parent, factory, label, and instance number.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewBoard(parent, factory, label, instance)
}

// init initializes the component factory by registering a new factory instance to the global registry.
func init() {
	registry.RegisterComponentFactory(NewFactory())
}
