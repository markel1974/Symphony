package references

import "github.com/markel1974/c64emu/src/config"

// IRomLoaderC64 is an interface that provides methods to load various ROM sections, including Kernal, Basic, and Char ROMs.
// LoadKernal loads the Kernal ROM bytes and returns the data as a slice of bytes.
// LoadBasic loads the Basic ROM bytes and returns the data as a slice of bytes.
// LoadChar loads the Character ROM bytes and returns the data as a slice of bytes.
type IRomLoaderC64 interface {
	Setup(cfg *config.Config) error

	Reset()

	LoadKernal() []byte

	LoadBasic() []byte

	LoadChar() []byte
}
