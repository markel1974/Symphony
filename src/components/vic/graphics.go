package mos6569

import "fmt"

// columnsMax defines the maximum number of columns used for graphical or text-based data buffers in the graphics rendering.
const columnsMax = 40

// rowsMax is the maximum number of rows used for video display handling and row counter operations in the graphics logic.
const rowsMax = 7

// _foregroundSequencer provides a sequence of rendering functions for various foreground drawing modes in the Graphics system.
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

// _backgroundSequencer is a sequence of functions for rendering backgrounds based on the current display mode of Graphics.
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

// Graphics represents the core structure handling graphical rendering and related operations in the system.
// It includes components for managing video memory, collisions, display buffer, and other graphical parameters.
type Graphics struct {
	core              *VIC
	collisions        *Collisions
	db                IDisplayBuffer
	gfxData           uint8
	colorData         uint8
	charData          uint8
	charDataLast      uint8
	offset            int     // Offset from bitmap spritesBuffer
	lineIndex         int     // Index in video matrix / color line
	videoMatrix       []uint8 // Video matrix spritesBuffer
	colorLine         []uint8 // Color line spritesBuffer
	rowCounter        uint16  // Row counter
	videoCounter      uint16  // Video counter
	videoCounterLatch uint16  // Video counter base
	displayAccess     bool    // Display state
	textBuffer        []byte
}

