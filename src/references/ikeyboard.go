package references

// IKeyboard defines an interface for managing keyboard input, including virtual key states and command processing.
// Reset reinitializes the keyboard state, clearing any prior configurations or input data.
// Emulate triggers an emulation process for the keyboard, typically used for virtualization behaviors.
// NumLockToggle toggles the state of the Num Lock key on the keyboard.
// CapitalToggle toggles the state of the Caps Lock key on the keyboard.
// SetVirtualKey processes a virtual key with a pressed state and associated virtual keycode.
// Poll retrieves the next key from the keyboard storage and indicates if a key is available.
// SetCommand processes and stores input commands based on their mapped key representations.
type IKeyboard interface {
	Setup() error

	Reset()

	Emulate()

	NumLockToggle()

	CapitalToggle()

	SetVirtualKey(pressed bool, vKey int)

	Poll() (uint32, bool)

	SetCommand(cmd string)
}
