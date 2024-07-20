package vic

type ForeMask struct {
	buf    []uint8 // Foreground mask for sprite-graphics collisions and priorities
	offset int     // Offset from buf
}

func NewForeMask() *ForeMask {
	return &ForeMask{
		buf:    make([]uint8, DisplayXFill+1),
		offset: 0,
	}
}

func (gr *ForeMask) Increment() {
	gr.offset++
}

func (gr *ForeMask) Clear() {
	copy(gr.buf, _emptyForeMaskBuffer)
	gr.offset = 0
}

func (gr *ForeMask) Update(a uint8, b uint8) {
	gr.buf[gr.offset] |= a
	gr.buf[gr.offset+1] |= b
}

func (gr *ForeMask) GetA(m int, s int) uint32 {
	f := ((uint32(gr.buf[m]) << 24) | (uint32(gr.buf[m+1]) << 16) | (uint32(gr.buf[m+2]) << 8) | (uint32(gr.buf[m+3]))) << s
	return f
}

func (gr *ForeMask) GetL(m int, s int) uint32 {
	f := (((uint32(gr.buf[m]) << 24) | (uint32(gr.buf[m+1]) << 16) | (uint32(gr.buf[m+2]) << 8) | (uint32(gr.buf[m+3]))) << s) | (uint32(gr.buf[m+4]) >> (8 - s))
	return f
}

func (gr *ForeMask) GetR(m int, s int) uint32 {
	f := (((uint32(gr.buf[m+4]) << 24) | (uint32(gr.buf[m+5]) << 16) | (uint32(gr.buf[m+6]) << 8) | (uint32(gr.buf[m+7]))) << s) | (uint32(gr.buf[m+8]) >> (8 - s))
	return f
}

//func (gr * ForeMask) GetB(m int, s int) uint32 {
//	f := (((uint32(gr.buf[m]) << 24) | (uint32(gr.buf[m+1]) << 16) | (uint32(gr.buf[m+2]) << 8) | (uint32(gr.buf[m+3]))) << s) | (uint32(gr.buf[m+4]) >> (8 - s))
//	return f
//}

//func (gr * ForeMask) GetC(m int, s int) uint32 {
//	f := (((uint32(gr.buf[m]) << 24) | (uint32(gr.buf[m+1]) << 16) | (uint32(gr.buf[m+2]) << 8) | (uint32(gr.buf[m+3]))) << s) | (uint32(gr.buf[m+4]) >> (8 - s))
//	return f
//}
