package icartridge

import (
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

type Interval int

const (
	ROM_LO   = Interval(1)
	ROM_HI_1 = Interval(2)
	ROM_HI_2 = Interval(4)
)

const (
	CMODE_8KGAME  = 0
	CMODE_16KGAME = 1
	CMODE_RAM     = 2
	CMODE_ULTIMAX = 3
)

type ICartridge interface {
	Setup(board IExpansion, ldr *loader.CRTLoader) error
	Write(i Interval, addr uint16, data uint8) bool
	Read(i Interval, addr uint16) (uint8, bool)
	IORead(addr uint16) (uint8, bool)
	IOWrite(addr uint16, data uint8) bool
	GetExRom() uint8
	GetGame() uint8
	Detach() error
}
