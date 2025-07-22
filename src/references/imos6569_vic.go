package references

import (
	"fmt"
)

// IMos6569Socket defines the interface for VIC's connection socket, providing methods for synchronization and signal control.
// Cycle returns the current cycle count of the VIC socket.
// GetBanks retrieves the associated IVICBanks for memory access operations.
// IRQTrigger signals an interrupt request trigger on the VIC socket.
// IRQClearTrigger clears an active interrupt request signal on the VIC socket.
// BALow sets or unsets the state of the BA (Bus Available) line.
// AECLow sets or unsets the state of the AEC (Address Enable Control) line.
// VBlank notifies the VIC to enter vertical blanking interval.
// LastCycle is invoked to finalize operations at the last cycle of the frame.
type IMos6569Socket interface {
	Cycle() uint64

	IRQTrigger()

	IRQClearTrigger()

	BALow(d bool)

	AECLow(d bool)

	VBlank()

	LastCycle()

	ScreenFreq() int

	TotalRaster() int

	ReadRam(addr uint16) uint8

	ReadColorRam(addr uint16) uint8

	ReadCharRom(addr uint16) uint8
}

// IMos6569 defines an interface for emulating a VIC component with methods for setup, interaction, state management, and operations.
type IMos6569 interface {
	Setup() error

	Bind(socket IMos6569Socket) error

	Connect() error

	Reset()

	Emulate()

	GetBALow() bool

	GetAECLow() bool

	WriteRegister(addr uint16, data uint8)

	ReadRegister(addr uint16) uint8

	LightPenTrigger()

	GetText() []byte

	ChangedVA(newVA uint8)

	GetVASignal() uint8
}

// IdIMos6569 generates a unique identifier for an IMos6569 interface based on the given label, instance, and interface name.
func IdIMos6569(v IMos6569, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIMos6569 converts an IComponent to an IMos6569 if the provided component matches the IMos6569 interface type.
// Returns an IMos6569 instance and nil error if successful, or nil and an error if the conversion fails.
// An error is also returned if the input component is nil.
func ComponentToIMos6569(component IComponent) (IMos6569, error) {
	if component == nil {
		return nil, fmt.Errorf("component IMos6569 is nil")
	}
	v, ok := component.(IMos6569)
	if !ok {
		return nil, fmt.Errorf("component is not a IMos6569")
	}
	return v, nil
}
