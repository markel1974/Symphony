package references

import (
	"fmt"
)

func IdIVIC(_ IVIC, label string, instance int) string {
	return IdInternalComponent(label, instance, "IVIC")
}

// IVIC defines an interface for a Video Interface Chip emulation, managing display rendering and register interactions.
// Setup configures the VIC with a socket and a configuration object.
// Reset reinitializes the VIC to its default state.
// Emulate performs a single cycle of VIC emulation.
// GetBALow returns the current state of the BA (Bus Available) line.
// GetAECLow returns the current state of the AEC (Address Enable) line.
// WriteRegister writes a value to a specified register address.
// ReadRegister reads a value from a specified register address.
// LightPenTrigger simulates the activation of the light pen input.
// GetText retrieves the text data associated with the VIC.
// ChangedVA notifies the VIC of a change in the video address.
// GetLastByte returns the last byte processed by the VIC.
type IVIC interface {
	Setup() error

	Bind(socket IVICSocket) error

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

	GetLastByte() uint8
}

// IVICBanks represents an interface providing methods to read data from various VIC memory regions, like character ROM and colors.
// ReadCharRom reads data from the character ROM memory mapped at a given address.
// ReadDirect reads data directly from memory at the specified address.
// ReadColor reads color information from memory at the specified address.
type IVICBanks interface {
	ReadCharRom(uint16) uint8

	ReadDirect(uint16) uint8

	ReadColor(uint16) uint8
}

// IVICSocket defines an interface for interacting with a VIC chip's socket, supporting display buffer and bank operations.
// Cycle returns the current cycle count.
// GetDisplayBuffer retrieves the associated display buffer.
// GetBanks retrieves the associated VIC bank interface.
// IRQTrigger triggers an interrupt request.
// IRQClear clears an interrupt request.
// BALow sets the BA (Bus Available) line low or high.
// AECLow sets the AEC (Address Enable) line low or high.
// VBlank signals the start of a vertical blanking interval.
// LastCycle performs logic associated with the last cycle in an operation.
type IVICSocket interface {
	Cycle() uint64

	GetBanks() IVICBanks

	IRQTrigger()

	IRQClear()

	BALow(d bool)

	AECLow(d bool)

	VBlank()

	LastCycle()
}

func ComponentToIVIC(component IComponent) (IVIC, error) {
	if component == nil {
		return nil, fmt.Errorf("component IVIC is nil")
	}
	v, ok := component.(IVIC)
	if !ok {
		return nil, fmt.Errorf("component is not a IVIC")
	}
	return v, nil
}
