package cia

const (
	IRQUnderflowTimerA = 0x1
	IRQUnderflowTimerB = 0x2
	IRQTODAlarmEqual   = 0x4
	IRQSDRFullOtEmpty  = 0x8
	IRQFlagPin         = 0x10
	IRQOccurred        = 0x80
)

const intrCia1Id = 4

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}
