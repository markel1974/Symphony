package mos6581

type ISidPlayer interface {
	GetCurrentPosition() int
	Write([]uint32, int, int)
	Pause()
	Resume()
}
