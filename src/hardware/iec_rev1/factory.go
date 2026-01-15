package iec_rev1

import (
	"github.com/markel1974/symphony/src/references"
	"github.com/markel1974/symphony/src/registry"
)

// Identifier returns a constant string "iec", typically used as a unique identifier for components or factories.
func Identifier() string {
	return "iec"
}

// Factory represents a construct used to create and manage component instances and their identifiers.
type Factory struct {
}

// NewFactory creates and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string constant that uniquely identifies the Factory type or its associated functionality.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes and returns a new Dispatcher instance while registering it with the specified parent and factory.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewDispatcher(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
