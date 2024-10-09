package void

type Void struct{}

func NewDisk() *Void {
	return &Void{}
}

func (e *Void) TrackLen() int {
	return 0
}

func (e *Void) TrackSectors() uint8 {
	return 0
}

func (e *Void) Read() uint8 {
	return 0
}

func (e *Void) Write(_ uint8) {
}

func (e *Void) Next() uint8 {
	return 0
}

func (e *Void) SetHeadTrack(uint8) int {
	return 0
}

func (e *Void) MicroSecPerByte() uint8 {
	return 0
}

func (e *Void) Rotate() {
}

func (e *Void) Usable() bool {
	return false
}
