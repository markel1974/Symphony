package vic

type IBanks interface {
	ReadCharRom(uint16) uint8
	ReadDirect(uint16) uint8
	ReadColor(uint16) uint8
}

type IInterrupts interface {
	TriggerVICIRQ()
	ClearVICIRQ()
}
