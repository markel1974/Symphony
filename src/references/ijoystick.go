package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIJoystick(_ IJoystick, label string, instance int) string {
	return IdInternalComponent(label, instance, "IJoystick")
}

type IJoystickSocket interface {
}

// IJoystick defines an interface for joystick operations including updates, resets, emulation, movement, key setting, and polling.
// Update defines a method to adjust sensitivity and recalibrate with minimum and maximum bounds.
// Reset defines a method to reinitialize the joystick state to default settings.
// Emulate defines a method to simulate joystick behavior or states.
// Move provides a method to update the joystick position and button states.
// SetKey adjusts the joystick state based on key presses or releases with a specific joystick ID.
// Poll retrieves the next joystick state and its validity, indicating if data is available.
type IJoystick interface {
	Setup(cc map[string]IComponent, cfg *config.Config) error

	Bind(socket IJoystickSocket) error

	Connect() error

	Update(min uint16, max uint16, sensitivity uint16)

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll() (uint8, bool)
}

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

func ComponentsToIJoystick(cc map[string]IComponent, label string, instance int) (IJoystick, error) {
	id := IdIJoystick(nil, label, instance)
	c, err := ComponentToIJoystick(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
