package references

import "github.com/markel1974/c64emu/src/config"

// IPlaC1541 represents an interface for handling PLA logic in a 1541 disk drive emulation.
// Setup initializes the interface by linking it to VIA components, ROM loader, and configuration data.
// Read retrieves the value from the specified memory address.
// Write writes a value to the specified memory address.
type IPlaC1541 interface {
	Setup(via1 IVia, via2 IVia, roms IRomLoaderC1541, cfg *config.Config) error

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)
}
