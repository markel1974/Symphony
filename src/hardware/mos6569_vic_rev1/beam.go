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
	scanlineSize        = 1 << 11 // 2048 Based on sequencer and border logic, a size of 576 is safe.
	scanlineMask        = scanlineSize - 1
)

type Beam struct {
	lineOffset             int
	displayBufferSet       func(int, uint8)
	displayBufferSet8      func(int, *[8]uint8)
	displayBufferSetMulti8 func(int, uint8)
	// scanline holds the pixel data for an entire scanline before committing it to the final display.
	scanline [scanlineSize]uint8
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

// Commit transfers the completed scanline from the internal buffer to the final display buffer.
// This should be called once at the very end of a scanline's rendering cycle.
func (s *Beam) Commit() {
	// In a real implementation, this would ideally be a single, highly optimized call
	// to the display buffer interface, e.g., displayBuffer.SetScanline(s.lineOffset, s.scanline[:]).
	// For now, we simulate it by iterating, which is less performant but functionally correct.
	const visibleWidth = 576
	for i := 0; i < visibleWidth; i++ {
		s.displayBufferSet(s.lineOffset+i, s.scanline[i])
	}
}

// SetOffset sets the line offset for the Beam to the specified value.
func (s *Beam) SetOffset(offset int) {
	s.lineOffset = offset
}

// Draw updates the internal scanline buffer at the computed location with the specified color value.
func (s *Beam) Draw(index int, color uint8) {
	s.scanline[index&scanlineMask] = s.colors[color]
}

// DrawMulti8 writes an 8-pixel multicolor value to the internal scanline buffer at the specified offset using the given color.
func (s *Beam) DrawMulti8(offset int, color uint8) {
	finalColor := s.colors[color]
	for i := 0; i < 8; i++ {
		s.scanline[(offset+i)&scanlineMask] = finalColor
	}
}

// Draw8Standard renders 8 pixels in standard mode into the internal scanline buffer.
// Draw8Standard renderizza 8 pixel in modalità standard nel buffer interno.
func (s *Beam) Draw8Standard(offset int, a uint8, b uint8, data uint8) {
	colorBuffer := [2]uint8{s.colors[a], s.colors[b]}
	colorIndex := s.standardIndex[data]
	s.scanline[(offset+0)&scanlineMask] = colorBuffer[colorIndex[0]]
	s.scanline[(offset+1)&scanlineMask] = colorBuffer[colorIndex[1]]
	s.scanline[(offset+2)&scanlineMask] = colorBuffer[colorIndex[2]]
	s.scanline[(offset+3)&scanlineMask] = colorBuffer[colorIndex[3]]
	s.scanline[(offset+4)&scanlineMask] = colorBuffer[colorIndex[4]]
	s.scanline[(offset+5)&scanlineMask] = colorBuffer[colorIndex[5]]
	s.scanline[(offset+6)&scanlineMask] = colorBuffer[colorIndex[6]]
	s.scanline[(offset+7)&scanlineMask] = colorBuffer[colorIndex[7]]
}

// Draw8Multi renders 8 pixels using a multicolor mode into the internal scanline buffer.
// Draw8Multi renderizza 8 pixel in modalità multicolor nel buffer interno.
func (s *Beam) Draw8Multi(offset int, a uint8, b uint8, c uint8, d uint8, data uint8) {
	colorBuffer := [4]uint8{s.colors[a], s.colors[b], s.colors[c], s.colors[d]}
	index := s.multicolorIndex[data]
	s.scanline[(offset+0)&scanlineMask] = colorBuffer[index[0]]
	s.scanline[(offset+1)&scanlineMask] = colorBuffer[index[1]]
	s.scanline[(offset+2)&scanlineMask] = colorBuffer[index[2]]
	s.scanline[(offset+3)&scanlineMask] = colorBuffer[index[3]]
	s.scanline[(offset+4)&scanlineMask] = colorBuffer[index[4]]
	s.scanline[(offset+5)&scanlineMask] = colorBuffer[index[5]]
	s.scanline[(offset+6)&scanlineMask] = colorBuffer[index[6]]
	s.scanline[(offset+7)&scanlineMask] = colorBuffer[index[7]]
}
