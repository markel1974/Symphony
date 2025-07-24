package mos6569

import (
	"fmt"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	// columnsMax defines the maximum number of columns used for graphical or text-based data buffers in the graphics rendering.
	// This constant represents the horizontal resolution in character units (40 columns).
	columnsMax = 40
	// rowsMax is the maximum number of rows used for video display handling and row counter operations in the graphics logic.
	// This constant represents the vertical resolution in character units (usually 8 pixel rows per character).
	rowsMax = 7
)

// GraphicsUnit represents the core structure handling graphical rendering and related operations in the system.
// It includes components for managing video memory, collisions, display buffer, and other graphical parameters.
// This struct encapsulates the state and behavior necessary for emulating the VIC-II's graphics rendering process.
type GraphicsUnit struct {
	*component.BaseComponent
	memory              *MemoryUnit
	collisions          *CollisionsUnit
	set8                func(int, *[8]uint8)
	setMulti8           func(int, uint8)
	gfxData             uint8
	colorData           uint8
	charData            uint8
	charDataLast        uint8
	ecmBackgroundColor  uint8
	ecmForegroundColor  uint8
	offset              int     // Offset from bitmap spritesPresence
	lineIndex           int     // Index in video matrix / color line
	videoMatrix         []uint8 // Video matrix spritesPresence
	colorLine           []uint8 // Color line spritesPresence
	rowCounter          uint16  // Row counter
	videoCounter        uint16  // Video counter
	videoCounterLatch   uint16  // Video counter base
	displayAccess       bool    // Display state
	textBuffer          []byte
	xScroll             uint16 // X scroll value
	yScroll             uint16 // Y scroll value
	displayMode         int    // Index of current display mode
	bmm                 bool
	ecm                 bool
	b0c                 uint8 // VIC register - graphics
	b1c                 uint8 // VIC register - graphics
	b2c                 uint8 // VIC register - graphics
	b3c                 uint8 // VIC register - graphics
	foregroundSequencer []func(int)
	backgroundSequencer []func(int)
}

// NewGraphics initializes and returns a new GraphicsUnit instance with the provided VIC core, collision handler, and display buffer.
func NewGraphics(parent references.IComponent, factory references.IComponentFactory, label string, instance int, memory *MemoryUnit, collisions *CollisionsUnit, displayBuffer references.IDisplayBuffer, rasterYMax uint16) *GraphicsUnit {
	gr := &GraphicsUnit{
		BaseComponent:      component.NewBaseComponent(),
		memory:             memory,
		collisions:         collisions,
		set8:               displayBuffer.Set8,
		setMulti8:          displayBuffer.SetMulti8,
		gfxData:            0,
		colorData:          0,
		charData:           0,
		charDataLast:       0,
		ecmBackgroundColor: 0,
		ecmForegroundColor: 0,
		offset:             0,
		lineIndex:          0,
		xScroll:            0,
		yScroll:            0,
		displayMode:        0,
		b0c:                0,
		b1c:                0,
		b2c:                0,
		b3c:                0,
		videoMatrix:        make([]uint8, columnsMax),
		colorLine:          make([]uint8, columnsMax),
		textBuffer:         make([]uint8, (rasterYMax/8)*columnsMax),
		rowCounter:         rowsMax,
		videoCounter:       0,
		videoCounterLatch:  0,
		displayAccess:      false,
		bmm:                false,
		ecm:                false,
	}

	// foregroundSequencer provides a sequence of rendering functions for various foreground drawing modes in the GraphicsUnit system.
	gr.foregroundSequencer = []func(int){
		gr.drawForegroundTextStandard,
		gr.drawForegroundTextMulticolor,
		gr.drawForegroundBitmapStandard,
		gr.drawForegroundBitmapMulticolor,
		gr.drawForegroundTextECM,
		gr.drawForegroundTextMulticolorInvalid,
		gr.drawForegroundBitmapStandardInvalid,
		gr.drawForegroundBitmapMulticolorInvalid,
	}

	// backgroundSequencer is a sequence of functions for rendering backgrounds based on the current display mode of GraphicsUnit.
	gr.backgroundSequencer = []func(int){
		gr.drawBackgroundTextStandard,
		gr.drawBackgroundTextMulticolor,
		gr.drawBackgroundBitmapStandard,
		gr.drawBackgroundBitmapMulticolor,
		gr.drawBackgroundTextECM,
		gr.drawBackgroundDefault,
		gr.drawBackgroundDefault,
		gr.drawBackgroundDefault,
	}
	gr.BaseComponent.Register(factory, parent, "graphics", gr, references.IdInternalComponent(label, instance, "GraphicsUnit"))
	return gr
}

