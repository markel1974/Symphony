package quartz

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Factory provides methods to identify and create new component instances within a system. It acts as a component factory.
type Factory struct {
}

// Identifier returns the string "quartz", typically used as a unique name or identifier for a component or module.
func Identifier() string {
	return "quartz"
}

// NewFactory initializes and returns a pointer to a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns the unique string identifier for the Factory instance.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new Quartz instance using the specified parent component, factory, and label.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label int) references.IComponent {
	return NewQuartz(parent, factory, label)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
