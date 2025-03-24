package board

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

// KeyboardSocket is a struct that embeds the IKeyboard interface for managing keyboard input operations.
// It acts as an adapter to connect and interact with an external keyboard functionality.
type KeyboardSocket struct {
	references.IKeyboard
}

// NewKeyboardSocket creates and returns a new instance of KeyboardSocket with an uninitialized IKeyboard reference.
func NewKeyboardSocket() *KeyboardSocket {
	c := &KeyboardSocket{
		IKeyboard: nil,
	}
	return c
}

// Setup assigns the provided IKeyboard interface implementation to the KeyboardSocket and initializes it.
func (w *KeyboardSocket) Setup(cc map[string]references.IComponent, _ *config.Config) error {
	var err error
	if w.IKeyboard, err = references.ComponentsToIKeyboard(cc, 0); err != nil {
		return err
	}
	if err = w.IKeyboard.Setup(); err != nil {
		return err
	}
	return nil
}

func (w *KeyboardSocket) Connect() error {
	return nil
}
