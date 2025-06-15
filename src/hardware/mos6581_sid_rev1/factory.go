package mos6581

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the unique identifier string "mos6581".
func Identifier() string {
	return "mos6581"
}

// Factory is a struct implementing the creation and identification of components in a modular system.
type Factory struct {
}

// NewFactory initializes and returns a pointer to a new Factory instance.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that serves as the identifier for the Factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create instantiates a new SID component with a specified parent, factory, and label, initializing it for emulation.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewSID(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
