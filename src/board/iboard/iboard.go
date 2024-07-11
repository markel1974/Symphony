package iboard

import "github.com/markel1974/c64emu/src/board/quartz"

type IBoard interface {
	RmwFlags() uint8

	Cycle() uint64
	CreateAlarm(string, quartz.AlarmCallback) *quartz.Alarm

	CpuRamWrite(uint16, uint8)
	CpuRamRead(uint16) uint8

	//BasicRomRead(uint16) uint8

	CharRomRead(uint16) uint8
	//KernalRomRead(uint16) uint8

	RamRead(uint16) uint8
	RamWrite(uint16, uint8)

	ColorRead(uint16) uint8
	//ColorWrite(uint16, uint8)

	ReadyEvent()

	VICLightPenTrigger()
	VICChangedVA(uint8)
	VICTriggerIRQ()
	VICClearIRQ()

	CIATriggerIRQ()
	CIAClearIRQ()

	NMITrigger()
	NMIClear()

	BusCpuWrite(uint8)
	BusCpuRead() uint8

	ExtRamWrite(int, uint16, uint8)
	ExtRamRead(int, uint16) uint8

	//VICLastByte() uint8
	//VICVBlank()
}
