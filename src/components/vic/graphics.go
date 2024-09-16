package mos6569

import "fmt"

const (
	modeTextStandard            = 0
	modeTextMulticolor          = 1
	modeBitmapStandard          = 2
	modeBitmapMulticolor        = 3
	modeTextECM                 = 4
	modeTextMulticolorInvalid   = 5
	modeBitmapStandardInvalid   = 6
	modeBitmapMulticolorInvalid = 7
)

const columnsMax = 40
const rowsMax = 7

type Graphics struct {
	core             *Core
	collisions       *Collisions
	db               IDisplayBuffer
	gfxData          uint8
	colorData        uint8
	charData         uint8
	charDataLast     uint8
	offset           int     // Offset from bitmap spritesBuffer
	lineIndex        int     // Index in video matrix / color line
	videoMatrix      []uint8 // Video matrix spritesBuffer
	colorLine        []uint8 // Color line spritesBuffer
	rowCounter       uint16  // Row counter
	videoCounter     uint16  // Video counter
	videoCounterBase uint16  // Video counter base
	displayAccess    bool    // Display state
	textBuffer       []byte
}

func NewGraphics(core *Core, collisions *Collisions, db IDisplayBuffer) *Graphics {
	gr := &Graphics{
		core:             core,
		collisions:       collisions,
		db:               db,
		gfxData:          0,
		colorData:        0,
		charData:         0,
		charDataLast:     0,
		offset:           0,
		lineIndex:        0,
		videoMatrix:      make([]uint8, columnsMax),
		colorLine:        make([]uint8, columnsMax),
		textBuffer:       make([]uint8, (RasterYMax/8)*columnsMax),
		rowCounter:       rowsMax,
		videoCounter:     0,
		videoCounterBase: 0,
		displayAccess:    false,
	}
	return gr
}

func (gr *Graphics) PrintText() {
	for x, v := range gr.textBuffer {
		if (x % 40) == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", v)
	}
	fmt.Println()
}

func (gr *Graphics) Setup() {
}

func (gr *Graphics) GetText() []byte {
	return gr.textBuffer
}

func (gr *Graphics) ResetVideoCounterBase() {
	gr.videoCounterBase = 0
}

func (gr *Graphics) ResetLineIndex() {
	gr.lineIndex = 0
}

func (gr *Graphics) UpdateVideoCounter() {
	gr.videoCounter = gr.videoCounterBase
}

func (gr *Graphics) SetOffset(offset int) {
	gr.offset = offset
}

func (gr *Graphics) UpdateLastCharData() {
	gr.charDataLast = gr.charData
}

func (gr *Graphics) TryResetRowCounter() {
	if gr.core.badLineCondition {
		gr.rowCounter = 0
	}
}

func (gr *Graphics) TryAcquireDisplayAccess() {
	if gr.core.badLineCondition {
		gr.displayAccess = true
	}
}

func (gr *Graphics) UpdateDisplayAccess() {
	// TODO VERIFY
	if gr.rowCounter == rowsMax {
		gr.videoCounterBase = gr.videoCounter
		gr.displayAccess = false
	}
	if gr.core.badLineCondition || gr.displayAccess {
		gr.rowCounter = (gr.rowCounter + 1) & rowsMax
		gr.displayAccess = true
	}
}

func (gr *Graphics) TryGraphicsAccess() {
	if gr.displayAccess {
		var addr uint16
		if (gr.core.cr1 & 0x20) != 0 {
			addr = ((gr.videoCounter & 0x03ff) << 3) | gr.core.bitmapBase | gr.rowCounter // Bitmap
		} else {
			addr = (uint16(gr.videoMatrix[gr.lineIndex]) << 3) | gr.core.charBase | gr.rowCounter // Text
		}
		if (gr.core.cr1 & 0x40) != 0 {
			addr &= 0xf9ff // ECM
		}
		gr.gfxData = gr.core.ReadByte(addr)
		gr.charData = gr.videoMatrix[gr.lineIndex]
		gr.colorData = gr.colorLine[gr.lineIndex]
		if gr.rowCounter == 0 {
			index := columnsMax * (gr.core.rasterY / 8)
			//https://sta.c64.org/cbm64scr.html
			gr.textBuffer[index+uint16(gr.lineIndex)] = _scCodesAscii[gr.charData&0x7f]
		}
		gr.lineIndex++
		gr.videoCounter++
	} else {
		if (gr.core.cr1 & 0x40) != 0 {
			gr.gfxData = gr.core.ReadByte(0x39ff)
		} else {
			gr.gfxData = gr.core.ReadByte(0x3fff)
		}
		gr.colorData = 0
		gr.charData = 0
	}
}

func (gr *Graphics) TryVideoMatrixAccess() {
	if gr.core.baLow {
		if gr.core.aecLow {
			addr := (gr.videoCounter & 0x03ff) | gr.core.matrixBase
			gr.videoMatrix[gr.lineIndex] = gr.core.ReadByte(addr)
			gr.colorLine[gr.lineIndex] = gr.core.banks.ReadColor(addr & 0x03ff)
		} else {
			gr.colorLine[gr.lineIndex] = 0xff
			gr.videoMatrix[gr.lineIndex] = 0xff
		}
	}
}

