package mechanics

type IBanks interface {
	Read(uint16) uint8
	ReadInterval(uint16, uint16) []uint8
}