// Connect establishes the necessary connections or dependencies for the GraphicsUnit component to function properly.
// Returns an error if the initialization fails.
func (gr *GraphicsUnit) Connect() error {
	return nil
}

// EmulationRequired determines if the current graphics configuration necessitates emulation for functionality.
func (gr *GraphicsUnit) EmulationRequired() bool {
	return false
}

// Emulate executes the main graphics rendering loop, processing video memory, updating counters, and rendering components.
func (gr *GraphicsUnit) Emulate() {
}

// Internal checks and returns a boolean indicating internal state or configuration for graphical operations.
func (gr *GraphicsUnit) Internal() bool {
	return true
}

// Reset reinitializes the GraphicsUnit instance to its default state, clearing any temporary data and resetting counters.
func (gr *GraphicsUnit) Reset() {
}

// Setup initializes the GraphicsUnit instance and prepares it for rendering operations.
// (Currently empty, but kept for consistency).
func (gr *GraphicsUnit) Setup() error {
	return nil
}

// ReadB0c returns the value of b0c with the high nibble set to 1 (binary 1111), effectively OR-ing the value with 0xf0.
func (gr *GraphicsUnit) ReadB0c() uint8 {
	return gr.b0c | 0xf0
}

// ReadB1c retrieves the `b1c` value from the GraphicsUnit struct, applies a bitwise OR with 0xf0, and returns the result.
func (gr *GraphicsUnit) ReadB1c() uint8 {
	return gr.b1c | 0xf0
}

// ReadB2c retrieves the value of the b2c property with a bitwise OR operation applied, returning a uint8 result.
func (gr *GraphicsUnit) ReadB2c() uint8 {
	return gr.b2c | 0xf0
}

// ReadB3c reads the b3c property of the GraphicsUnit receiver and applies a bitwise OR operation with the value 0xf0.
func (gr *GraphicsUnit) ReadB3c() uint8 {
	return gr.b3c | 0xf0
}

// WriteB0c sets the value of the b0c field in the GraphicsUnit instance to the specified data value.
func (gr *GraphicsUnit) WriteB0c(data uint8) {
	gr.b0c = data
}

// WriteB1c sets the value of the b1c field in the GraphicsUnit struct to the provided uint8 data.
func (gr *GraphicsUnit) WriteB1c(data uint8) {
	gr.b1c = data
}

// WriteB2c sets the value of the b2c field with the provided data parameter.
func (gr *GraphicsUnit) WriteB2c(data uint8) {
	gr.b2c = data
}

// WriteB3c sets the b3c property of the GraphicsUnit object to the given data value.
func (gr *GraphicsUnit) WriteB3c(data uint8) {
	gr.b3c = data
}

// SetXScroll sets the horizontal scroll offset for the graphics rendering system to the specified value.
func (gr *GraphicsUnit) SetXScroll(xScroll uint16) {
	gr.xScroll = xScroll
}

// GetXScroll retrieves the horizontal scroll offset of the graphics rendering system.
// Returns the current value of `xScroll`.
func (gr *GraphicsUnit) GetXScroll() uint16 {
	return gr.xScroll
}

// SetYScroll sets the vertical scroll offset for the graphics rendering system to the specified value.
func (gr *GraphicsUnit) SetYScroll(yScroll uint16) {
	gr.yScroll = yScroll
}

