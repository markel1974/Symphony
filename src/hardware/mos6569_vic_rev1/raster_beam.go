package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
	"strings"
)

const (
	RGBAIndexRed      = 0
	RGBAIndexGreen    = 1
	RGBAIndexBlue     = 2
	RGBAIndexAlpha    = 3
	RGBABytesPerPixel = 4
)

const (
	colorMax            = 0xf
	paletteSize         = 1 << 8
	standardIndexSize   = 1 << 8
	multicolorIndexSize = 1 << 8
)

type Palette [RGBABytesPerPixel][16]uint8

var PaletteBright = Palette{
	{0x00, 0xff, 0x99, 0x00, 0xcc, 0x44, 0x11, 0xff, 0xaa, 0x66, 0xff, 0x40, 0x80, 0x66, 0x77, 0xc0},
	{0x00, 0xff, 0x00, 0xff, 0x00, 0xcc, 0x00, 0xff, 0x55, 0x33, 0x66, 0x40, 0x80, 0xff, 0x77, 0xc0},
	{0x00, 0xff, 0x00, 0xcc, 0xcc, 0x44, 0x99, 0x00, 0x00, 0x00, 0x66, 0x40, 0x80, 0x66, 0xff, 0xc0},
	{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
}

var PaletteOriginal = Palette{
	{0x00, 0xfc, 0x80, 0x87, 0x82, 0x6e, 0x39, 0xdc, 0x8a, 0x52, 0xb7, 0x52, 0x7d, 0xbb, 0x79, 0xaf},
	{0x00, 0xfc, 0x41, 0xc3, 0x46, 0xa9, 0x2d, 0xe9, 0x5c, 0x40, 0x79, 0x52, 0x7d, 0xf9, 0x6c, 0xaf},
	{0x00, 0xfc, 0x32, 0xd2, 0xb4, 0x39, 0xa3, 0x6c, 0x22, 0x03, 0x6a, 0x52, 0x7d, 0x83, 0xea, 0xaf},
	{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
}

// RasterBeam is a structure used to manage raster-based rendering with support for color palettes and scanline operations.
type RasterBeam struct {
	*component.BaseComponent
	reflect            *RasterBeamReflect
	lineWidthRGBA      int
	lineOffsetRGBA     int
	displayBufferArray func(idx int, data *[]uint8, width int)
	scanline           *[]uint8
	colorsRGBA         [paletteSize][RGBABytesPerPixel]uint8
	standardColorIndex [standardIndexSize][8]uint8
	multiColorIndex    [multicolorIndexSize][8]uint8
}

// NewBeam creates and initializes a new RasterBeam instance using the provided display buffer for rendering operations.
func NewBeam(parent references.IComponent, factory references.IComponentFactory, label string, instance int, displayBuffer references.IDisplayBuffer, width int, paletteId string) *RasterBeam {
	palette := PaletteOriginal
	switch strings.ToLower(paletteId) {
	case "bright":
		palette = PaletteBright
	}

	lineWidthRGBA := width * RGBABytesPerPixel
	scanlineSize := uint32(1)
	for scanlineSize <= uint32(lineWidthRGBA*2) {
		scanlineSize <<= 1
	}

	//scanlineMask := int(scanlineSize - 1)
	scanline := make([]uint8, scanlineSize)
	s := &RasterBeam{
		BaseComponent:      component.NewBaseComponent(),
		lineWidthRGBA:      lineWidthRGBA,
		lineOffsetRGBA:     0,
		displayBufferArray: displayBuffer.SetArray,
		scanline:           &scanline,
	}
	for idx := range s.colorsRGBA {
		if idx <= colorMax {
			s.colorsRGBA[idx] = [RGBABytesPerPixel]uint8{
				palette[RGBAIndexRed][idx],
				palette[RGBAIndexGreen][idx],
				palette[RGBAIndexBlue][idx],
				palette[RGBAIndexAlpha][idx],
			}
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
	s.reflect = NewRasterBeamReflect(s, factory, parent, "rasterBeam", instance, references.IdInternalComponent(label, instance, "RasterBeam"))
	return s
}

func (s *RasterBeam) Setup() error {
	return nil
}

func (s *RasterBeam) Connect() error {
	return nil
}

func (s *RasterBeam) EmulationRequired() bool {
	return false
}

func (s *RasterBeam) Emulate() {
}

func (s *RasterBeam) Internal() bool {
	return true
}

func (s *RasterBeam) Reset() {
}

// ResetLineOffset resets the line offset to 0, typically used to prepare for a new rendering cycle or frame.
func (s *RasterBeam) ResetLineOffset() {
	s.lineOffsetRGBA = 0
}

// Draw updates the internal scanline buffer at the computed location with the specified color value.
func (s *RasterBeam) Draw(offset int, color uint8) {
	copy((*s.scanline)[offset*RGBABytesPerPixel:], s.colorsRGBA[color][:])
}

// Draw8 writes an 8-pixel multicolor value to the internal scanline buffer at the specified offset using the given color.
func (s *RasterBeam) Draw8(offset int, color uint8) {
	for i := 0; i < 8; i++ {
		copy((*s.scanline)[(offset+i)*RGBABytesPerPixel:], s.colorsRGBA[color][:])
	}
}

// DrawStandard renders 8 pixels in standard mode into the internal scanline buffer.
func (s *RasterBeam) DrawStandard(offset int, a uint8, b uint8, data uint8) {
	cb := [2][RGBABytesPerPixel]uint8{s.colorsRGBA[a], s.colorsRGBA[b]}
	si := s.standardColorIndex[data]
	for i := 0; i < 8; i++ {
		copy((*s.scanline)[(offset+i)*RGBABytesPerPixel:], (cb[si[i]])[:])
	}
}

// DrawMultiColor renders 8 pixels using a multicolor mode into the internal scanline buffer.
func (s *RasterBeam) DrawMultiColor(offset int, a uint8, b uint8, c uint8, d uint8, data uint8) {
	cb := [4][RGBABytesPerPixel]uint8{s.colorsRGBA[a], s.colorsRGBA[b], s.colorsRGBA[c], s.colorsRGBA[d]}
	mi := s.multiColorIndex[data]
	for i := 0; i < 8; i++ {
		copy((*s.scanline)[(offset+i)*RGBABytesPerPixel:], (cb[mi[i]])[:])
	}
}

// Commit transfers the completed scanline from the internal buffer to the final display buffer.
// This should be called once at the very end of a scanline's rendering cycle.
func (s *RasterBeam) Commit() {
	s.displayBufferArray(s.lineOffsetRGBA, s.scanline, s.lineWidthRGBA)
	s.lineOffsetRGBA += s.lineWidthRGBA
}
