package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// JoystickSocket is a struct that embeds references.IJoystick to facilitate joystick operations and interactions.
type JoystickSocket struct {
	references.IJoystick
	instance int
}

// NewJoystickSocket creates and returns a new instance of JoystickSocket with no IJoystick implementation assigned.
func NewJoystickSocket(instance int) *JoystickSocket {
	c := &JoystickSocket{
		IJoystick: nil,
		instance:  instance,
	}
	return c
}

// Setup assigns the given IJoystick object to the JoystickSocket and prepares it for further operations.
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

func (w *JoystickSocket) Connect() error {
	if err := w.IJoystick.Connect(); err != nil {
		return err
	}
	return nil
}
