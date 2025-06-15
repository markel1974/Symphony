package references

import (
	"fmt"
)

// IPlaC64Socket represents an interface for socket communication with the PLA, facilitating component interaction and integration.
type IPlaC64Socket interface {
}

// IPlaC64 defines the interface for the Programmable Logic Array (PLA) component of the C64 system.
// Setup initializes the PLA component for operation.
// Bind links the PLA to various system components, including socket, VIC, SID, CIAs, cartridge manager, RAM, and ROM loader.
// Connect establishes the necessary connections for the PLA to interact with other components.
// Reset reinitializes the PLA to its default state.
// GetMemoryConfig retrieves the current memory configuration as a slice of 8-bit values.
// SetMemoryEntry sets a specific memory configuration entry using an 8-bit value.
// SetMemoryConfig updates the entire memory configuration using a slice of 8-bit values.
// RebuildMemoryConfig reconstructs the memory configuration based on current system state.
// Write writes an 8-bit data value to a specified 16-bit memory address.
// Read reads an 8-bit data value from a specified 16-bit memory address.
// ReadCharRom reads an 8-bit character ROM value from a specified 16-bit memory address.
// ReadDirect reads an 8-bit value directly from a specified 16-bit memory address bypassing memory configuration.
// ReadColor retrieves an 8-bit color value from a specified 16-bit memory address.
// SetWriteTrigger assigns a trigger function that executes when writing to a specific 16-bit memory address.
// RemoveRamTrigger removes a previously set write trigger function for a specified 16-bit memory address identified by a trigger ID.
type IPlaC64 interface {
	Setup() error

	Bind(socket IPlaC64Socket, vic IVIC, sid ISID, cia1 ICIA, cia2 ICIA, cartMan ICartridgeManagerC64, ram IRamC64, colorRam IColorRamC64, roms IRomsC64) error

	Connect() error

	Reset()

	GetMemoryConfig() []uint8

	SetMemoryEntry(m uint8)

	SetMemoryConfig(m []uint8)

	RebuildMemoryConfig()

	Write(addr uint16, data uint8)

	Read(addr uint16) uint8

	SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int

	RemoveRamTrigger(addr uint16, id int)
}

// IdIPlaC64 generates a unique identifier for an IPlaC64 interface instance using the specified label and instance index.
func IdIPlaC64(v IPlaC64, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIPLAc64 converts an IComponent to an IPlaC64 if possible, returning an error if the conversion fails.
func ComponentToIPLAc64(component IComponent) (IPlaC64, error) {
	if component == nil {
		return nil, fmt.Errorf("component IPlaC64 is nil")
	}
	v, ok := component.(IPlaC64)
	if !ok {
		return nil, fmt.Errorf("component is not a IPlaC64")
	}
	return v, nil
}
