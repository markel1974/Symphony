package dynamic_throttle_rev1

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns a constant string "dynamic_throttle" used for uniquely identifying the component or factory.
func Identifier() string {
	return "dynamic_throttle"
}

// Factory is a type responsible for creating and managing instances of components in a hierarchical structure.
type Factory struct {
}

// NewFactory creates and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that uniquely identifies the dynamic throttle functionality.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create instantiates a new DynamicThrottle component using the provided parent, factory, and label parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewDynamicThrottle(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
