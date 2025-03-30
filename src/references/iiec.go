package references

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
)

const (
	IECAtnABit = uint8(0x10)
)

const (
	IECSidecarEnabled     = uint16(0x100)
	IECSidecarAtnAEnabled = uint16(0x200)
	IECSidecarDrop        = uint16(0x400)
	IECSidecarAtnABit     = uint16(uint16(IECAtnABit) << 8)
)

func IdIIec(_ IIec, label string, instance int) string {
	return IdInternalComponent(label, instance, "IIec")
}

type IIecSocket interface {
}

// IIec defines an interface for managing communication between a CPU and peripherals in an emulated environment.
// Setup initializes the IIec instance with provided quartz and configuration objects.
// AddPeripheral adds a new peripheral device with specified parameters, associating it with a unique device ID.
// RemovePeripheral removes a peripheral device using its unique device ID.
// Reset reinitializes the state of the IIec instance and all associated peripherals.
// Emulate performs the emulation cycle for the IIec and connected peripherals.
// CpuRead retrieves the data from the CPU bus.
// CpuWrite writes data to the CPU bus for transmission to peripherals.
// PeripheralRead retrieves data sent from peripherals to the CPU.
// PeripheralWrite writes data from the CPU to a specific peripheral identified by its device number.
// LEDSignal provides access to the LED signal, allowing observation of state changes via a SignalUint32 instance.
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

	LEDSignal() *signals.SignalUint32
}

func IdIIecDevice(_ IIecDevice, label string, instance int) string {
	return IdInternalComponent(label, instance, "IIecDevices")
}

type IIecDeviceSocket interface {
}

// IIecDevice represents the interface for a virtual drive in the emulation environment.
// Setup initializes the virtual drive with the given configuration.
// Reset resets the state of the virtual drive.
// Emulate executes the emulation cycle for the virtual drive.
// Ready checks if the virtual drive is ready for operation.
// GetDeviceNumber returns the device number associated with the virtual drive.
// AtnStateChanged handles changes in the Attention (ATN) line state.
// LEDSignal provides access to the LED signal, allowing observation of state changes via a SignalUint32 instance.
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

	LEDSignal() *signals.SignalUint32
}

func IdIIecProtocolDevice(_ IIecProtocolDevice, label string, instance int) string {
	return IdInternalComponent(label, instance, "IIecProtocolDevice")
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

func ComponentsToIEC(cc map[string]IComponent, label string, instance int) (IIec, error) {
	id := IdIIec(nil, label, instance)
	c, err := ComponentToIEC(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
