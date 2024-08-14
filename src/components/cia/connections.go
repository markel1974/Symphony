package cia

const intrCia1Id = 4

type IBus interface {
	CpuRead() uint8
	CpuWrite(uint8)
}
