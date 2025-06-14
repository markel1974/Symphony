package references

import (
	"fmt"
)

// IJoystickSocket defines an interface for connecting a joystick device to external systems or components.
type IJoystickSocket interface {
}

// IJoystick defines an interface for managing joystick-related operations and interactions in a system.
type IJoystick interface {
	Setup() error

	Bind(socket IJoystickSocket) error

	Connect() error

	Update(min uint16, max uint16, sensitivity uint16)

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll() (uint8, bool)
}

// IdIJoystick generates a unique hardware identifier for an IJoystick instance using its label, instance number, and interface name.
func IdIJoystick(v IJoystick, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIJoystick converts an IComponent to an IJoystick interface if possible, returning an error on type mismatch or nil.
func ComponentToIJoystick(component IComponent) (IJoystick, error) {
	if component == nil {
		return nil, fmt.Errorf("component IJoystick is nil")
	}
	v, ok := component.(IJoystick)
	if !ok {
		return nil, fmt.Errorf("component is not a IJoystick")
	}
	return v, nil
}