// NewGraphics initializes and returns a new Graphics instance with the provided VIC core, collision handler, and display buffer.
func NewGraphics(core *VIC, collisions *Collisions, db IDisplayBuffer) *Graphics {
	gr := &Graphics{
		core:              core,
		collisions:        collisions,
		db:                db,
		gfxData:           0,
		colorData:         0,
		charData:          0,
		charDataLast:      0,
		offset:            0,
		lineIndex:         0,
		videoMatrix:       make([]uint8, columnsMax),
		colorLine:         make([]uint8, columnsMax),
		textBuffer:        make([]uint8, (RasterYMax/8)*columnsMax),
		rowCounter:        rowsMax,
		videoCounter:      0,
		videoCounterLatch: 0,
		displayAccess:     false,
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

// PrintText outputs the contents of the textBuffer in a formatted manner, wrapping lines at every 40 characters.
func (gr *Graphics) PrintText() {
	for x, v := range gr.textBuffer {
		if (x % 40) == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", v)
	}
	fmt.Println()
}

// Setup initializes the Graphics instance and prepares it for rendering operations.
func (gr *Graphics) Setup() {
}

// GetText retrieves the text buffer content from the Graphics instance as a slice of bytes.
func (gr *Graphics) GetText() []byte {
	return gr.textBuffer
}

// ResetVideoCounterLatch resets the video counter latch to zero.
func (gr *Graphics) ResetVideoCounterLatch() {
	gr.videoCounterLatch = 0
}

// UpdateVideoCounter updates the video counter to match the current video counter latch value.
func (gr *Graphics) UpdateVideoCounter() {
	gr.videoCounter = gr.videoCounterLatch
}

// ResetLineIndex resets the line index to zero, typically used to start a new line in the video matrix or color line.
func (gr *Graphics) ResetLineIndex() {
	gr.lineIndex = 0
}

// SetOffset sets the offset value for the Graphics instance to the given parameter.
func (gr *Graphics) SetOffset(offset int) {
	gr.offset = offset
}

// UpdateCharDataLast updates the `charDataLast` field with the current value of `charData` to retain the last frame's data.
func (gr *Graphics) UpdateCharDataLast() {
	gr.charDataLast = gr.charData
}

// TryResetRowCounter resets the row counter to 0 if the badLineCondition in the core is true.
func (gr *Graphics) TryResetRowCounter() {
	if gr.core.badLineCondition {
		gr.rowCounter = 0
	}
}

// TryAcquireDisplayAccess sets the displayAccess flag to true if the badLineCondition flag in the core is active.
func (gr *Graphics) TryAcquireDisplayAccess() {
	if gr.core.badLineCondition {
		gr.displayAccess = true
	}
}

// UpdateDisplayAccess updates the display access state and row counter based on the current conditions and row limit.
func (gr *Graphics) UpdateDisplayAccess() {
	if gr.rowCounter == rowsMax {
		gr.videoCounterLatch = gr.videoCounter
		gr.displayAccess = false
	}
	if gr.core.badLineCondition || gr.displayAccess {
		gr.rowCounter = (gr.rowCounter + 1) & rowsMax
		gr.displayAccess = true
	}
}

// TryGraphicsAccess fetches and processes graphics data from memory based on the current raster and display state.
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

// TryPhi2Access handles Phi2 clock phase access, updating videoMatrix and colorLine based on core conditions and memory.
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

// DrawBackground renders the background using the current display mode and updates the graphics offset and collision state.
func (gr *Graphics) DrawBackground() {
	_backgroundSequencer[gr.core.displayMode](gr)
	gr.offset += 8
	gr.collisions.IncrementGraphicsOffset()
}

// DrawForeground renders the foreground graphics based on the current display mode and x-scroll offset.
// It also increments the offset and updates the collisions.
func (gr *Graphics) DrawForeground() {
	offset := gr.offset + int(gr.core.xScroll)
	_foregroundSequencer[gr.core.displayMode](gr, offset)
	gr.offset += 8
	gr.collisions.IncrementGraphicsOffset()
}

// drawBackgroundTextStandard renders the background text using the standard text mode based on the current Graphics settings.
// It utilizes the offset and core attributes of the Graphics instance to determine the drawing configuration.
func drawBackgroundTextStandard(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundTextMulticolor renders a multicolor text background for the given Graphics object.
// Internally, it utilizes the _drawDefault function to set multi-color data based on the specified parameters.
// The Graphics parameter contains all the necessary data like offset, core state, and color information.
func drawBackgroundTextMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundBitmapMulticolor draws a multicolor bitmap background using the provided Graphics object.
// Updates the display buffer with colors based on the provided Graphics state and configuration.
func drawBackgroundBitmapMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundBitmapStandard renders a standard bitmap background using the provided Graphics instance.
// It utilizes the _drawDefault function to handle the drawing process based on the current Graphics state.
func drawBackgroundBitmapStandard(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.charDataLast)
}

// drawBackgroundTextECM renders the background text in ECM (Extended Color Mode) based on character data and bitmask checks.
// It selects the appropriate color source from the Graphics core and applies it using the _drawDefault helper function.
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

// drawBackgroundDefault draws the default background by delegating to the _drawDefault function with the current offset.
func drawBackgroundDefault(gr *Graphics) {
	_drawDefault(gr, gr.offset, 0)
}

// drawForegroundTextStandard renders foreground text using the standard graphics mode at the specified offset.
func drawForegroundTextStandard(gr *Graphics, offset int) {
	_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
}

// drawForegroundTextMulticolor renders multicolor text for the foreground depending on the color mode and provided offset.
// If the color mode indicates multicolor, it invokes `_drawMulticolor`; otherwise, `_drawStandard` is used.
func drawForegroundTextMulticolor(gr *Graphics, offset int) {
	if (gr.colorData & 8) != 0 {
		_drawMulticolor(gr, offset, gr.core.b0c, gr.core.b1c, gr.core.b2c, gr.colorData&7)
	} else {
		_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
	}
}

// drawForegroundBitmapStandard renders a standard foreground bitmap using character and offset data for pixel mapping.
func drawForegroundBitmapStandard(gr *Graphics, offset int) {
	_drawStandard(gr, offset, gr.charData, gr.charData>>4)
}

// drawForegroundBitmapMulticolor renders a foreground bitmap in multicolor mode using the specified graphics and offset.
func drawForegroundBitmapMulticolor(gr *Graphics, offset int) {
	_drawMulticolor(gr, offset, gr.core.b0c, gr.charData>>4, gr.charData, gr.colorData)
}

// drawForegroundTextECM handles rendering of foreground text in Extended Color Mode (ECM) by selecting the correct color mapping.
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

// drawForegroundTextMulticolorInvalid draws invalid multicolor or standard text foreground based on colorData.
func drawForegroundTextMulticolorInvalid(gr *Graphics, offset int) {
	if (gr.colorData & 8) != 0 {
		_drawInvalidMulticolor(gr, offset, 0)
	} else {
		_drawInvalidStandard(gr, offset, 0)
	}
}

// drawForegroundBitmapStandardInvalid renders an invalid-standard-mode bitmap to the specified offset in the graphics buffer.
// It calls the internal _drawInvalidStandard function to handle rendering logic with the provided Graphics object.
func drawForegroundBitmapStandardInvalid(gr *Graphics, offset int) {
	_drawInvalidStandard(gr, offset, 0)
}

// drawForegroundBitmapMulticolorInvalid renders an invalid multicolor bitmap at the specified offset using Graphics object.
func drawForegroundBitmapMulticolorInvalid(gr *Graphics, offset int) {
	_drawInvalidMulticolor(gr, offset, 0)
}

// _drawDefault sets a color value from the _colors array into the display buffer at the specified offset.
func _drawDefault(gr *Graphics, offset int, a uint8) {
	gr.db.SetMulti8(offset, _colors[a])
}

// _drawInvalidStandard updates graphics buffer based on x-scroll and sets color values in the display buffer.
func _drawInvalidStandard(gr *Graphics, offset int, a uint8) {
	p1 := gr.gfxData >> gr.core.xScroll
	p2 := gr.gfxData << (7 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, _colors[a])
}

// _drawInvalidMulticolor processes invalid multicolor graphics and updates collision and display buffers accordingly.
func _drawInvalidMulticolor(gr *Graphics, offset int, a uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.core.xScroll
	p2 := p << (8 - gr.core.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	gr.db.SetMulti8(offset, _colors[a])
}

// _drawStandard renders 8 pixels to a graphics display using two color indices, updating graphics collisions and buffer.
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

// _drawMulticolor draws a multicolor graphics sequence into the display buffer using the provided Graphics context.
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