// GetYScroll retrieves the vertical scroll offset of the graphics rendering system.
// Returns the current value of `yScroll`.
func (gr *GraphicsUnit) GetYScroll() uint16 {
	return gr.yScroll
}

// SetDisplayMode sets the current graphical display mode for the GraphicsUnit instance to the specified integer value.
func (gr *GraphicsUnit) SetDisplayMode(displayMode int) {
	gr.displayMode = displayMode
}

// SetBmm sets the bmm property of the GraphicsUnit object to the specified boolean value.
func (gr *GraphicsUnit) SetBmm(bmm bool) {
	gr.bmm = bmm
}

// SetEcm sets the ECM (Error Correction Mode) state for the GraphicsUnit instance.
func (gr *GraphicsUnit) SetEcm(ecm bool) {
	gr.ecm = ecm
}

// GetText retrieves the text buffer content from the GraphicsUnit instance as a slice of bytes.
func (gr *GraphicsUnit) GetText() []byte {
	return gr.textBuffer
}

// PrintText outputs the contents of the textBuffer in a formatted manner, wrapping lines at every 40 characters.
// This is a helper function for debugging.
func (gr *GraphicsUnit) PrintText() {
	for x, v := range gr.textBuffer {
		if (x % 40) == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", v)
	}
	fmt.Println()
}

// ResetVideoCounterLatch resets the video counter latch to zero.  This happens at the start of each frame.
func (gr *GraphicsUnit) ResetVideoCounterLatch() {
	gr.videoCounterLatch = 0
}

// UpdateVideoCounter updates the video counter to match the current video counter latch value.
func (gr *GraphicsUnit) UpdateVideoCounter() {
	gr.videoCounter = gr.videoCounterLatch
}

// ResetLineIndex resets the line index to zero. This happens at the beginning of each scanline.
func (gr *GraphicsUnit) ResetLineIndex() {
	gr.lineIndex = 0
}

// SetOffset sets the offset value for the GraphicsUnit instance to the given parameter.
// The offset represents the horizontal starting position (in pixels) for rendering.
func (gr *GraphicsUnit) SetOffset(offset int) {
	gr.offset = offset
}

// UpdateCharDataLast updates the `charDataLast` field with the current value of `charData`.
// This is needed for Extended Color Mode (ECM) where the background color is determined by the *previous* character.
func (gr *GraphicsUnit) UpdateCharDataLast() {
	gr.charDataLast = gr.charData
	// In ECM, the background color is determined by bits 7 and 6 of the *previous* character code.
	if (gr.charDataLast & 0x80) != 0 {
		if (gr.charDataLast & 0x40) != 0 {
			gr.ecmBackgroundColor = gr.b3c // Background color 3.
		} else {
			gr.ecmBackgroundColor = gr.b2c // Background color 2.
		}
	} else {
		if (gr.charDataLast & 0x40) != 0 {
			gr.ecmBackgroundColor = gr.b1c // Background color 1.
		} else {
			gr.ecmBackgroundColor = gr.b0c // Background color 0.
		}
	}
}

// setCharData sets the character data and determines the foreground color in Extended Color Mode based on character code bits.
func (gr *GraphicsUnit) setCharData(data uint8) {
	gr.charData = data
	// In ECM, the foreground color is determined by bits 7 and 6 of the character code.
	if (gr.charData & 0x80) != 0 {
		if (gr.charData & 0x40) != 0 {
			gr.ecmForegroundColor = gr.b3c
		} else {
			gr.ecmForegroundColor = gr.b2c
		}
	} else if (gr.charData & 0x40) != 0 {
		gr.ecmForegroundColor = gr.b1c
	} else {
		gr.ecmForegroundColor = gr.b0c
	}
}

// setColorData sets the color data for the GraphicsUnit object using the provided 8-bit unsigned integer value.
func (gr *GraphicsUnit) setColorData(data uint8) {
	gr.colorData = data
}

