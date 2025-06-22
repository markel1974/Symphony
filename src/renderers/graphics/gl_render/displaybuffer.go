package gl_render

import (
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

// DisplayBuffer represents a buffer for rendering pixel data using predefined colors and coordinates.
type DisplayBuffer struct {
	p       *pixels.Picture
	coords  []int
	colors  [256]*[4]uint8
	colors8 [256]*[32]uint8
	mask    uint32
}

// NewDisplayBuffer creates and initializes a new DisplayBuffer using the given pixels.Picture instance.
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

// Set updates a single pixel at the specified index with the color data provided, using a masked index calculation.
func (db *DisplayBuffer) Set(idx int, data uint8) {
	idx = idx & int(db.mask)
	db.p.SetRGBA4DirectArrayPtr(db.coords[idx], db.colors[data])
}

// Set8 updates a group of 8 pixels in the display buffer using the provided 8-color palette data.
func (db *DisplayBuffer) Set8(idx int, data *[8]uint8) {
	idx = idx & int(db.mask)
	colors := db.colors
	z0 := colors[data[0]]
	z1 := colors[data[1]]
	z2 := colors[data[2]]
	z3 := colors[data[3]]
	z4 := colors[data[4]]
	z5 := colors[data[5]]
	z6 := colors[data[6]]
	z7 := colors[data[7]]
	t := [32]uint8{
		z0[0], z0[1], z0[2], z0[3],
		z1[0], z1[1], z1[2], z1[3],
		z2[0], z2[1], z2[2], z2[3],
		z3[0], z3[1], z3[2], z3[3],
		z4[0], z4[1], z4[2], z4[3],
		z5[0], z5[1], z5[2], z5[3],
		z6[0], z6[1], z6[2], z6[3],
		z7[0], z7[1], z7[2], z7[3],
	}
	db.p.SetRGBA32DirectArrayPtr(db.coords[idx], &t)
}

// SetMulti8 modifies the pixel data at the specified index using an 8-color RGBA32 palette mapped by the data value.
func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
	idx = idx & int(db.mask)
	db.p.SetRGBA32DirectArrayPtr(db.coords[idx], db.colors8[data])
}

// computeMask calculates and returns the bitmask for a given number, which is the largest power of two minus one.
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
