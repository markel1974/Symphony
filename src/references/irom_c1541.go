package references

import "github.com/markel1974/c64emu/src/config"

// IROMLoaderC1541 is an interface for handling ROM loading functionality specific to the C1541 drive emulation.
// Setup configures the ROM loader using the provided configuration.
// Load retrieves the raw byte data of the ROM.
type IROMLoaderC1541 interface {
	Setup(cfg *config.Config) error

	Load() []byte
}
