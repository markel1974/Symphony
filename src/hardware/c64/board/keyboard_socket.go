package board

import (
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

// Connect assigns the provided IKeyboard interface implementation to the KeyboardSocket and initializes it.
func (w *KeyboardSocket) Connect(k references.IKeyboard) error {
	w.IKeyboard = k
	if err := w.IKeyboard.Setup(); err != nil {
		return err
	}
	return nil
}
