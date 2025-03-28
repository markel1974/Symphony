package mos6526

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the identifier string for the component, which in this case is "mos6526".
func Identifier() string {
	return "mos6526"
}

// Factory is a struct designed for creating and handling components in a hierarchical and modular fashion.
type Factory struct {
}

// NewFactory creates and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string representing the unique identifier for the Factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new IComponent instance, associating it with the given parent, factory, and label.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewCIA(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
