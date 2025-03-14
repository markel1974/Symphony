package references

import (
	"github.com/markel1974/c64emu/src/config"
)

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
