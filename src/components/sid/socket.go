package mos6581

type IPlayer interface {
	GetCurrentPosition() int
	Write([]uint32, int, int)
	Pause()
	Resume()
}

type ISocket interface {
	GetPlayer() IPlayer
}
