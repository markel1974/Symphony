package cia

const IntrCiaId = 0x2

type IVic interface {
	LightPenTrigger()
	ChangedVA(uint8)
}

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}
