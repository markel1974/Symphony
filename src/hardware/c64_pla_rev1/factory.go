package c64_pla_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a fixed string identifier to represent the component type or functionality in the system.
func Identifier() string {
	return "pla_c64"
}

// Factory represents a component factory responsible for creating and managing hierarchical components.
type Factory struct {
}

// NewFactory creates a new instance of the Factory struct and returns a pointer to it.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that uniquely identifies the factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes a new PLA instance with the specified parent component, factory, and label, and returns it.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewPLA(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
