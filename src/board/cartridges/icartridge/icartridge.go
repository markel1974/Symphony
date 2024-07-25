package icartridge

import (
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

type ICartridge interface {
	Setup(board IExpansion, ldr *loader.CRTLoader) error
	GetId() string
	Write(i RomInterval, addr uint16, data uint8) bool
	Read(i RomInterval, addr uint16) (uint8, bool)
	IORead(addr uint16) (uint8, bool)
	IOWrite(addr uint16, data uint8) bool
	GetExRom() uint8
	GetGame() uint8
	EmulationRequired() bool
	Emulate()
	Detach() error
}
