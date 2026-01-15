package mos6522

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

// Identifier returns the static identifier string for the component, typically used for registration and reference purposes.
func Identifier() string {
	return "mos6522"
}

// Factory provides methods and structures for creating and managing components in a system.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that uniquely identifies the Factory type or its associated component.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create instantiates a new IComponent using the provided parent component, factory, and label identifier.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewVIA(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