// TryResetRowCounter resets the row counter (RC) to 0 if the badLineCondition in the core is true.
func (gr *GraphicsUnit) TryResetRowCounter(badLineCondition bool) {
	if badLineCondition {
		gr.rowCounter = 0
	}
}

// TryAcquireDisplayAccess sets the displayAccess flag to true if the badLineCondition flag in the core is active.
// This gives the CPU access to video memory during "bad lines".
func (gr *GraphicsUnit) TryAcquireDisplayAccess(badLineCondition bool) {
	if badLineCondition {
		gr.displayAccess = true
	}
}

// UpdateDisplayAccess updates the display access state and row counter based on the current conditions and row limit.
// This function manages the timing of when the display access is granted to the CPU and when the row counter is incremented.
// It's called at the *end* of each scanline (cycle 58).
func (gr *GraphicsUnit) UpdateDisplayAccess(badLineCondition bool) {
	if gr.rowCounter == rowsMax {
		// If we've reached the end of a character row (8 pixel rows), latch the video counter.
		gr.videoCounterLatch = gr.videoCounter
		gr.displayAccess = false
	}
	if badLineCondition || gr.displayAccess {
		// The & operator has precedence
		gr.rowCounter = (gr.rowCounter + 1) & rowsMax
		gr.displayAccess = true
	}
}

// TryGraphicsAccess fetches and processes graphics data from memory based on the current raster and display state.
// This function is the core of the VIC-II's graphics data fetching logic.
func (gr *GraphicsUnit) TryGraphicsAccess(rasterY uint16) {
	if gr.displayAccess {
		var addr uint16
		if gr.bmm {
			// Bitmap Mode: Calculate the address based on the video counter, bitmap base address, and row counter.
			addr = ((gr.videoCounter & 0x03ff) << 3) | gr.memory.GetBitmapBase() | gr.rowCounter // Bitmap
		} else {
			// Text Mode: Calculate the address based on the character code from the video matrix, character base address, and row counter.
			addr = (uint16(gr.videoMatrix[gr.lineIndex]) << 3) | gr.memory.GetCharBase() | gr.rowCounter // Text
		}
		if gr.ecm {
			// Extended Color Mode (ECM): Mask the address to use only the lower 13 bits of the character ROM address.
			addr &= 0xf9ff
		}
		gr.gfxData = gr.memory.ReadByte(addr)    // Read the graphics data (pixel data or character pattern).
		charData := gr.videoMatrix[gr.lineIndex] // Get the character code from the video matrix.
		colorData := gr.colorLine[gr.lineIndex]  // Get the color data from the color RAM.
		gr.setCharData(charData)
		gr.setColorData(colorData)
		if gr.rowCounter == 0 {
			// At the beginning of a new character row (rowCounter == 0),
			// store the character code in the text buffer for debugging/display purposes.
			// https://sta.c64.org/cbm64scr.html
			index := columnsMax * (rasterY / 8)
			gr.textBuffer[index+uint16(gr.lineIndex)] = _scCodesAscii[charData&0x7f] // Convert screen code to ASCII
		}

		gr.lineIndex++    // Increment the line index to point to the next character/color data.
		gr.videoCounter++ // Increment the video counter.
	} else {
		// If display access is *not* granted, read from a "dummy" address.  The values read are not used.
		if gr.ecm {
			gr.gfxData = gr.memory.ReadByte(0x39ff) // Fake read (ECM).
		} else {
			gr.gfxData = gr.memory.ReadByte(0x3fff) // Fake read (non-ECM).
		}
		gr.setColorData(0) // Dummy values.
		gr.setCharData(0)  // Dummy values.
	}
}

