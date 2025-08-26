package references

import (
	"fmt"
)

// IECAtnABit represents the bit mask for the ATN (Attention) line used in peripheral communication.
const (
	IECAtnABit = uint8(0x10)
)

// IECSidecarEnabled represents the flag for enabling the IEC sidecar functionality.
// IECSidecarAtnAEnabled represents the flag for enabling attention A in the IEC sidecar.
// IECSidecarDrop represents the flag for dropping the IEC sidecar functionality.
// IECSidecarAtnABit represents the shifted value of IECAtnABit for use in the IEC sidecar configuration.
const (
	IECSidecarEnabled     = uint16(0x100)
	IECSidecarAtnAEnabled = uint16(0x200)
	IECSidecarDrop        = uint16(0x400)
	IECSidecarAtnABit     = uint16(uint16(IECAtnABit) << 8)
)

// IIecSocket represents an interface for handling LED activity status for devices in an IEC communication setup.
// LedActivity manipulates the LED state for a specified device, toggling it on or off.
type IIecSocket interface {
	LedActivity(deviceNumber uint8, led bool)
}

// IIec defines an interface for managing IEC bus operations, including setup, peripheral control, and data communication.
type IIec interface {
	Setup() error

	Bind(socket IIecSocket) error

	Connect() error

	CreatePeripherals() error

	AddPeripheral(kind string, deviceId uint8) error

	RemovePeripheral(deviceId uint8)

	Reset()

	Emulate()

	CpuRead() uint8

	CpuWrite(data uint8)

	PeripheralRead() uint8

	PeripheralWrite(deviceNumber uint8, data uint16)

	LedActivity(deviceNumber uint8, led bool)
}

// IdIIecDevice generates a unique identifier for an IIecDevice by combining its label, instance, and a fixed Id string.
func IdIIecDevice(_ IIecDevice, label string, instance int) string {
	return IdInternalComponent(label, instance, "IIecDevices")
}

// IIecDeviceSocket represents an interface for an IEC device connection socket, providing a binding mechanism for devices.
type IIecDeviceSocket interface {
}

// IIecDevice defines the interface for an IEC device, providing methods for setup, binding, connection, and emulation.
type IIecDevice interface {
	Setup() error

	Bind(socket IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error

	Connect() error

	Reset()

	EmulationRequired() bool

	Emulate()

	Shutdown()

	Ready() bool

	GetDeviceNumber() uint8

	AtnStateChanged(bool)

	LedActivity(led bool)
}

// IdIIecProtocolDevice generates a unique identifier string for an IIecProtocolDevice using a label and instance number.
func IdIIecProtocolDevice(_ IIecProtocolDevice, label string, instance int) string {
	return IdInternalComponent(label, instance, "IIecProtocolDevice")
}

// IIecProtocolDevice defines an interface for handling IEC protocol device operations.
// Write sends data to a device at a given secondary address and returns a status or result code.
// Read retrieves data from a device at a given secondary address along with a status or reply code.
// EOI signals the End-Of-Identify (EOI) condition for the specified channel and returns a result code.
// Open initializes or prepares a device at the specified secondary address for communication.
// Close terminates communication with a device at the specified secondary address.
// Listen sets the device to listen mode at the specified secondary address.
// Unlisten disables the listen mode for the device at the specified secondary address.
// Talk sets the device to talk mode at the specified secondary address.
// Untalk disables the talk mode for the device at the specified secondary address.
type IIecProtocolDevice interface {
	Write(sec uint8, b uint8) uint8

	Read(sec uint8) (uint8, uint8)

	EOI(channel uint8) uint8

	Open(sec uint8) uint8

	Close(sec uint8) uint8

	Listen(sec uint8) uint8

	Unlisten(sec uint8) uint8

	Talk(sec uint8) uint8

	Untalk(sec uint8) uint8
}

// IdIIec generates a unique identifier for an IIec component using the provided label, instance, and interface name.
func IdIIec(v IIec, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIEC converts an IComponent to an IIec type, returning an error if the conversion is invalid or nil.
func ComponentToIEC(component IComponent) (IIec, error) {
	if component == nil {
		return nil, fmt.Errorf("component IIec is nil")
	}
	v, ok := component.(IIec)
	if !ok {
		return nil, fmt.Errorf("component is not a IIec")
	}
	return v, nil
}
