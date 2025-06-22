package gl_render

import (
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

type DisplayBuffer struct {
	p       *pixels.Picture
	coords  []int
	colors  [256]*[4]uint8
	colors8 [256]*[32]uint8
	mask    uint32
}

func NewDisplayBuffer(p *pixels.Picture) *DisplayBuffer {
	paletteR := []byte{0x00, 0xff, 0x99, 0x00, 0xcc, 0x44, 0x11, 0xff, 0xaa, 0x66, 0xff, 0x40, 0x80, 0x66, 0x77, 0xc0}
	paletteG := []byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xcc, 0x00, 0xff, 0x55, 0x33, 0x66, 0x40, 0x80, 0xff, 0x77, 0xc0}
	paletteB := []byte{0x00, 0xff, 0x00, 0xcc, 0xcc, 0x44, 0x99, 0x00, 0x00, 0x00, 0x66, 0x40, 0x80, 0x66, 0xff, 0xc0}
	//Original
	//paletteR := []byte{0x00, 0xfc, 0x80, 0x87, 0x82, 0x6e, 0x39, 0xdc, 0x8a, 0x52, 0xb7, 0x52, 0x7d, 0xbb, 0x79, 0xaf}
	//paletteG := []byte{0x00, 0xfc, 0x41, 0xc3, 0x46, 0xa9, 0x2d, 0xe9, 0x5c, 0x40, 0x79, 0x52, 0x7d, 0xf9, 0x6c, 0xaf}
	//paletteB := []byte{0x00, 0xfc, 0x32, 0xd2, 0xb4, 0x39, 0xa3, 0x6c, 0x22, 0x03, 0x6a, 0x52, 0x7d, 0x83, 0xea, 0xaf}

	db := &DisplayBuffer{
		p: p,
	}
	h := p.Height()
	w := p.Width()
	db.mask = db.computeMask(uint32(h * w))
	db.coords = make([]int, db.mask)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			index := (y * w) + x
			db.coords[index] = p.ComputeIndex(x, y)
		}
	}

	for j := 0; j < 16; j++ {
		red := paletteR[j]
		green := paletteG[j]
		blue := paletteB[j]
		alfa := uint8(255)
		rgba := [4]uint8{red, green, blue, alfa}
		db.colors[j] = &rgba
		rgba8 := [32]uint8{
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
			red, green, blue, alfa,
		}
		db.colors8[j] = &rgba8
	}
	for k := 16; k < 256; k++ {
		db.colors[k] = db.colors[k&0x0f]
		db.colors8[k] = db.colors8[k&0x0f]
	}
	return db
}

func (db *DisplayBuffer) Set(idx int, data uint8) {
	idx = idx & int(db.mask)
	db.p.SetRGBA4DirectArrayPtr(db.coords[idx], db.colors[data])
}

func (db *DisplayBuffer) Set8(idx int, data *[8]uint8) {
	idx = idx & int(db.mask)
	colors := db.colors
	d0 := data[0]
	d1 := data[1]
	d2 := data[2]
	d3 := data[3]
	d4 := data[4]
	d5 := data[5]
	d6 := data[6]
	d7 := data[7]

	//z0 := colors[data[0]]
	//z1 := colors[data[1]]
	t := [32]uint8{
		colors[d0][0], colors[d0][1], colors[d0][2], colors[d0][3],
		colors[d1][0], colors[d1][1], colors[d1][2], colors[d1][3],
		colors[d2][0], colors[d2][1], colors[d2][2], colors[d2][3],
		colors[d3][0], colors[d3][1], colors[d3][2], colors[d3][3],
		colors[d4][0], colors[d4][1], colors[d4][2], colors[d4][3],
		colors[d5][0], colors[d5][1], colors[d5][2], colors[d5][3],
		colors[d6][0], colors[d6][1], colors[d6][2], colors[d6][3],
		colors[d7][0], colors[d7][1], colors[d7][2], colors[d7][3],
	}
	db.p.SetRGBA32DirectArrayPtr(db.coords[idx], &t)
}

func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
	idx = idx & int(db.mask)
	db.p.SetRGBA32DirectArrayPtr(db.coords[idx], db.colors8[data])
}

func (db *DisplayBuffer) computeMask(num uint32) uint32 {
	if num == 0 {
		return 0
	}
	power := uint32(1)
	for power <= num {
		power <<= 1
	}
	return power - 1
}
