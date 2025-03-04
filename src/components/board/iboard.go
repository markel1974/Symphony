package board

import (
	"github.com/markel1974/c64emu/src/config"
)

// IDisplayBuffer defines methods for interacting with a display buffer, allowing data manipulation at specified indices.
// Set sets a single byte of data at the given index.
// SetMulti8 sets a single byte of data and applies it across multiple relevant sections.
// Set8 sets an array of 8 bytes of data starting at the given index.
type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}

// IBoard defines a hardware abstraction interface for managing emulation, input/output events, and configurations.
type IBoard interface {
	Setup(db IDisplayBuffer, p IPlayer, cfg *config.Config) error

	Emulate() bool

	GetText() []byte

	GetScreenSize() (int, int)

	Joystick1Move(x uint, y uint, buttons uint)
	Joystick2Move(x uint, y uint, buttons uint)
	Joy1SetKey(pressed bool, vKey int)
	Joy2SetKey(pressed bool, vKey int)
	JoySwap(p bool)

	KeyboardPaste(pressed bool)
	KeyboardSetCommand(cmd string)
	KeyboardNumLockToggle()
	KeyboardCapitalToggle()
	KeyboardSetVirtualKey(pressed bool, vKey int)

	SetMouse(x uint8, y uint8)

	Throttle() IThrottling

	DiskChange()
}
