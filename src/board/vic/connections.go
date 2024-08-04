package vic

const IntrVicId = 0x1

type IBanks interface {
	ReadCharRom(uint16) uint8
	ReadDirect(uint16) uint8
	ReadColor(uint16) uint8
}

//type IInterrupts interface {
//	TriggerIRQ(uint32)
//	ClearIRQ(uint32)
//}

type IDisplayBuffer interface {
	Set(idx int, data uint8)
	SetMulti8(idx int, data uint8)
	Set8(idx int, data [8]uint8)
}
