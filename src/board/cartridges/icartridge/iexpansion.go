package icartridge

import "github.com/markel1974/c64emu/src/board/quartz"

// https://www.c64-wiki.com/wiki/Expansion_Port

type IExpansion interface {
	GameExRomConfigChanged()
	RamRead(uint16) uint8
	RamWrite(uint16, uint8)
	NMITrigger()
	DMALow(bool)
	BusAvailable() bool
	IRQTrigger()
	IRQClear()
	HasIRQ()
	ResetTrigger()
	GetIrqCycleDistance(int) uint64
	Cycle() uint64
	CreateAlarm(string, quartz.AlarmCallback) *quartz.Alarm
	RmwFlags() uint8 //TODO NOT STANDARD
}
