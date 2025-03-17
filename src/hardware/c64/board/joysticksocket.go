package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// JoystickSocket is a struct that embeds references.IJoystick to facilitate joystick operations and interactions.
type JoystickSocket struct {
	references.IJoystick
}

// NewJoystickSocket creates and returns a new instance of JoystickSocket with no IJoystick implementation assigned.
func NewJoystickSocket() *JoystickSocket {
	c := &JoystickSocket{
		IJoystick: nil,
	}
	return c
}

// Connect assigns the given IJoystick object to the JoystickSocket and prepares it for further operations.
func (w *JoystickSocket) Connect(j references.IJoystick) error {
	w.IJoystick = j
	return nil
}
