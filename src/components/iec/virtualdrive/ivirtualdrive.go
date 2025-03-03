package virtualdrive

import "github.com/markel1974/c64emu/src/config"

// IIec defines an interface for reading from and writing to peripherals using device-specific implementations.
// It provides methods for reading a byte and writing a byte to a designated device.
type IIec interface {
	PeripheralRead() uint8
	PeripheralWrite(deviceNumber uint8, data uint8)
}

// IVirtualDrive represents the interface for a virtual drive in the emulation environment.
// Setup initializes the virtual drive with the given configuration.
// Reset resets the state of the virtual drive.
// Emulate executes the emulation cycle for the virtual drive.
// Ready checks if the virtual drive is ready for operation.
// GetDeviceNumber returns the device number associated with the virtual drive.
// AtnStateChanged handles changes in the Attention (ATN) line state.
// BusStateChanged handles changes in the CPU bus state.
type IVirtualDrive interface {
	Setup(*config.Config)
	Reset()
	Emulate()
	Ready() bool
	GetDeviceNumber() uint8
	AtnStateChanged(bool, bool)
	BusStateChanged(uint8)

	//GetPath() string
	//Open(uint8, []uint8) uint8
	//Close(uint8) uint8
	//Read(uint8) (uint8, uint8)
	//Write(uint8, uint8, bool) uint8
}
