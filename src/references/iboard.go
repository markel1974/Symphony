package references

// IBoard defines the interface for managing and interacting with a board component in the system.
// Setup initializes the board component and returns an error if setup fails.
// Mount connects IBoard with IBoardConnections and returns an error if the process fails.
// Start begins the operation of the board and returns an error if starting fails.
// Reset resets the state of the board to its initial configuration.
// Emulate triggers the emulation process managed by the board.
// GetText retrieves the current text buffer from the board as a slice of bytes.
// Joystick1Move sends joystick 1 movement data with x, y coordinates and button states.
// Joystick2Move sends joystick 2 movement data with x, y coordinates and button states.
// Joy1SetKey sets the state of a joystick 1 key, specifying if it is pressed and the virtual key.
// Joy2SetKey sets the state of a joystick 2 key, specifying if it is pressed and the virtual key.
// JoySwap swaps the configurations of joystick 1 and joystick 2.
// HardwareButton processes hardware button inputs, specifying if pressed and its value.
// KeyboardSetCommand sends a keyboard command to the board.
// KeyboardNumLockToggle toggles the Num Lock state on the board's keyboard.
// KeyboardCapitalToggle toggles the Caps Lock state on the board's keyboard.
// KeyboardSetKey sets the state of a keyboard key, specifying if it is pressed and the virtual key.
// SetMouse updates the mouse position with x, y coordinates.
type IBoard interface {
	Setup() error

	Mount(conn IBoardConnections) error

	Start() error

	Reset()

	Emulate()

	GetText() []byte

	Joystick1Move(x uint, y uint, buttons uint)
	Joystick2Move(x uint, y uint, buttons uint)
	Joy1SetKey(pressed bool, vKey int)
	Joy2SetKey(pressed bool, vKey int)
	JoySwap()

	HardwareButton(pressed bool, val uint8)

	KeyboardSetCommand(cmd string)
	KeyboardNumLockToggle()
	KeyboardCapitalToggle()
	KeyboardSetKey(pressed bool, vKey int)

	SetMouse(x uint8, y uint8)
}

// IBoardConnections defines the interface for board connection interactions, such as VBlank execution and LED activity control.
// VBlank handles vertical blanking operations during rendering.
// LedActivity manages LED state for specified device numbers.
type IBoardConnections interface {
	VBlank()

	LedActivity(deviceNumber uint8, led bool)
}

// IdIBoardC64 generates a unique identifier for an IBoard interface using its label, instance, and interface name.
func IdIBoardC64(v IBoard, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

//func IdIBoardC1541(_ IIecDevice, label string, instance int) string {
//	return IdInternalComponent(label, instance, "IBoardC1541")
//}

//func IdIBoardVIC20(_ IBoard, label string, instance int) string {
//	return IdInternalComponent(label, instance, "IBoardVIC20")
//}
