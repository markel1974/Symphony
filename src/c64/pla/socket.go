package pla

import "github.com/markel1974/c64emu/src/c64/cartridges/icartridge"

// IExpansionSocket defines an interface for managing expansion socket operations and interactions in a given system.
// Config retrieves the configuration data of the expansion socket and additional state details.
// Read fetches a value from a specified ROM interval and address, returning the value and success state.
// IORead performs a read operation from an I/O address, returning the value and success state.
// IOWrite executes a write operation to an I/O address with the specified data, indicating success.
type IExpansionSocket interface {
	Config() (uint8, uint8, bool)
	Read(interval icartridge.RomInterval, addr uint16) (uint8, bool)
	IORead(addr uint16) (uint8, bool)
	IOWrite(addr uint16, data uint8) bool
}

// ISocket represents an interface for reading and writing to hardware registers and retrieving the last accessed byte.
// WriteRegister writes an 8-bit value to a specified 16-bit address register.
// ReadRegister reads an 8-bit value from a specified 16-bit address register.
// GetLastByte retrieves the last byte that was accessed or operated upon.
type ISocket interface {
	WriteRegister(addr uint16, data uint8)
	ReadRegister(addr uint16) uint8
	GetLastByte() uint8
}

// IRomSocket is an interface that provides methods to load various ROM sections, including Kernal, Basic, and Char ROMs.
// LoadKernal loads the Kernal ROM bytes and returns the data as a slice of bytes.
// LoadBasic loads the Basic ROM bytes and returns the data as a slice of bytes.
// LoadChar loads the Character ROM bytes and returns the data as a slice of bytes.
type IRomSocket interface {
	LoadKernal() []byte
	LoadBasic() []byte
	LoadChar() []byte
}
