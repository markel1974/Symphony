package banks

type IWiring interface {
	WriteRegister(addr uint16, data uint8)
	ReadRegister(addr uint16) uint8
	GetLastByte() uint8
}
