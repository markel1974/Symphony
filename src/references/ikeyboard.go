package references

import (
	"fmt"
)

// IKeyboardSocket represents an interface that allows binding with a keyboard for input handling and operations.
type IKeyboardSocket interface {
}

// IKeyboard defines an interface for keyboard functionality, including setup, binding, connection, and key state management.
type IKeyboard interface {
	Setup() error

	Bind(socket IKeyboardSocket) error

	Connect() error

	Reset()

	Emulate()

	NumLockToggle()

	CapitalToggle()

	SetKey(pressed bool, vKey int)

	Poll() (uint32, bool)

	SetCommand(cmd string)
}

// IdIKeyboard generates a unique identifier for an IKeyboard interface by combining the label, instance, and interface name.
func IdIKeyboard(v IKeyboard, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIKeyboard converts an IComponent implementation to an IKeyboard interface or returns an error if invalid.
func ComponentToIKeyboard(component IComponent) (IKeyboard, error) {
	if component == nil {
		return nil, fmt.Errorf("component IKeyboard is nil")
	}
	v, ok := component.(IKeyboard)
	if !ok {
		return nil, fmt.Errorf("component is not a IKeyboard")
	}
	return v, nil
}
