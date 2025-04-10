package references

func IdIBoardC64(_ IBoard, label string, instance int) string {
	return IdInternalComponent(label, instance, "IBoardC64")
}

func IdIBoardC1541(_ IIecDevice, label string, instance int) string {
	return IdInternalComponent(label, instance, "IBoardC1541")
}

func IdIBoardVIC20(_ IBoard, label string, instance int) string {
	return IdInternalComponent(label, instance, "IBoardVIC20")
}

// IBoard represents the interface for controlling and managing a board-based system's input, output, and state.
// Setup initializes the board with the provided display buffer, player, and configuration settings.
// Emulate executes the emulation cycle and returns whether it should continue running.
// GetText retrieves the textual representation of the board's screen content.
// GetScreenSize returns the width and height of the board's screen in pixels.
// Joystick1Move updates the state of joystick 1, including axis positions and button states.
// Joystick2Move updates the state of joystick 2, including axis positions and button states.
// Joy1SetKey sets the key state for joystick 1 based on the given virtual key and pressed state.
// Joy2SetKey sets the key state for joystick 2 based on the given virtual key and pressed state.
// JoySwap toggles or sets the swapping of joystick 1 and 2 controls.
// KeyboardPaste simulates pasting input to the keyboard.
// KeyboardSetCommand sends a specific command string to the keyboard.
// KeyboardNumLockToggle toggles the state of the keyboard's Num Lock functionality.
// KeyboardCapitalToggle toggles the state of the keyboard's Caps Lock functionality.
// KeyboardSetKey updates the state of a specific virtual keyboard key.
// SetMouse updates the mouse position relative to the board's display.
// Throttle retrieves the throttle interface for managing execution rates.
// DiskChange triggers an event for changing the current disk in the board's system.
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

	KeyboardSetCommand(cmd string)
	KeyboardNumLockToggle()
	KeyboardCapitalToggle()
	KeyboardSetKey(pressed bool, vKey int)

	SetMouse(x uint8, y uint8)
}

type IBoardConnections interface {
	VBlank()

	LedActivity(deviceNumber uint8, led bool)
}
