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

type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
}
