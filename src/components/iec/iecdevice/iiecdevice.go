package iecdevice

import "github.com/markel1974/c64emu/src/config"

// IIecDevice represents the interface for a virtual drive in the emulation environment.
// Setup initializes the virtual drive with the given configuration.
// Reset resets the state of the virtual drive.
// Emulate executes the emulation cycle for the virtual drive.
// Ready checks if the virtual drive is ready for operation.
// GetDeviceNumber returns the device number associated with the virtual drive.
// AtnStateChanged handles changes in the Attention (ATN) line state.
// BusStateChanged handles changes in the CPU bus state.
type IIecDevice interface {
	Setup(IIec, *config.Config)

	Reset()

	Emulate()

	Ready() bool

	GetDeviceNumber() uint8

	AtnStateChanged(bool)

	BusStateChanged(uint8)
}
