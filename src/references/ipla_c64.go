package references

import "github.com/markel1974/c64emu/src/config"

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
	Setup(vic IVIC, sid ISID, cia1 ICIA, cia2 ICIA, cartMan ICartridgeManagerC64, roms IROMLoaderC64, cfg *config.Config) error

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
