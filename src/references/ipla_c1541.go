package references

import "github.com/markel1974/c64emu/src/config"

const IdIPLAc1541 = "IPLAc1541"

// IPLAc1541 represents an interface for handling PLA logic in a 1541 disk drive emulation.
// Setup initializes the interface by linking it to VIA components, ROM loader, and configuration data.
// Read retrieves the value from the specified memory address.
// Write writes a value to the specified memory address.
type IPLAc1541 interface {
	Setup(via1 IVIA, via2 IVIA, roms IROMLoaderC1541, cfg *config.Config) error

	Read(addr uint16) uint8

	Write(addr uint16, data uint8)
}
