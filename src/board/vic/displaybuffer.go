package vic

type DisplayBuffer struct {
	colorMultiplier [][]uint8
	displayBuffer   []uint8
}

func NewDisplayBuffer() *DisplayBuffer {
	db := &DisplayBuffer{
		displayBuffer:   make([]uint8, DisplaySize),
		colorMultiplier: make([][]uint8, 0xff),
	}
	for idx := range db.colorMultiplier {
		x := uint8(idx)
		db.colorMultiplier[idx] = []uint8{x, x, x, x, x, x, x, x}
	}
	return db
}

func (db *DisplayBuffer) Get() []uint8 {
	return db.displayBuffer
}

func (db *DisplayBuffer) Set(idx int, data uint8) {
	db.displayBuffer[idx] = data
}

func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
	copy(db.displayBuffer[idx:], db.colorMultiplier[data])
}
