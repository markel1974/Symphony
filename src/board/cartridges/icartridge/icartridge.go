package icartridge

import (
	"github.com/markel1974/c64emu/src/board/cartridges/loader"
)

type Interval int

//OFF
//GAME = 1, EXROM = 1

//8K Cartridge, $8000-$9FFF (ROML).
//GAME = 1, EXROM = 0
//ROML is read only. Basic ROM and Kernal ROM are available.

//16K Cartridge, $8000-$9FFF / $A000-$BFFF (ROML / ROMH).
//GAME = 0, EXROM = 0
//ROML/ROMH are read only, Basic ROM is overwritten by ROMH.

//16K Cartridge, $8000-$9FFF / $E000-$FFFF (ROML / ROMH). Ultimax mode.
//GAME = 0, EXROM = 1
//Ultimax mode is an emulation of the Japanese CBM machine called “MAX”. It is a predecessor of the C64 with less RAM. In Ultimax mode ROMH replaces the kernal at $E000. You do not need ROML for a cartridge to function and can be left out.

type ModeType int

const (
	ROM_LO   = Interval(1)
	ROM_HI_1 = Interval(2)
	ROM_HI_2 = Interval(4)
)

const (
	Mode16K     = 0
	Mode8K      = 1
	ModeUltimax = 2
	ModeOff     = 3
)

type Mode struct {
	Game         uint8
	ExRom        uint8
	IntervalLow  Interval
	IntervalHigh Interval
}

var Modes = []Mode{
	{Game: 0, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_1}, //16 K
	{Game: 0, ExRom: 1, IntervalLow: ROM_LO, IntervalHigh: 0},        //8 K
	{Game: 1, ExRom: 0, IntervalLow: ROM_LO, IntervalHigh: ROM_HI_2}, //ULTIMAX
	{Game: 1, ExRom: 1, IntervalLow: 0, IntervalHigh: 0},             //ALL RAM
}

type ICartridge interface {
	Setup(board IExpansion, ldr *loader.CRTLoader) error
	GetId() string
	Write(i Interval, addr uint16, data uint8) bool
	Read(i Interval, addr uint16) (uint8, bool)
	IORead(addr uint16) (uint8, bool)
	IOWrite(addr uint16, data uint8) bool
	GetExRom() uint8
	GetGame() uint8
	Detach() error
}
