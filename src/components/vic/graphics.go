package mos6569

import "fmt"

const columnsMax = 40
const rowsMax = 7

var _foregroundSequencer = []func(*Graphics, int){
	drawForegroundTextStandard,
	drawForegroundTextMulticolor,
	drawForegroundBitmapStandard,
	drawForegroundBitmapMulticolor,
	drawForegroundTextECM,
	drawForegroundTextMulticolorInvalid,
	drawForegroundBitmapStandardInvalid,
	drawForegroundBitmapMulticolorInvalid,
}

var _backgroundSequencer = []func(*Graphics){
	drawBackgroundTextStandard,
	drawBackgroundTextMulticolor,
	drawBackgroundBitmapStandard,
	drawBackgroundBitmapMulticolor,
	drawBackgroundTextECM,
	drawBackgroundDefault,
	drawBackgroundDefault,
	drawBackgroundDefault,
}

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
	//foregroundSequencer []func(*Graphics, int)
	//backgroundSequencer []func(*Graphics)
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
	//gr.foregroundSequencer = make([]func(*Graphics, int), 8)
	//gr.foregroundSequencer[modeTextStandard] = drawForegroundTextStandard
	//gr.foregroundSequencer[modeTextMulticolor] = drawForegroundTextMulticolor
	//gr.foregroundSequencer[modeBitmapStandard] = drawForegroundBitmapStandard
	//gr.foregroundSequencer[modeBitmapMulticolor] = drawForegroundBitmapMulticolor
	//gr.foregroundSequencer[modeTextECM] = drawForegroundTextECM
	//gr.foregroundSequencer[modeTextMulticolorInvalid] = drawForegroundTextMulticolorInvalid
	//gr.foregroundSequencer[modeBitmapStandardInvalid] = drawForegroundBitmapStandardInvalid
	//gr.foregroundSequencer[modeBitmapMulticolorInvalid] = drawForegroundBitmapMulticolorInvalid

	//gr.backgroundSequencer = make([]func(*Graphics), 8)
	//gr.backgroundSequencer[modeTextStandard] = drawBackgroundTextStandard
	//gr.backgroundSequencer[modeTextMulticolor] = drawBackgroundTextMulticolor
	//gr.backgroundSequencer[modeBitmapStandard] = drawBackgroundBitmapStandard
	//gr.backgroundSequencer[modeBitmapMulticolor] = drawBackgroundBitmapMulticolor
	//gr.backgroundSequencer[modeTextECM] = drawBackgroundTextECM
	//gr.backgroundSequencer[modeTextMulticolorInvalid] = drawBackgroundDefault
	//gr.backgroundSequencer[modeBitmapStandardInvalid] = drawBackgroundDefault
	//gr.backgroundSequencer[modeBitmapMulticolorInvalid] = drawBackgroundDefault

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
		if gr.core.bmm {
			addr = ((gr.videoCounter & 0x03ff) << 3) | gr.core.bitmapBase | gr.rowCounter // Bitmap
		} else {
			addr = (uint16(gr.videoMatrix[gr.lineIndex]) << 3) | gr.core.charBase | gr.rowCounter // Text
		}
		if gr.core.ecm {
			addr &= 0xf9ff
		}
		gr.gfxData = gr.core.ReadByte(addr)
		gr.charData = gr.videoMatrix[gr.lineIndex]
		gr.colorData = gr.colorLine[gr.lineIndex]
		if gr.rowCounter == 0 {
			//https://sta.c64.org/cbm64scr.html
			index := columnsMax * (gr.core.rasterY / 8)
			gr.textBuffer[index+uint16(gr.lineIndex)] = _scCodesAscii[gr.charData&0x7f]
		}
		gr.lineIndex++
		gr.videoCounter++
	} else {
		if gr.core.ecm {
			gr.gfxData = gr.core.ReadByte(0x39ff)
		} else {
			gr.gfxData = gr.core.ReadByte(0x3fff)
		}
		gr.colorData = 0
		gr.charData = 0
	}
}

