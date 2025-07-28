package mos6569

import "github.com/markel1974/c64emu/src/references"

// Beam is a data structure for graphical rendering operations using line-based pixel manipulation.
// lineOffset specifies the current offset used for line rendering calculations.
// displayBufferSet is a function reference for setting pixel data in the display buffer at a specified position.
// displayBufferSet8 is a function for writing 8 consecutive pixel values into the display buffer.
// displayBufferSetMulti8 is used for writing multicolor pixel values to the display buffer.
// colors is an array of uint8 values used to store indexed color mappings for rendering.
// standardColorIndex is a lookup table mapping 8-bit values to arrays of 8 single-bit representations for standard modes.
// multiColorIndex maps 8-bit values to 2-bit pixel indices for multicolor rendering operations.

const (
	paletteSize         = 1 << 8
	standardIndexSize   = 1 << 8
	multicolorIndexSize = 1 << 8
	rgbSize             = 4
)

// Beam represents a rendering component for managing and storing scanline data during graphical operations.
// It handles rendering tasks including multicolor and standard pixel drawing operations into a buffer.
type Beam struct {
	lineWidthRGBA      int
	lineOffsetRGBA     int
	displayBufferArray func(int, []uint8)
	scanline           []uint8
	palette            [paletteSize]uint8
	colorsRGBA         [paletteSize][rgbSize]uint8
	standardColorIndex [standardIndexSize][8]uint8
	multiColorIndex    [multicolorIndexSize][8]uint8
}

// NewBeam creates and initializes a new Beam instance using the provided display buffer for rendering operations.
func NewBeam(displayBuffer references.IDisplayBuffer, width int) *Beam {
	const colorMax = 0xf
	//Bright
	//paletteR := []byte{0x00, 0xff, 0x99, 0x00, 0xcc, 0x44, 0x11, 0xff, 0xaa, 0x66, 0xff, 0x40, 0x80, 0x66, 0x77, 0xc0}
	//paletteG := []byte{0x00, 0xff, 0x00, 0xff, 0x00, 0xcc, 0x00, 0xff, 0x55, 0x33, 0x66, 0x40, 0x80, 0xff, 0x77, 0xc0}
	//paletteB := []byte{0x00, 0xff, 0x00, 0xcc, 0xcc, 0x44, 0x99, 0x00, 0x00, 0x00, 0x66, 0x40, 0x80, 0x66, 0xff, 0xc0}
	//paletteA := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	//Original
	paletteR := []byte{0x00, 0xfc, 0x80, 0x87, 0x82, 0x6e, 0x39, 0xdc, 0x8a, 0x52, 0xb7, 0x52, 0x7d, 0xbb, 0x79, 0xaf}
	paletteG := []byte{0x00, 0xfc, 0x41, 0xc3, 0x46, 0xa9, 0x2d, 0xe9, 0x5c, 0x40, 0x79, 0x52, 0x7d, 0xf9, 0x6c, 0xaf}
	paletteB := []byte{0x00, 0xfc, 0x32, 0xd2, 0xb4, 0x39, 0xa3, 0x6c, 0x22, 0x03, 0x6a, 0x52, 0x7d, 0x83, 0xea, 0xaf}
	paletteA := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	lineWidthRGBA := width * rgbSize
	scanlineSize := uint32(1)
	for scanlineSize <= uint32(lineWidthRGBA*2) {
		scanlineSize <<= 1
	}
	//scanlineMask := int(scanlineSize - 1)
	s := &Beam{
		lineWidthRGBA:      lineWidthRGBA,
		lineOffsetRGBA:     0,
		displayBufferArray: displayBuffer.SetArray,
		scanline:           make([]uint8, scanlineSize),
	}
	for idx := range s.palette {
		s.palette[idx] = (uint8)(idx & colorMax)
	}
	for idx := range s.colorsRGBA {
		if idx <= colorMax {
			s.colorsRGBA[idx] = [rgbSize]uint8{paletteR[idx], paletteG[idx], paletteB[idx], paletteA[idx]}
		} else {
			s.colorsRGBA[idx] = s.colorsRGBA[idx&colorMax]
		}
	}
	for i := range s.multiColorIndex {
		data := i
		idx := uint8(data & 3)
		s.multiColorIndex[i][7] = idx
		s.multiColorIndex[i][6] = idx
		data >>= 2
		idx = uint8(data & 3)
		s.multiColorIndex[i][5] = idx
		s.multiColorIndex[i][4] = idx
		data >>= 2
		idx = uint8(data & 3)
		s.multiColorIndex[i][3] = idx
		s.multiColorIndex[i][2] = idx
		data >>= 2
		idx = uint8(data)
		s.multiColorIndex[i][1] = idx
		s.multiColorIndex[i][0] = idx
	}
	for i := range s.standardColorIndex {
		data := uint8(i)
		s.standardColorIndex[i][7] = data & 1
		data >>= 1
		s.standardColorIndex[i][6] = data & 1
		data >>= 1
		s.standardColorIndex[i][5] = data & 1
		data >>= 1
		s.standardColorIndex[i][4] = data & 1
		data >>= 1
		s.standardColorIndex[i][3] = data & 1
		data >>= 1
		s.standardColorIndex[i][2] = data & 1
		data >>= 1
		s.standardColorIndex[i][1] = data & 1
		data >>= 1
		s.standardColorIndex[i][0] = data & 1
	}
	return s
}

// ResetLineOffset resets the line offset to 0, typically used to prepare for a new rendering cycle or frame.
func (s *Beam) ResetLineOffset() {
	s.lineOffsetRGBA = 0
}

// Draw updates the internal scanline buffer at the computed location with the specified color value.
func (s *Beam) Draw(offset int, color uint8) {
	copy(s.scanline[offset*rgbSize:], s.colorsRGBA[color][:])
}

// DrawMulti8 writes an 8-pixel multicolor value to the internal scanline buffer at the specified offset using the given color.
func (s *Beam) DrawMulti8(offset int, color uint8) {
	for i := 0; i < 8; i++ {
		copy(s.scanline[(offset+i)*rgbSize:], s.colorsRGBA[color][:])
	}
}

// Draw8Standard renders 8 pixels in standard mode into the internal scanline buffer.
func (s *Beam) Draw8Standard(offset int, a uint8, b uint8, data uint8) {
	cb := [2]uint8{s.palette[a], s.palette[b]}
	si := s.standardColorIndex[data]
	for i := 0; i < 8; i++ {
		copy(s.scanline[(offset+i)*rgbSize:], (s.colorsRGBA[cb[si[i]]])[:])
	}
}

// Draw8Multi renders 8 pixels using a multicolor mode into the internal scanline buffer.
func (s *Beam) Draw8Multi(offset int, a uint8, b uint8, c uint8, d uint8, data uint8) {
	cb := [4]uint8{s.palette[a], s.palette[b], s.palette[c], s.palette[d]}
	mi := s.multiColorIndex[data]
	for i := 0; i < 8; i++ {
		copy(s.scanline[(offset+i)*rgbSize:], (s.colorsRGBA[cb[mi[i]]])[:])
	}
}

// Commit transfers the completed scanline from the internal buffer to the final display buffer.
// This should be called once at the very end of a scanline's rendering cycle.
func (s *Beam) Commit() {
	s.displayBufferArray(s.lineOffsetRGBA, s.scanline[:s.lineWidthRGBA])
	s.lineOffsetRGBA += s.lineWidthRGBA
}
