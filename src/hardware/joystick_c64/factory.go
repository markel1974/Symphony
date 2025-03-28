package joystick_c64

import (
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/registry"
)

// Identifier returns the unique identifier string for the joystick component.
func Identifier() string {
	return "joystick_c64"
}

// Factory is a type that implements component creation and identification functionalities.
type Factory struct {
}

// NewFactory initializes and returns a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{}
}

// Identifier returns a string that represents the unique identifier of the Factory.
func (t *Factory) Identifier() string {
	return Identifier()
}

// Create initializes a new joystick component using the specified parent, factory, and label parameters.
func (t *Factory) Create(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IComponent {
	return NewJoystick(parent, factory, label, instance)
}

func init() {
	registry.RegisterComponentFactory(NewFactory())
}
