package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// JoystickSocket represents a joystick interface socket that integrates joystick functionality with a specific instance Id.
type JoystickSocket struct {
	references.IC64Joystick
	label     string
	parent    references.IComponent
	component references.IComponent
	instance  int
	hwId      string
}

// NewJoystickSocket creates and returns a new JoystickSocket instance with a specified joystick reference and instance number.
func NewJoystickSocket(parent references.IComponent, label string, instance int) *JoystickSocket {
	c := &JoystickSocket{
		IC64Joystick: nil,
		parent:       parent,
		label:        label,
		instance:     instance,
	}
	c.hwId = references.IdIC64Joystick(c.IC64Joystick, c.label, c.instance)
	return c
}

func (w *JoystickSocket) HardwareId() string {
	return w.hwId
}

// Wire initializes the JoystickSocket by resolving its IC64Joystick component and calling its Bind method with configuration.
func (w *JoystickSocket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IC64Joystick, err = references.ComponentToIC64Joystick(w.component); err != nil {
		return err
	}
	if err = w.IC64Joystick.Bind(w); err != nil {
		return err
	}
	return nil
}
