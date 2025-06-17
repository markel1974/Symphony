package mos6510_pic_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a static string value representing the identifier "mos6510_pic".
func Identifier() string {
	return "mos6510_pic"
}

// Factory represents an entity responsible for creating and managing components in the system.
type Factory struct {
}

// NewFactory creates and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a unique string identifier for the Factory instance.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create creates a new instance of IComponent by initializing a Pic object with the specified parent, factory, and instance.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewPIC(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
