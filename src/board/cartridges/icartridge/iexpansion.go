package icartridge

import (
	"github.com/markel1974/c64emu/src/board/quartz"
)

// https://www.c64-wiki.com/wiki/Expansion_Port

type IExpansion interface {
	GameExRomConfigChanged()
	Read(uint16) uint8
	Write(uint16, uint8)
	RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int
	RamRemoveWriteTrigger(addr uint16, id int)
	NMITrigger()
	DMALow(bool)
	BusAvailable() bool
	IRQTrigger()
	IRQClear()
	//HasIRQ() bool
	//IRQLine() uint32
	ResetTrigger()
	GetQuartz() *quartz.Quartz
	CreateAlarm(string, quartz.AlarmCallback) *quartz.Alarm
	RmwFlags() uint8 //TODO NOT STANDARD
}
