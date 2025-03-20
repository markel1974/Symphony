package references

func IdIJoystick(_ IJoystick) string {
	return "IJoystick"
}

// IJoystick defines an interface for joystick operations including updates, resets, emulation, movement, key setting, and polling.
// Update defines a method to adjust sensitivity and recalibrate with minimum and maximum bounds.
// Reset defines a method to reinitialize the joystick state to default settings.
// Emulate defines a method to simulate joystick behavior or states.
// Move provides a method to update the joystick position and button states.
// SetKey adjusts the joystick state based on key presses or releases with a specific joystick ID.
// Poll retrieves the next joystick state and its validity, indicating if data is available.
type IJoystick interface {
	Setup() error

	Update(min uint16, max uint16, sensitivity uint16)

	Reset()

	Emulate()

	Move(x uint, y uint, buttons uint)

	SetKey(pressed bool, jId int)

	Poll() (uint8, bool)
}
