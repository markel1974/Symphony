package wasm_render

const rgbaLen = 4
const rgba8Len = rgbaLen * 8

type DisplayBuffer struct {
	colors  [][rgbaLen]uint8
	colors8 [][rgba8Len]uint8
	surface []byte
	maxLen  int
}

func NewDisplayBuffer(w int, h int) *DisplayBuffer {
	paletteR := []byte{0x00, 0xff, 0x99, 0x00, 0xcc, 0x44, 0x11, 0xff, 0xaa, 0x66, 0xff, 0x40, 0x80, 0x66, 0x77, 0xc0}
	paletteG := []byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xcc, 0x00, 0xff, 0x55, 0x33, 0x66, 0x40, 0x80, 0xff, 0x77, 0xc0}
	paletteB := []byte{0x00, 0xff, 0x00, 0xcc, 0xcc, 0x44, 0x99, 0x00, 0x00, 0x00, 0x66, 0x40, 0x80, 0x66, 0xff, 0xc0}
	//Original
	//paletteR := []byte{0x00, 0xfc, 0x80, 0x87, 0x82, 0x6e, 0x39, 0xdc, 0x8a, 0x52, 0xb7, 0x52, 0x7d, 0xbb, 0x79, 0xaf}
	//paletteG := []byte{0x00, 0xfc, 0x41, 0xc3, 0x46, 0xa9, 0x2d, 0xe9, 0x5c, 0x40, 0x79, 0x52, 0x7d, 0xf9, 0x6c, 0xaf}
	//paletteB := []byte{0x00, 0xfc, 0x32, 0xd2, 0xb4, 0x39, 0xa3, 0x6c, 0x22, 0x03, 0x6a, 0x52, 0x7d, 0x83, 0xea, 0xaf}
	colors := make([][rgbaLen]uint8, 256)
	colors8 := make([][rgba8Len]uint8, 256)
	for j := 0; j < 16; j++ {
		red := paletteR[j]
		green := paletteG[j]
		blue := paletteB[j]
		alfa := uint8(255)
		rgba := [rgbaLen]uint8{red, green, blue, alfa}
		colors[j] = rgba
		rgba8 := [rgba8Len]uint8{
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
		}
		colors8[j] = rgba8
	}
	for k := 16; k < 256; k++ {
		colors[k] = colors[k&0x0f]
		colors8[k] = colors8[k&0x0f]
	}
	surface := make([]byte, w*h*4)
	return &DisplayBuffer{
		colors:  colors,
		colors8: colors8,
		surface: surface,
		maxLen:  len(surface),
	}
}

func (db *DisplayBuffer) GetSurface() []byte {
	return db.surface
}

func (db *DisplayBuffer) Set(idx int, data uint8) {
	target := idx * rgbaLen
	if (target + rgbaLen) > db.maxLen {
		return
	}
	copy(db.surface[target:], db.colors[data][:])
}

func (db *DisplayBuffer) Set8(idx int, data [8]uint8) {
	if max := (idx * rgbaLen) + (8 * rgbaLen); max > db.maxLen {
		return
	}
	copy(db.surface[(idx+0)*rgbaLen:], db.colors[data[0]][:])
	copy(db.surface[(idx+1)*rgbaLen:], db.colors[data[1]][:])
	copy(db.surface[(idx+2)*rgbaLen:], db.colors[data[2]][:])
	copy(db.surface[(idx+3)*rgbaLen:], db.colors[data[3]][:])
	copy(db.surface[(idx+4)*rgbaLen:], db.colors[data[4]][:])
	copy(db.surface[(idx+5)*rgbaLen:], db.colors[data[5]][:])
	copy(db.surface[(idx+6)*rgbaLen:], db.colors[data[6]][:])
	copy(db.surface[(idx+7)*rgbaLen:], db.colors[data[7]][:])
}

func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
	target := idx * rgbaLen
	if (target + rgba8Len) > db.maxLen {
		return
	}
	copy(db.surface[target:], db.colors8[data][:])
}
