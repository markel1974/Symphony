package mos6510fn

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
	GetNMICycleDistance(int) uint64
	GetIrqCycleDistance(int) uint64
}
