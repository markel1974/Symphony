package pla_c1541

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a static string "pla_c1541", representing the unique identifier of the component.
func Identifier() string {
	return "pla_c1541"
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
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewPLA(parent, factory, label)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
