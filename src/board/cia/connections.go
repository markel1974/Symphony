package cia

const intrCiaId = 0x8

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}
