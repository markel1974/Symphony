package mos6510

type IBanks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

type IPic interface {
	HasAny() bool
	HasReset() bool
	HasNMI() bool
	HasIRQ() bool
	ClearNMI()
	GetCycle() uint64
	GetNMICycleDistance(int) uint64
	GetIrqCycleDistance(int) uint64
}
