package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// KeyboardSocket represents a concrete implementation of IC64Keyboard, facilitating keyboard interactions in a system.
type KeyboardSocket struct {
	references.IC64Keyboard
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewKeyboardSocket initializes and returns a new instance of KeyboardSocket with IC64Keyboard set to nil.
func NewKeyboardSocket(parent references.IComponent, label string) *KeyboardSocket {
	c := &KeyboardSocket{
		IC64Keyboard: nil,
		parent:       parent,
		label:        label,
	}
	c.hwId = references.IdIC64Keyboard(c.IC64Keyboard, c.label, 0)
	return c
}

func (w *KeyboardSocket) HardwareId() string {
	return w.hwId
}

// Wire initializes the KeyboardSocket by resolving and setting IC64Keyboard and invoking its Setup method with provided config.
func (w *KeyboardSocket) Wire() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IC64Keyboard, err = references.ComponentToIC64Keyboard(w.component); err != nil {
		return err
	}
	if err = w.IC64Keyboard.Bind(w); err != nil {
		return err
	}
	return nil
}
