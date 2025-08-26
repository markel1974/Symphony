package references

import (
	"fmt"
)

// IC64PlaSocket represents an interface for socket communication with the PLA, facilitating component interaction and integration.
type IC64PlaSocket interface {
}

type IC64PlaVASignals interface {
	GetVASignal() uint8
}

type IC64PlaChipSelect interface {
	WriteRegister(uint16, uint8)

	ReadRegister(uint16) uint8
}

// IC64Pla defines the interface for the Programmable Logic Array (PLA) component of the C64 system.
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
// RemoveRamTrigger removes a previously set write trigger function for a specified 16-bit memory address identified by a trigger Id.
type IC64Pla interface {
	Setup() error

	Bind(socket IC64PlaSocket, va IC64PlaVASignals, cartMan IC64CartridgeManager, ram IC64Ram, roms IC64Roms, cs0 IC64PlaChipSelect, cs1 IC64PlaChipSelect, cs2 IC64PlaChipSelect, cs3 IC64PlaChipSelect, cs4 IC64PlaChipSelect) error

	Connect() error

	Reset()

	//GetMemoryConfig() []uint8

	//SetMemoryEntry(m uint8)

	//SetMemoryConfig(m []uint8)

	RebuildMemoryConfig()

	Write(addr uint16, data uint8)

	Read(addr uint16) uint8

	WriteExt(memoryConfig int, addr uint16, data uint8)

	ReadExt(memoryConfig int, addr uint16) uint8

	SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int

	RemoveRamTrigger(addr uint16, id int)
}

// IdIC64Pla generates a unique identifier for an IC64Pla interface instance using the specified label and instance index.
func IdIC64Pla(v IC64Pla, label string, instance int) string {
	return IdInternalComponent(label, instance, InterfaceName(&v))
}

// ComponentToIC64Pla converts an IComponent to an IC64Pla if possible, returning an error if the conversion fails.
func ComponentToIC64Pla(component IComponent) (IC64Pla, error) {
	if component == nil {
		return nil, fmt.Errorf("component IC64Pla is nil")
	}
	v, ok := component.(IC64Pla)
	if !ok {
		return nil, fmt.Errorf("component is not a IC64Pla")
	}
	return v, nil
}
