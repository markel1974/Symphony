package board

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a constant string that serves as a unique identifier for specific components or factories.
func Identifier() string {
	return "c64_board"
}

// Factory represents a struct primarily responsible for creating and managing components in the system.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a unique string identifier for the Factory type.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new instance of IComponent using the provided parent, factory, and label parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewBoard(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
