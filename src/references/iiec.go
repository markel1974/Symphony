package references

import "github.com/markel1974/c64emu/src/config"

// IIec defines an interface for reading from and writing to peripherals using device-specific implementations.
// It provides methods for reading a byte and writing a byte to a designated device.
type IIec interface {
	Setup(IQuartz, *config.Config) error

	Reset()

	Emulate()

	CpuRead() uint8

	CpuWrite(data uint8)

	PeripheralRead() uint8

	PeripheralWrite(deviceNumber uint8, data uint8)
}

// IIecDevice represents the interface for a virtual drive in the emulation environment.
// Setup initializes the virtual drive with the given configuration.
// Reset resets the state of the virtual drive.
// Emulate executes the emulation cycle for the virtual drive.
// Ready checks if the virtual drive is ready for operation.
// GetDeviceNumber returns the device number associated with the virtual drive.
// AtnStateChanged handles changes in the Attention (ATN) line state.
// BusStateChanged handles changes in the CPU bus state.
type IIecDevice interface {
	Setup(IIec, IQuartz, uint8, uint8, string, *config.Config) error

	Reset()

	Emulate()

	Ready() bool

	GetDeviceNumber() uint8

	AtnStateChanged(bool)

	BusStateChanged(uint8)
}

// IIecProtocolDevice defines an interface for interacting with devices using the IEC protocol.
// Write sends a byte to a specific secondary address on the device.
// Read retrieves a byte from a specific secondary address on the device.
// Open initiates communication with a specific secondary address on the device.
// Close terminates communication with a specific secondary address on the device.
// Listen signals the device to start listening on a specific secondary address.
// Unlisten signals the device to stop listening on a specific secondary address.
// Talk instructs the device to enter talking mode on a specific secondary address.
// Untalk instructs the device to exit talking mode on a specific secondary address.
type IIecProtocolDevice interface {
	Write(sec uint8, b uint8) uint8

	Read(sec uint8) (uint8, uint8)

	Open(sec uint8) uint8

	Close(sec uint8) uint8

	Listen(sec uint8)

	Unlisten(sec uint8)

	Talk(sec uint8)

	Untalk(sec uint8)
}