// TryPhi2Access handles Phi2 clock phase access, updating videoMatrix and colorLine based on core conditions and memory.
// This function handles the memory access during the PHI2 phase of the CPU clock cycle.
func (gr *GraphicsUnit) TryPhi2Access(baLow bool, aecLow bool) {
	// Check if the Bus-Available (BA) signal is low.
	if baLow {
		// Check if the Address Enable Control (AEC) signal is low.
		if aecLow {
			// If both BA and AEC are low, the VIC-II has access to the address bus.
			addr := (gr.videoCounter & 0x3ff) | gr.memory.GetMatrixBase() // Calculate address in video matrix.
			gr.videoMatrix[gr.lineIndex] = gr.memory.ReadByte(addr)       // Read character code from video matrix.
			gr.colorLine[gr.lineIndex] = gr.memory.readColorRam(addr)     // Read color data from color RAM.
		} else {
			// If AEC is high, the CPU has access to the address bus, so we fill with dummy data.
			gr.colorLine[gr.lineIndex] = 0xff
			gr.videoMatrix[gr.lineIndex] = 0xff
		}
	}
}

// DrawBackground renders the background using the current display mode
// and updates the graphics offset and collision state.
func (gr *GraphicsUnit) DrawBackground() {
	// Call the appropriate background drawing function based on the current display mode.
	gr.backgroundSequencer[gr.displayMode](gr.offset)
	// Increment the pixel offset by 8 (one character width).
	gr.offset += 8
	// Update the collision detection system's offset.
	gr.collisions.IncrementGraphicsOffset()
}

// DrawForeground renders the foreground graphics based on the current display mode and x-scroll offset.
// It also increments the offset and updates the collisions.
func (gr *GraphicsUnit) DrawForeground() {
	// Calculate the final offset, including x-scrolling.
	offset := gr.offset + int(gr.xScroll)
	// Call the appropriate foreground drawing function.
	gr.foregroundSequencer[gr.displayMode](offset)
	// Increment the pixel offset by 8.
	gr.offset += 8
	// Update the collision detection system's offset.
	gr.collisions.IncrementGraphicsOffset()
}

