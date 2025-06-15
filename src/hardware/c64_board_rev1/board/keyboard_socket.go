package board

import (
	"github.com/markel1974/c64emu/src/references"
)

// KeyboardSocket represents a concrete implementation of IKeyboard, facilitating keyboard interactions in a system.
type KeyboardSocket struct {
	references.IKeyboard
	label     string
	parent    references.IComponent
	component references.IComponent
	hwId      string
}

// NewKeyboardSocket initializes and returns a new instance of KeyboardSocket with IKeyboard set to nil.
func NewKeyboardSocket(parent references.IComponent, label string) *KeyboardSocket {
	c := &KeyboardSocket{
		IKeyboard: nil,
		parent:    parent,
		label:     label,
	}
	c.hwId = references.IdIKeyboard(c.IKeyboard, c.label, 0)
	return c
}

func (w *KeyboardSocket) HardwareId() string {
	return w.hwId
}

// Mount initializes the KeyboardSocket by resolving and setting IKeyboard and invoking its Setup method with provided config.
func (w *KeyboardSocket) Mount() error {
	var err error
	w.component = w.parent.GetChildByHardwareId(w.HardwareId())
	if w.IKeyboard, err = references.ComponentToIKeyboard(w.component); err != nil {
		return err
	}
	if err = w.IKeyboard.Bind(w); err != nil {
		return err
	}
	return nil
}
