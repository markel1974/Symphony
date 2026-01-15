package board

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

// componentId is a constant that holds the unique identifier for a component.
const componentId = "c1541"

// Identifier returns the fixed identifier string for the component.
func Identifier() string {
	return componentId
}

// Factory is a type responsible for creating and managing components in a system.
type Factory struct {
}

// NewFactory initializes and returns a new instance of the Factory struct.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns the unique identifier for the component.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes a new Board component using the given parent, factory, and label, and returns it as IComponent.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewBoard(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