// drawBackgroundTextStandard renders the background text using the standard text mode based on the current GraphicsUnit settings.
// It uses the offset and core attributes of the GraphicsUnit instance to determine the drawing configuration.
// Used in Standard Character Mode.
func (gr *GraphicsUnit) drawBackgroundTextStandard(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundTextMulticolor renders a multicolor text background for the given GraphicsUnit object.
// Internally, it uses the _drawDefault function to set multi-color data based on the specified parameters.
// The GraphicsUnit parameter contains all the necessary data like offset, core state, and color information.
// Used in Multicolor Mode.
func (gr *GraphicsUnit) drawBackgroundTextMulticolor(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundBitmapMulticolor draws a multicolor bitmap background using the provided GraphicsUnit object.
// Updates the display buffer with colors based on the provided GraphicsUnit state and configuration.
// Used in Multicolor Bitmap Mode.
func (gr *GraphicsUnit) drawBackgroundBitmapMulticolor(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundBitmapStandard renders a standard bitmap background using the provided GraphicsUnit instance.
// It uses the _drawDefault function to handle the drawing process based on the current GraphicsUnit state.
// Used in Standard Bitmap Mode.
func (gr *GraphicsUnit) drawBackgroundBitmapStandard(offset int) {
	// In standard bitmap mode, the background color for each 8x8 block is taken from the *previous* character's data.
	gr.drawDefault(offset, gr.charDataLast)
}

// drawBackgroundTextECM renders the background text in ECM (Extended Color Mode) based on character data and bitmask checks.
// It selects the appropriate color source from the GraphicsUnit core and applies it using the _drawDefault helper function.
// Used in Extended Background Color Mode (ECM).
func (gr *GraphicsUnit) drawBackgroundTextECM(offset int) {
	gr.drawDefault(offset, gr.ecmBackgroundColor)
}

// drawBackgroundDefault draws the default background by delegating to the _drawDefault function with the current offset.
// This is used for invalid modes.
func (gr *GraphicsUnit) drawBackgroundDefault(offset int) {
	// Draw 8 pixels with color 0 (usually black).
	gr.drawDefault(offset, 0)
}

// drawForegroundTextStandard renders foreground text using the standard graphics mode at the specified offset.
// Used in Standard Character Mode.
func (gr *GraphicsUnit) drawForegroundTextStandard(offset int) {
	gr.drawStandard(offset, gr.b0c, gr.colorData)
}

// drawForegroundTextMulticolor renders multicolor text for the foreground depending on the color mode and provided offset.
// If the color mode indicates multicolor, it invokes `drawMulticolor`; otherwise, `drawStandard` is used.
// Used in Multicolor Mode.
func (gr *GraphicsUnit) drawForegroundTextMulticolor(offset int) {
	if (gr.colorData & 8) != 0 {
		gr.drawMulticolor(offset, gr.b0c, gr.b1c, gr.b2c, gr.colorData&7)
	} else {
		gr.drawStandard(offset, gr.b0c, gr.colorData)
	}
}

// drawForegroundTextMulticolorInvalid draws invalid multicolor or standard text foreground based on colorData.
func (gr *GraphicsUnit) drawForegroundTextMulticolorInvalid(offset int) {
	if (gr.colorData & 8) != 0 {
		gr.drawInvalidMulticolor(offset, 0)
	} else {
		gr.drawInvalidStandard(offset, 0)
	}
}

// drawForegroundBitmapStandard renders a standard foreground bitmap using character and offset data for pixel mapping.
func (gr *GraphicsUnit) drawForegroundBitmapStandard(offset int) {
	gr.drawStandard(offset, gr.charData, gr.charData>>4)
}

// drawForegroundBitmapMulticolor renders a foreground bitmap in multicolor mode using the specified graphics and offset.
func (gr *GraphicsUnit) drawForegroundBitmapMulticolor(offset int) {
	gr.drawMulticolor(offset, gr.b0c, gr.charData>>4, gr.charData, gr.colorData)
}

// drawForegroundTextECM handles rendering of foreground text in Extended Color Mode (ECM) by selecting the correct color mapping.
func (gr *GraphicsUnit) drawForegroundTextECM(offset int) {
	gr.drawStandard(offset, gr.ecmForegroundColor, gr.colorData)
}

// drawForegroundBitmapStandardInvalid renders an invalid-standard-mode bitmap to the specified offset in the graphics buffer.
// It calls the internal _drawInvalidStandard function to handle rendering logic with the provided GraphicsUnit object.
func (gr *GraphicsUnit) drawForegroundBitmapStandardInvalid(offset int) {
	gr.drawInvalidStandard(offset, 0)
}

// drawForegroundBitmapMulticolorInvalid renders an invalid multicolor bitmap at the specified offset using GraphicsUnit object.
func (gr *GraphicsUnit) drawForegroundBitmapMulticolorInvalid(offset int) {
	gr.drawInvalidMulticolor(offset, 0)
}

// drawDefault sets a color value from the _colors array into the display buffer at the specified offset.
func (gr *GraphicsUnit) drawDefault(offset int, a uint8) {
	gr.setMulti8(offset, _colors[a])
}

// drawInvalidStandard updates graphics buffer based on x-scroll and sets color values in the display buffer.
func (gr *GraphicsUnit) drawInvalidStandard(offset int, a uint8) {
	p1 := gr.gfxData >> gr.xScroll
	p2 := gr.gfxData << (7 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)
	gr.setMulti8(offset, _colors[a])
}

// drawInvalidMulticolor processes invalid multicolor graphics and updates collision and display buffers accordingly.
func (gr *GraphicsUnit) drawInvalidMulticolor(offset int, a uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.xScroll
	p2 := p << (8 - gr.xScroll)
	gr.collisions.UpdateGraphics(p1, p2)
	gr.setMulti8(offset, _colors[a])
}

// drawStandard renders 8 pixels in standard mode (1 bit per pixel).
// Uses colors 'a' (for 0 bits) and 'b' (for 1 bit).
func (gr *GraphicsUnit) drawStandard(offset int, a uint8, b uint8) {
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

// drawMulticolor renders 8 pixels in multicolor mode (2 bits per pixel).
// Uses colors 'a', 'b', 'c', and 'd'.
func (gr *GraphicsUnit) drawMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
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