func (gr *Graphics) DrawBackground() {
	switch gr.core.displayIdx {
	case modeTextStandard, modeTextMulticolor, modeBitmapMulticolor:
		gr.db.SetMulti8(gr.offset, gr.core.b0cColor)
	case modeBitmapStandard:
		gr.db.SetMulti8(gr.offset, gr.core.colors[gr.charDataLast])
	case modeTextECM:
		if (gr.charDataLast & 0x80) != 0 {
			if (gr.charDataLast & 0x40) != 0 {
				gr.db.SetMulti8(gr.offset, gr.core.b3cColor)
			} else {
				gr.db.SetMulti8(gr.offset, gr.core.b2cColor)
			}
		} else {
			if (gr.charDataLast & 0x40) != 0 {
				gr.db.SetMulti8(gr.offset, gr.core.b1cColor)
			} else {
				gr.db.SetMulti8(gr.offset, gr.core.b0cColor)
			}
		}
	default:
		gr.db.SetMulti8(gr.offset, gr.core.colors[0])
	}

	gr.incrementOffset()
}

func (gr *Graphics) DrawForeground() {
	offset := gr.offset + int(gr.core.xScroll)
	switch gr.core.displayIdx {
	case modeTextStandard:
		gr.drawStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
	case modeTextMulticolor:
		if (gr.colorData & 8) != 0 {
			gr.drawMulticolor(offset, gr.core.b0cColor, gr.core.b1cColor, gr.core.b2cColor, gr.core.colors[gr.colorData&7])
		} else {
			gr.drawStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
		}
	case modeBitmapStandard:
		gr.drawStandard(offset, gr.core.colors[gr.charData], gr.core.colors[gr.charData>>4])
	case modeBitmapMulticolor:
		gr.drawMulticolor(offset, gr.core.b0cColor, gr.core.colors[gr.charData>>4], gr.core.colors[gr.charData], gr.core.colors[gr.colorData])
	case modeTextECM:
		if (gr.charData & 0x80) != 0 {
			if (gr.charData & 0x40) != 0 {
				gr.drawStandard(offset, gr.core.b3cColor, gr.core.colors[gr.colorData])
			} else {
				gr.drawStandard(offset, gr.core.b2cColor, gr.core.colors[gr.colorData])
			}
		} else if (gr.charData & 0x40) != 0 {
			gr.drawStandard(offset, gr.core.b1cColor, gr.core.colors[gr.colorData])
		} else {
			gr.drawStandard(offset, gr.core.b0cColor, gr.core.colors[gr.colorData])
		}
	case modeTextMulticolorInvalid:
		if (gr.colorData & 8) != 0 {
			gr.drawInvalidMulticolor(offset, gr.core.colors[0])
		} else {
			gr.drawInvalidStandard(offset, gr.core.colors[0])
		}
	case modeBitmapStandardInvalid:
		gr.drawInvalidStandard(offset, gr.core.colors[0])
	case modeBitmapMulticolorInvalid:
		gr.drawInvalidMulticolor(offset, gr.core.colors[0])
	}

	gr.incrementOffset()
}

func (gr *Graphics) incrementOffset() {
	gr.offset += 8
	gr.collisions.IncrementGraphicsOffset()
}

func (gr *Graphics) drawInvalidStandard(offset int, a uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, a)
}

func (gr *Graphics) drawInvalidMulticolor(offset int, a uint8) {
	p := (gr.gfxData & 0xAA) | ((gr.gfxData & 0xAA) >> 1)
	p1 := p >> gr.core.xScroll
	p2 := p << (8 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, a)
}

func (gr *Graphics) drawStandard(offset int, a uint8, b uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{a, b, 0, 0}
	data := gr.gfxData
	gr.db.Set(offset+7, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+6, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+5, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+4, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+3, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+2, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset+1, colorBuffer[data&1])
	data >>= 1
	gr.db.Set(offset, colorBuffer[data])
}

func (gr *Graphics) drawMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
	p := (gr.gfxData & 0xAA) | ((gr.gfxData & 0xAA) >> 1)
	p1 := p >> gr.core.xScroll
	p2 := p << (8 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{a, b, c, d}
	data := gr.gfxData
	color := colorBuffer[data&3]
	gr.db.Set(offset+7, color)
	gr.db.Set(offset+6, color)
	data >>= 2
	color = colorBuffer[data&3]
	gr.db.Set(offset+5, color)
	gr.db.Set(offset+4, color)
	data >>= 2
	color = colorBuffer[data&3]
	gr.db.Set(offset+3, color)
	gr.db.Set(offset+2, color)
	data >>= 2
	color = colorBuffer[data]
	gr.db.Set(offset+1, color)
	gr.db.Set(offset, color)
}
