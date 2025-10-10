package gl_render

import (
	"github.com/markel1974/c64emu/src/renderers/graphics/gl_render/pixels"
)

// DisplayBuffer is a structure used for managing a buffer of display pixels with an associated Picture and bitmask.
type DisplayBuffer struct {
	p    *pixels.Picture
	mask uint32
}

// NewDisplayBuffer initializes and returns a new DisplayBuffer using the given pixels.Picture.
// It computes and assigns a mask based on the Picture's dimensions.
func NewDisplayBuffer(p *pixels.Picture) *DisplayBuffer {
	db := &DisplayBuffer{
		p: p,
	}
	h := p.Height()
	w := p.Width()
	db.mask = db.computeMask(uint32(h * w))
	return db
}

// SetArray updates the display buffer at the specified index using the provided RGBA data array and width.
func (db *DisplayBuffer) SetArray(idx int, data *[]uint8, width int) {
	db.p.SetRGBADirectArrayPtr(idx, data, width)
}

// computeMask calculates and returns a mask value that is one less than the smallest power of two greater than or equal to num.
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
