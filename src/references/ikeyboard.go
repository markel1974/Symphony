package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/config"
)

func IdIKeyboard(_ IKeyboard, label string, instance int) string {
	return IdInternalComponent(label, instance, "IKeyboard")
}

type IKeyboardSocket interface {
}

// IKeyboard defines an interface for managing keyboard input, including virtual key states and command processing.
// Reset reinitializes the keyboard state, clearing any prior configurations or input data.
// Emulate triggers an emulation process for the keyboard, typically used for virtualization behaviors.
// NumLockToggle toggles the state of the Num Lock key on the keyboard.
// CapitalToggle toggles the state of the Caps Lock key on the keyboard.
// SetVirtualKey processes a virtual key with a pressed state and associated virtual keycode.
// Poll retrieves the next key from the keyboard storage and indicates if a key is available.
// SetCommand processes and stores input commands based on their mapped key representations.
type IKeyboard interface {
	Setup(cc map[string]IComponent, _ *config.Config) error

	Bind(socket IKeyboardSocket) error

	Connect() error

	Reset()

	Emulate()

	NumLockToggle()

	CapitalToggle()

	SetKey(pressed bool, vKey int)

	Poll() (uint32, bool)

	SetCommand(cmd string)
}

func ComponentToIKeyboard(component IComponent) (IKeyboard, error) {
	if component == nil {
		return nil, fmt.Errorf("component IKeyboard is nil")
	}
	v, ok := component.(IKeyboard)
	if !ok {
		return nil, fmt.Errorf("component is not a IKeyboard")
	}
	return v, nil
}

func ComponentsToIKeyboard(cc map[string]IComponent, label string, instance int) (IKeyboard, error) {
	id := IdIKeyboard(nil, label, instance)
	c, err := ComponentToIKeyboard(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
