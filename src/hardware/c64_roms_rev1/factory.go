package c64_roms_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a string identifier for the component, typically used for registration or reference purposes.
func Identifier() string {
	return "roms_c64"
}

// Factory represents a construct used to create and initialize components within a system.
type Factory struct {
}

// NewFactory creates and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string representing the unique identifier of the Factory instance.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new Roms component with the specified parent, factory, and label.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewRomLoader(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
