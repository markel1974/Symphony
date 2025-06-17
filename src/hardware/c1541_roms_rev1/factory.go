package c1541_roms_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the unique identifier string "c1541_roms" for the associated component.
func Identifier() string {
	return "c1541_roms"
}

// Factory is a type responsible for creating and initializing components in the system.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that represents the identifier for the Factory instance.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create instantiates a new Roms component using the provided parent, factory, and label as parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewRoms(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
