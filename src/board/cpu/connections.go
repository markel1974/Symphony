package cpu

type IBanks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}
