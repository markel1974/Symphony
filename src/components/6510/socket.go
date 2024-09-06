package mos6510

type ISocketBanks interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

type ISocketPic interface {
	Reset()
	VerifyIrq(uint8) uint8
	ClearOPFlags()
	ClearNMI()
	HasNMI() bool
	SetOpFlagIrqDisabled()
	SetOpFlagIrqEnabled()
	SetOpFlagIntDelayed()
}
