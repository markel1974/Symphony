package cia

type IVic interface {
	TriggerLightPen()
	ChangedVA(uint8)
}

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}

type IInterrupts interface {
	ClearCIAIRQ()
	TriggerCIAIRQ()
	ClearNMI()
	TriggerNMI()
}
