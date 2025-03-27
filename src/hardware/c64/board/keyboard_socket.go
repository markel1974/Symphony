package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// KeyboardSocket represents a concrete implementation of IKeyboard, facilitating keyboard interactions in a system.
type KeyboardSocket struct {
	references.IKeyboard
}

// NewKeyboardSocket initializes and returns a new instance of KeyboardSocket with IKeyboard set to nil.
func NewKeyboardSocket() *KeyboardSocket {
	c := &KeyboardSocket{
		IKeyboard: nil,
	}
	return c
}

// Setup initializes the KeyboardSocket by resolving and setting IKeyboard and invoking its Setup method with provided config.
func (w *KeyboardSocket) Setup(cc map[string]references.IComponent, cfg *config.Config) error {
	var err error
	if w.IKeyboard, err = references.ComponentsToIKeyboard(cc, 0); err != nil {
		return err
	}
	if err = w.IKeyboard.Setup(w, cfg); err != nil {
		return err
	}
	return nil
}
