package c1541_pla_rev1

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

// Identifier returns a static string "c1541_pla", representing the unique identifier of the component.
func Identifier() string {
	return "c1541_pla"
}

// Factory represents a generic structure for creating and managing components.
type Factory struct {
}

// NewFactory creates and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that represents the identifier of the PLA factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new PLA component with the specified parent, factory, and label parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewPLA(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
