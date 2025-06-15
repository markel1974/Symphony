package references

import (
	"fmt"
)

// IC64KeyboardSocket represents an interface that allows binding with a keyboard for input handling and operations.
type IC64KeyboardSocket interface {
}

// IC64Keyboard defines an interface for keyboard functionality, including setup, binding, connection, and key state management.
type IC64Keyboard interface {
	Setup() error

	Bind(socket IC64KeyboardSocket) error

	Connect() error

	Reset()

	Emulate()

	NumLockToggle()

	CapitalToggle()

	SetKey(pressed bool, vKey int)

	Poll() (uint32, bool)

	SetCommand(cmd string)
}

// IdIC64Keyboard generates a unique identifier for an IC64Keyboard interface by combining the label, instance, and interface name.
func IdIC64Keyboard(v IC64Keyboard, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Keyboard converts an IComponent implementation to an IC64Keyboard interface or returns an error if invalid.
func ComponentToIC64Keyboard(component IComponent) (IC64Keyboard, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Keyboard is nil")
	}
	v, ok := component.(IC64Keyboard)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Keyboard")
	}
	return v, nil
}
