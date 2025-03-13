package iecdevice

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
	Write(sec uint8, b uint8)

	Read(sec uint8) uint8

	Open(sec uint8)

	Close(sec uint8)

	Listen(sec uint8)

	Unlisten(sec uint8)

	Talk(sec uint8)

	Untalk(sec uint8)
}
