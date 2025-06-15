package references

import (
	"fmt"
)

// IC64JoystickSocket defines an interface for connecting a joystick device to external systems or components.
type IC64JoystickSocket interface {
}

// IC64Joystick defines an interface for managing joystick-related operations and interactions in a system.
type IC64Joystick interface {
	Setup() error

	Bind(socket IC64JoystickSocket) error

	Connect() error

	Update(min uint16, max uint16, sensitivity uint16)

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll() (uint8, bool)
}

// IdIC64Joystick generates a unique hardware identifier for an IC64Joystick instance using its label, instance number, and interface name.
func IdIC64Joystick(v IC64Joystick, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Joystick converts an IComponent to an IC64Joystick interface if possible, returning an error on type mismatch or nil.
func ComponentToIC64Joystick(component IComponent) (IC64Joystick, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Joystick is nil")
	}
	v, ok := component.(IC64Joystick)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Joystick")
	}
	return v, nil
}
