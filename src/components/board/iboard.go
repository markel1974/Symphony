package board

import (
	"github.com/markel1974/c64emu/src/config"
)

// IDisplayBuffer is an interface for managing display buffer operations in graphical rendering systems.
// Set sets a single value at the specified index in the display buffer.
// SetMulti8 sets a multi-bit value at the specified index in the display buffer.
// Set8 sets an array of 8-bit values starting at the specified index in the display buffer.
type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}

// IPlayer is an interface for managing player-related operations in a game or multimedia context.
// GetCurrentPosition returns the current position of the player.
// Write writes audio or data buffer with specified parameters.
type IPlayer interface {
	GetCurrentPosition() int
	Write([]uint32, int, int)
}

// IBoard is an interface representing the functionalities of an emulated board system.
type IBoard interface {
	Setup(db IDisplayBuffer, p IPlayer, cfg *config.Config) error

	Emulate() bool

	GetText() []byte

	GetScreenSize() (int, int)

	GetFrameInterval() int

	Joy1SetKey(pressed bool, vKey int)
	Joy2SetKey(pressed bool, vKey int)
	JoySwap(p bool)

	KeyboardPaste(pressed bool)
	KeyboardSetCommand(cmd string)
	KeyboardNumLockToggle()
	KeyboardCapitalToggle()
	KeyboardSetVirtualKey(pressed bool, vKey int)

	SetMouse(x uint8, y uint8)

	DiskChange()
}
