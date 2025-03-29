package references

import (
	"fmt"
)

func IdIPlaC64(_ IPlaC64, label string, instance int) string {
	return IdInternalComponent(label, instance, "IPlaC64")
}

type IPlaC64Socket interface {
}

// IPlaC64 defines the interface for a Programmable Logic Array implementation specifically for the C64 system.
// Setup initializes the PLA with required components like VIC, SID, CIAs, cartridge manager, ROM loader, and configuration.
// Reset reinitializes the state of the PLA to default values.
// GetMemoryConfig retrieves the current memory configuration array as a slice of uint8.
// SetMemoryEntry sets a single memory entry in the configuration to the specified value.
// SetMemoryConfig updates the entire memory configuration array with a new slice of uint8.
// RebuildMemoryConfig recalculates the memory configuration based on current settings and inputs.
// Write performs a write operation for a given 16-bit address and 8-bit data value.
// Read performs a read operation from a given 16-bit memory address, returning the corresponding 8-bit value.
// ReadCharRom retrieves data from the Character ROM at a specified 16-bit address.
// ReadDirect performs a direct memory read from the specified 16-bit address without any access filtering.
// ReadColor reads the color RAM value from a specified address in the range of the memory map.
// SetWriteTrigger associates a callback function to trigger on writes at the specified 16-bit address.
// RemoveRamTrigger removes a write trigger callback associated with a 16-bit address by its identifier.
type IPlaC64 interface {
	Setup() error

	Bind(socket IPlaC64Socket, vic IVIC, sid ISID, cia1 ICIA, cia2 ICIA, cartMan ICartridgeManagerC64, roms IROMLoaderC64) error

	Connect() error

	Reset()

	GetMemoryConfig() []uint8

	SetMemoryEntry(m uint8)

	SetMemoryConfig(m []uint8)

	RebuildMemoryConfig()

	Write(addr uint16, data uint8)

	Read(addr uint16) uint8

	ReadCharRom(addr uint16) uint8

	ReadDirect(addr uint16) uint8

	ReadColor(addr uint16) uint8

	SetWriteTrigger(addr uint16, fn func(uint16, uint8)) int

	RemoveRamTrigger(addr uint16, id int)
}

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

func ComponentsToIPLAc64(cc map[string]IComponent, label string, instance int) (IPlaC64, error) {
	id := IdIPlaC64(nil, label, instance)
	c, err := ComponentToIPLAc64(cc[id])
	if err != nil {
		return nil, err
	}
	return c, nil
}
