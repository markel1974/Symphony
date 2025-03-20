package references

import "github.com/markel1974/c64emu/src/config"

func IdIIec(_ IIec) string {
	return "IIec"
}

// IIec defines the interface for interacting with Input/Output devices on the IEC (serial bus) in an emulation environment.
// Setup initializes the IEC instance with the provided quartz instance and configuration, returning an error if unsuccessful.
// Reset resets the state of the IEC interface to its initial default configuration.
// Emulate executes an IEC emulation cycle, processing communication between the CPU and peripherals.
// CpuRead handles read operations initiated by the CPU and returns the respective data byte.
// CpuWrite performs write operations from the CPU to the IEC interface using the provided data byte.
// PeripheralRead manages read operations performed by IEC peripherals and returns the relevant data byte.
// PeripheralWrite handles write operations to peripherals on the IEC bus using the specified device number and data byte.
type IIec interface {
	Setup(IQuartz, *config.Config) error

	Reset()

	Emulate()

	CpuRead() uint8

	CpuWrite(data uint8)

	PeripheralRead() uint8

	PeripheralWrite(deviceNumber uint8, data uint8)
}

func IdIIecDevice(_ IIecDevice) string {
	return "IIecDevices"
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

func IdIIecProtocolDevice(_ IIecProtocolDevice) string {
	return "IIecProtocolDevice"
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
