package mos6569

// _foregroundSequencer provides a sequence of rendering functions for various foreground drawing modes in the Graphics system.
// Each function in this slice corresponds to a different VIC-II display mode.
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
// Each function in this slice corresponds to a different VIC-II display mode.
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

// drawBackgroundTextStandard renders the background text using the standard text mode based on the current Graphics settings.
// It uses the offset and core attributes of the Graphics instance to determine the drawing configuration.
// Used in Standard Character Mode.
func drawBackgroundTextStandard(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundTextMulticolor renders a multicolor text background for the given Graphics object.
// Internally, it uses the _drawDefault function to set multi-color data based on the specified parameters.
// The Graphics parameter contains all the necessary data like offset, core state, and color information.
// Used in Multicolor Mode.
func drawBackgroundTextMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundBitmapMulticolor draws a multicolor bitmap background using the provided Graphics object.
// Updates the display buffer with colors based on the provided Graphics state and configuration.
// Used in Multicolor Bitmap Mode.
func drawBackgroundBitmapMulticolor(gr *Graphics) {
	_drawDefault(gr, gr.offset, gr.core.b0c)
}

// drawBackgroundBitmapStandard renders a standard bitmap background using the provided Graphics instance.
// It uses the _drawDefault function to handle the drawing process based on the current Graphics state.
// Used in Standard Bitmap Mode.
func drawBackgroundBitmapStandard(gr *Graphics) {
	// In standard bitmap mode, the background color for each 8x8 block is taken from the *previous* character's data.
	_drawDefault(gr, gr.offset, gr.charDataLast)
}

// drawBackgroundTextECM renders the background text in ECM (Extended Color Mode) based on character data and bitmask checks.
// It selects the appropriate color source from the Graphics core and applies it using the _drawDefault helper function.
// Used in Extended Background Color Mode (ECM).
func drawBackgroundTextECM(gr *Graphics) {
	// In ECM, the background color is determined by bits 7 and 6 of the *previous* character code.
	if (gr.charDataLast & 0x80) != 0 {
		if (gr.charDataLast & 0x40) != 0 {
			// Background color 3.
			_drawDefault(gr, gr.offset, gr.core.b3c)
		} else {
			// Background color 2.
			_drawDefault(gr, gr.offset, gr.core.b2c)
		}
	} else {
		if (gr.charDataLast & 0x40) != 0 {
			// Background color 1.
			_drawDefault(gr, gr.offset, gr.core.b1c)
		} else {
			// Background color 0.
			_drawDefault(gr, gr.offset, gr.core.b0c)
		}
	}
}

// drawBackgroundDefault draws the default background by delegating to the _drawDefault function with the current offset.
// This is used for invalid modes.
func drawBackgroundDefault(gr *Graphics) {
	// Draw 8 pixels with color 0 (usually black).
	_drawDefault(gr, gr.offset, 0)
}

// drawForegroundTextStandard renders foreground text using the standard graphics mode at the specified offset.
// Used in Standard Character Mode.
func drawForegroundTextStandard(gr *Graphics, offset int) {
	_drawStandard(gr, offset, gr.core.b0c, gr.colorData)
}

// drawForegroundTextMulticolor renders multicolor text for the foreground depending on the color mode and provided offset.
// If the color mode indicates multicolor, it invokes `_drawMulticolor`; otherwise, `_drawStandard` is used.
// Used in Multicolor Mode.
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
	gr.setMulti8(offset, _colors[a])
}

// _drawInvalidStandard updates graphics buffer based on x-scroll and sets color values in the display buffer.
func _drawInvalidStandard(gr *Graphics, offset int, a uint8) {
	p1 := gr.gfxData >> gr.xScroll
	p2 := gr.gfxData << (7 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)
	gr.setMulti8(offset, _colors[a])
}

// _drawInvalidMulticolor processes invalid multicolor graphics and updates collision and display buffers accordingly.
func _drawInvalidMulticolor(gr *Graphics, offset int, a uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.xScroll
	p2 := p << (8 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)
	gr.setMulti8(offset, _colors[a])
}

// _drawStandard renders 8 pixels in standard mode (1 bit per pixel).
// Uses colors 'a' (for 0 bits) and 'b' (for 1 bit).
func _drawStandard(gr *Graphics, offset int, a uint8, b uint8) {
	p1 := gr.gfxData >> gr.xScroll
	p2 := gr.gfxData << (7 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{_colors[a], _colors[b], 0, 0}
	index := _standardIndex[gr.gfxData]
	drawBuffer := [8]uint8{
		colorBuffer[index[0]], colorBuffer[index[1]], colorBuffer[index[2]], colorBuffer[index[3]],
		colorBuffer[index[4]], colorBuffer[index[5]], colorBuffer[index[6]], colorBuffer[index[7]],
	}
	gr.set8(offset, &drawBuffer)
}

// _drawMulticolor renders 8 pixels in multicolor mode (2 bits per pixel).
// Uses colors 'a', 'b', 'c', and 'd'.
func _drawMulticolor(gr *Graphics, offset int, a uint8, b uint8, c uint8, d uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.xScroll
	p2 := p << (8 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)

	colorBuffer := [4]uint8{_colors[a], _colors[b], _colors[c], _colors[d]}
	index := _multicolorIndex[gr.gfxData]
	drawBuffer := [8]uint8{
		colorBuffer[index[0]], colorBuffer[index[1]], colorBuffer[index[2]], colorBuffer[index[3]],
		colorBuffer[index[4]], colorBuffer[index[5]], colorBuffer[index[6]], colorBuffer[index[7]],
	}
	gr.set8(offset, &drawBuffer)
}
