package mos6569

type IBanks interface {
	ReadCharRom(uint16) uint8
	ReadDirect(uint16) uint8
	ReadColor(uint16) uint8
}

type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}

type ISocket interface {
	Cycle() uint64
	GetDisplayBuffer() IDisplayBuffer
	GetBanks() IBanks
	Ready()
	IRQTrigger()
	IRQClear()
	BALow(d bool)
	AECLow(d bool)
	VBlank()
	LastCycle()
}
