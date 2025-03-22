package gl_render

import (
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

type DisplayBuffer struct {
	p       *pixels.Picture
	coords  []int
	colors  [][]uint8
	colors8 [][]uint8
	maxLen  int
}

func NewDisplayBuffer(p *pixels.Picture) *DisplayBuffer {
	paletteR := []byte{0x00, 0xff, 0x99, 0x00, 0xcc, 0x44, 0x11, 0xff, 0xaa, 0x66, 0xff, 0x40, 0x80, 0x66, 0x77, 0xc0}
	paletteG := []byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xcc, 0x00, 0xff, 0x55, 0x33, 0x66, 0x40, 0x80, 0xff, 0x77, 0xc0}
	paletteB := []byte{0x00, 0xff, 0x00, 0xcc, 0xcc, 0x44, 0x99, 0x00, 0x00, 0x00, 0x66, 0x40, 0x80, 0x66, 0xff, 0xc0}
	//Original
	//paletteR := []byte{0x00, 0xfc, 0x80, 0x87, 0x82, 0x6e, 0x39, 0xdc, 0x8a, 0x52, 0xb7, 0x52, 0x7d, 0xbb, 0x79, 0xaf}
	//paletteG := []byte{0x00, 0xfc, 0x41, 0xc3, 0x46, 0xa9, 0x2d, 0xe9, 0x5c, 0x40, 0x79, 0x52, 0x7d, 0xf9, 0x6c, 0xaf}
	//paletteB := []byte{0x00, 0xfc, 0x32, 0xd2, 0xb4, 0x39, 0xa3, 0x6c, 0x22, 0x03, 0x6a, 0x52, 0x7d, 0x83, 0xea, 0xaf}

	var coords []int
	h := p.Height()
	w := p.Width()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			coords = append(coords, p.ComputeIndex(x, y))
		}
	}
	colors := make([][]uint8, 256)
	colors8 := make([][]uint8, 256)
	for j := 0; j < 16; j++ {
		red := paletteR[j]
		green := paletteG[j]
		blue := paletteB[j]
		alfa := uint8(255)
		rgba := []uint8{red, green, blue, alfa}
		colors[j] = rgba
		rgba8 := []uint8{
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
	return &DisplayBuffer{
		p:       p,
		coords:  coords,
		colors:  colors,
		colors8: colors8,
		maxLen:  len(coords),
	}
}

func (db *DisplayBuffer) Set(idx int, data uint8) {
	if idx >= db.maxLen {
		return
	}
	db.p.SetRGBADirectArray(db.coords[idx], db.colors[data])
}

func (db *DisplayBuffer) Set8(idx int, data [8]uint8) {
	if idx+7 >= db.maxLen {
		return
	}
	//t := []uint8{
	//	db.colors[data[0]][0], db.colors[data[0]][1], db.colors[data[0]][2], db.colors[data[0]][3],
	//	db.colors[data[1]][0], db.colors[data[1]][1], db.colors[data[1]][2], db.colors[data[1]][3],
	//	db.colors[data[2]][0], db.colors[data[2]][1], db.colors[data[2]][2], db.colors[data[2]][3],
	//	db.colors[data[3]][0], db.colors[data[3]][1], db.colors[data[3]][2], db.colors[data[3]][3],
	//	db.colors[data[4]][0], db.colors[data[4]][1], db.colors[data[4]][2], db.colors[data[4]][3],
	//	db.colors[data[5]][0], db.colors[data[5]][1], db.colors[data[5]][2], db.colors[data[5]][3],
	//	db.colors[data[6]][0], db.colors[data[6]][1], db.colors[data[6]][2], db.colors[data[6]][3],
	//	db.colors[data[7]][0], db.colors[data[7]][1], db.colors[data[7]][2], db.colors[data[7]][3],
	//}
	//db.p.SetRGBADirectArray(db.coords[idx], t)
	for x := 0; x < 8; x++ {
		db.p.SetRGBADirectArray(db.coords[idx+x], db.colors[data[x]])
	}
}

func (db *DisplayBuffer) SetMulti8(idx int, data uint8) {
	if idx+7 >= db.maxLen {
		return
	}
	db.p.SetRGBADirectArray(db.coords[idx], db.colors8[data])
	//color := db.colors[data]
	//for x := 0; x < 8; x++ {
	//	db.p.SetRGBADirectArray(db.coords[idx+x], color)
	//}
}
