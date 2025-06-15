package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// JoystickSocket represents a joystick interface socket that integrates joystick functionality with a specific instance ID.
type JoystickSocket struct {
	references.IJoystick
	label     string
	parent    references.IComponent
	component references.IComponent
	instance  int
	hwId      string
}

// NewJoystickSocket creates and returns a new JoystickSocket instance with a specified joystick reference and instance number.
func NewJoystickSocket(parent references.IComponent, label string, instance int) *JoystickSocket {
	c := &JoystickSocket{
		IJoystick: nil,
		parent:    parent,
		label:     label,
		instance:  instance,
	}
	c.hwId = references.IdIJoystick(c.IJoystick, c.label, c.instance)
	return c
}

func (w *JoystickSocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the JoystickSocket by resolving its IJoystick component and calling its Setup method with configuration.
func (w *JoystickSocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IJoystick, err = references.ComponentToIJoystick(w.component); err != nil {
		return err
	}
	if err = w.IJoystick.Bind(w); err != nil {
		return err
	}
	return nil
}
