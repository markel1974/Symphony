package mos6510

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the string "mos6510", representing a unique identifier.
func Identifier() string {
	return "mos6510"
}

// Factory is a structure responsible for creating and managing hierarchical components for a system.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that uniquely identifies the component type "mos6510".
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create instantiates and returns a new IComponent of type CPU with the specified parent, factory, and label.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewCPU(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
