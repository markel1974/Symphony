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

type IInterrupts interface {
	ClearIRQ(uint32)
	TriggerIRQ(uint32)
	ClearNMI()
	TriggerNMI()
}
