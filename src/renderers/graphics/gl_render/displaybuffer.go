package gl_render

import (
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

// DisplayBuffer represents a buffer for rendering pixel data using predefined colors and coordinates.
type DisplayBuffer struct {
	p    *pixels.Picture
	mask uint32
}

// NewDisplayBuffer creates and initializes a new DisplayBuffer using the given pixels.Picture instance.
func NewDisplayBuffer(p *pixels.Picture) *DisplayBuffer {
	db := &DisplayBuffer{
		p: p,
	}
	h := p.Height()
	w := p.Width()
	db.mask = db.computeMask(uint32(h * w))
	return db
}

func (db *DisplayBuffer) SetArray(idx int, data []uint8) {
	db.p.SetRGBADirectArrayPtr(idx, &data)
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
