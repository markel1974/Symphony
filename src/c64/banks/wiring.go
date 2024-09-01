package banks

import "github.com/markel1974/c64emu/src/c64/cartridges/icartridge"

type ISystemWiring interface {
	WriteRegister(addr uint16, data uint8)
	ReadRegister(addr uint16) uint8
	GetLastByte() uint8
}

type IExpansionWiring interface {
	Config() (uint8, uint8, bool)
	Read(interval icartridge.RomInterval, addr uint16) (uint8, bool)
	IORead(addr uint16) (uint8, bool)
	IOWrite(addr uint16, data uint8) bool
}
