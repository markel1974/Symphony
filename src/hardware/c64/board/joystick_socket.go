package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// JoystickSocket represents a joystick interface socket that integrates joystick functionality with a specific instance ID.
type JoystickSocket struct {
	references.IJoystick
	instance int
}

// NewJoystickSocket creates and returns a new JoystickSocket instance with a specified joystick reference and instance number.
func NewJoystickSocket(instance int) *JoystickSocket {
	c := &JoystickSocket{
		IJoystick: nil,
		instance:  instance,
	}
	return c
}

// Setup initializes the JoystickSocket by resolving its IJoystick component and calling its Setup method with configuration.
func (w *JoystickSocket) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if w.IJoystick, err = references.ComponentsToIJoystick(cc, w.instance); err != nil {
		return err
	}
	if err = w.IJoystick.Setup(w, cfg); err != nil {
		return err
	}
	return nil
}
