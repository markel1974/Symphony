package asciirender

type DisplayBuffer struct {
}

func NewDisplayBuffer() *DisplayBuffer {
	return &DisplayBuffer{}
}

func (db *DisplayBuffer) Set(idx int, data uint8) {
}

func (db *DisplayBuffer) Set8(idx int, data [8]uint8) {
}

func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
}
