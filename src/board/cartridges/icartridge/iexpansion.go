package icartridge

import (
	"github.com/markel1974/c64emu/src/components/quartz"
)

// https://www.c64-wiki.com/wiki/Expansion_Port

type IExpansion interface {
	Cycle() uint64
	CycleAlarm(string, quartz.AlarmCallback) *quartz.Alarm
	GameExRomConfigChanged()
	Read(uint16) uint8
	Write(uint16, uint8)
	RamSetWriteTrigger(addr uint16, fn func(uint16, uint8)) int
	RamRemoveWriteTrigger(addr uint16, id int)
	ResetTrigger()
	NMITrigger()
	IRQTrigger()
	IRQClear()
	IRQTriggerBind(fn func(uint32))
	IRQClearBind(fn func(uint32))
	SetDMALow(bool)
	BusAvailable() bool
	AECAvailable() bool //TODO NOT STANDARD
	RmwFlags() uint8    //TODO NOT STANDARD
}
