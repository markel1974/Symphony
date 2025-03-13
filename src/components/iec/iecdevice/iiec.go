package iecdevice

// IIec defines an interface for reading from and writing to peripherals using device-specific implementations.
// It provides methods for reading a byte and writing a byte to a designated device.
type IIec interface {
	PeripheralRead() uint8
	PeripheralWrite(deviceNumber uint8, data uint8)
}
