package mos6569

import "github.com/markel1974/c64emu/src/references"

// Beam is a data structure for graphical rendering operations using line-based pixel manipulation.
// lineOffset specifies the current offset used for line rendering calculations.
// displayBufferSet is a function reference for setting pixel data in the display buffer at a specified position.
// displayBufferSet8 is a function for writing 8 consecutive pixel values into the display buffer.
// displayBufferSetMulti8 is used for writing multicolor pixel values to the display buffer.
// colors is an array of uint8 values used to store indexed color mappings for rendering.
// standardIndex is a lookup table mapping 8-bit values to arrays of 8 single-bit representations for standard modes.
// multicolorIndex maps 8-bit values to 2-bit pixel indices for multicolor rendering operations.

const (
	colorSize           = 1 << 8
	standardIndexSize   = 1 << 8
	multicolorIndexSize = 1 << 8
)

type Beam struct {
	lineOffset             int
	displayBufferSet       func(int, uint8)
	displayBufferSet8      func(int, *[8]uint8)
	displayBufferSetMulti8 func(int, uint8)
	// colors is an array that holds 256 color values, initialized with specific uint8 values for rendering purposes.
	colors [colorSize]uint8
	// standardIndex is a lookup table that maps an 8-bit value to an array of 8 single-bit values extracted from it.
	standardIndex [standardIndexSize][8]uint8
	// multicolorIndex is a lookup table that maps each 8-bit value to an array of 8 corresponding 2-bit pixel indices.
	multicolorIndex [multicolorIndexSize][8]uint8
}

// NewBeam creates and initializes a new Beam instance using the provided display buffer for rendering operations.
func NewBeam(displayBuffer references.IDisplayBuffer) *Beam {
	s := &Beam{
		lineOffset:             0,
		displayBufferSet:       displayBuffer.Set,
		displayBufferSet8:      displayBuffer.Set8,
		displayBufferSetMulti8: displayBuffer.SetMulti8,
	}
	for i := range s.colors {
		s.colors[i] = (uint8)(i & 0xf)
	}
	for i := range s.multicolorIndex {
		data := i
		idx := uint8(data & 3)
		s.multicolorIndex[i][7] = idx
		s.multicolorIndex[i][6] = idx
		data >>= 2
		idx = uint8(data & 3)
		s.multicolorIndex[i][5] = idx
		s.multicolorIndex[i][4] = idx
		data >>= 2
		idx = uint8(data & 3)
		s.multicolorIndex[i][3] = idx
		s.multicolorIndex[i][2] = idx
		data >>= 2
		idx = uint8(data) // non serve &3, sono gli ultimi 2 bit
		s.multicolorIndex[i][1] = idx
		s.multicolorIndex[i][0] = idx
	}

	for i := range s.standardIndex {
		data := uint8(i)
		s.standardIndex[i][7] = data & 1
		data >>= 1
		s.standardIndex[i][6] = data & 1
		data >>= 1
		s.standardIndex[i][5] = data & 1
		data >>= 1
		s.standardIndex[i][4] = data & 1
		data >>= 1
		s.standardIndex[i][3] = data & 1
		data >>= 1
		s.standardIndex[i][2] = data & 1
		data >>= 1
		s.standardIndex[i][1] = data & 1
		data >>= 1
		s.standardIndex[i][0] = data & 1
	}
	return s
}

// SetOffset sets the line offset for the Beam to the specified value.
func (s *Beam) SetOffset(offset int) {
	s.lineOffset = offset
}

// Draw updates the display buffer at the computed location with the specified color value.
func (s *Beam) Draw(index int, color uint8) {
	s.displayBufferSet(s.lineOffset+index, s.colors[color])
}

// DrawMulti8 writes an 8-pixel multicolor value to the display buffer at the specified offset using the given color.
func (s *Beam) DrawMulti8(offset int, color uint8) {
	s.displayBufferSetMulti8(s.lineOffset+offset, s.colors[color])
}

// Draw8Standard renders 8 pixels in standard mode using two colors based on the bit values in the provided data byte.
func (s *Beam) Draw8Standard(offset int, a uint8, b uint8, data uint8) {
	colorBuffer := [4]uint8{s.colors[a], s.colors[b], 0, 0}
	colorIndex := s.standardIndex[data]
	drawBuffer := [8]uint8{
		colorBuffer[colorIndex[0]], colorBuffer[colorIndex[1]], colorBuffer[colorIndex[2]], colorBuffer[colorIndex[3]],
		colorBuffer[colorIndex[4]], colorBuffer[colorIndex[5]], colorBuffer[colorIndex[6]], colorBuffer[colorIndex[7]],
	}
	s.displayBufferSet8(s.lineOffset+offset, &drawBuffer)
}

// Draw8Multi renders 8 pixels using a multicolor mode where each pixel is selected from four input colors (a, b, c, d).
// It maps an 8-bit data value to 2-bit pixel indices using a multicolorIndex lookup and writes the resulting color values.
// The rendering output is written to the display buffer starting from the calculated offset position.
func (s *Beam) Draw8Multi(offset int, a uint8, b uint8, c uint8, d uint8, data uint8) {
	colorBuffer := [4]uint8{s.colors[a], s.colors[b], s.colors[c], s.colors[d]}
	index := s.multicolorIndex[data]
	drawBuffer := [8]uint8{
		colorBuffer[index[0]], colorBuffer[index[1]], colorBuffer[index[2]], colorBuffer[index[3]],
		colorBuffer[index[4]], colorBuffer[index[5]], colorBuffer[index[6]], colorBuffer[index[7]],
	}
	s.displayBufferSet8(s.lineOffset+offset, &drawBuffer)
}
