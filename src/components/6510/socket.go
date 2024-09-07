package mos6510

type IBanks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

type IPic interface {
	Reset()
	VerifyIrq(uint8, uint8) uint8
	ClearNMI()
	HasNMI() bool
}

type ISocket interface {
	GetBanks() IBanks
	GetPic() IPic
}
