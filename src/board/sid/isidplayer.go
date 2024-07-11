package sid

type ISidPlayer interface {
	GetCurrentPosition() int
	Write([]uint32, int, int)
	Pause()
	Resume()
}