func (gr *Graphics) TryPhi2Access() {
	if gr.core.baLow {
		if gr.core.aecLow {
			addr := (gr.videoCounter & 0x3ff) | gr.core.matrixBase
			gr.videoMatrix[gr.lineIndex] = gr.core.ReadByte(addr)
			gr.colorLine[gr.lineIndex] = gr.core.banks.ReadColor(addr & 0x3ff)
		} else {
			gr.colorLine[gr.lineIndex] = 0xff
			gr.videoMatrix[gr.lineIndex] = 0xff
		}
	}
}

func (gr *Graphics) DrawBackground() {
	_backgroundSequencer[gr.core.displayMode](gr)
	gr.offset += 8
	gr.collisions.IncrementGraphicsOffset()
}

func (gr *Graphics) DrawForeground() {
	offset := gr.offset + int(gr.core.xScroll)
	_foregroundSequencer[gr.core.displayMode](gr, offset)
	gr.offset += 8
	gr.collisions.IncrementGraphicsOffset()
}

func drawBackgroundTextStandard(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

func drawBackgroundTextMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

func drawBackgroundBitmapMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

func drawBackgroundBitmapStandard(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.charDataLast)
}

func drawBackgroundTextECM(gr *Graphics) {
	if (gr.charDataLast & 0x80) != 0 {
		if (gr.charDataLast & 0x40) != 0 {
			_drawDefault(gr, gr.offset, gr.core.b3c)
		} else {
			_drawDefault(gr, gr.offset, gr.core.b2c)
		}
	} else {
		if (gr.charDataLast & 0x40) != 0 {
			_drawDefault(gr, gr.offset, gr.core.b1c)
		} else {
			_drawDefault(gr, gr.offset, gr.core.b0c)
		}
	}
}

func drawBackgroundDefault(gr *Graphics) {
	_drawDefault(gr, gr.offset, 0)
}

func drawForegroundTextStandard(gr *Graphics, offset int) {
	_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
}

func drawForegroundTextMulticolor(gr *Graphics, offset int) {
	if (gr.colorData & 8) != 0 {
		_drawMulticolor(gr, offset, gr.core.b0c, gr.core.b1c, gr.core.b2c, gr.colorData&7)
	} else {
		_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
	}
}

func drawForegroundBitmapStandard(gr *Graphics, offset int) {
	_drawStandard(gr, offset, gr.charData, gr.charData>>4)
}

func drawForegroundBitmapMulticolor(gr *Graphics, offset int) {
	_drawMulticolor(gr, offset, gr.core.b0c, gr.charData>>4, gr.charData, gr.colorData)
}

func drawForegroundTextECM(gr *Graphics, offset int) {
	if (gr.charData & 0x80) != 0 {
		if (gr.charData & 0x40) != 0 {
			_drawStandard(gr, offset, gr.core.b3c, gr.colorData)
		} else {
			_drawStandard(gr, offset, gr.core.b2c, gr.colorData)
		}
	} else if (gr.charData & 0x40) != 0 {
		_drawStandard(gr, offset, gr.core.b1c, gr.colorData)
	} else {
		_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
	}
}

func drawForegroundTextMulticolorInvalid(gr *Graphics, offset int) {
	if (gr.colorData & 8) != 0 {
		_drawInvalidMulticolor(gr, offset, 0)
	} else {
		_drawInvalidStandard(gr, offset, 0)
	}
}

func drawForegroundBitmapStandardInvalid(gr *Graphics, offset int) {
	_drawInvalidStandard(gr, offset, 0)
}

func drawForegroundBitmapMulticolorInvalid(gr *Graphics, offset int) {
	_drawInvalidMulticolor(gr, offset, 0)
}

func _drawDefault(gr *Graphics, offset int, a uint8) {
	gr.db.SetMulti8(offset, _colors[a])
}

func _drawInvalidStandard(gr *Graphics, offset int, a uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, _colors[a])
}

func _drawInvalidMulticolor(gr *Graphics, offset int, a uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.core.xScroll
	p2 := p << (8 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, _colors[a])
}

func _drawStandard(gr *Graphics, offset int, a uint8, b uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{_colors[a], _colors[b], 0, 0}
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

func _drawMulticolor(gr *Graphics, offset int, a uint8, b uint8, c uint8, d uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.core.xScroll
	p2 := p << (8 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{_colors[a], _colors[b], _colors[c], _colors[d]}
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
