package mos6569

import (
	"fmt"

	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	// columnsMax defines the maximum number of columns used for graphical or text-based data buffers in the bgrState rendering.
	// This constant represents the horizontal resolution in character units (40 columns).
	columnsMax = 40
	// rowsMax is the maximum number of rows used for video display handling and row counter operations in the bgrState logic.
	// This constant represents the vertical resolution in character units (usually 8 pixel rows per character).
	rowsMax = 7

	memoryUnitSize = 1 << 8
)

const (
	characterPixelWidth = 8
)

// GraphicsUnit represents the core structure handling graphical rendering and related operations in the system.
// It includes components for managing video memory, collisions, display buffer, and other graphical parameters.
// This struct encapsulates the state and behavior necessary for emulating the VIC-II's bgrState rendering process.
type GraphicsUnit struct {
	*component.BaseComponent
	reflect               *GraphicsUnitReflect
	memory                *MemoryUnit
	collisions            *CollisionsUnit
	beam                  *RasterBeam
	gfxData               uint8
	colorData             uint8
	charData              uint8
	charDataLatch         uint8
	ecmBackgroundColor    uint8
	ecmForegroundColor    uint8
	baseOffset            int     // Offset from bitmap sprPresence
	lineIndex             int     // Index in video matrix / color line
	videoMatrix           []uint8 // Video matrix sprPresence
	colorLine             []uint8 // Color line sprPresence
	rowCounter            uint16  // Row counter
	videoCounter          uint16  // Video counter
	videoCounterLatch     uint16  // Video counter base
	displayAccess         uint8   // Display state
	textBuffer            []byte
	xScroll               uint16 // X scroll value
	yScroll               uint16 // Y scroll value
	displayMode           uint8  // Index of current display mode
	bmm                   bool
	ecm                   bool
	b0c                   uint8 // VIC register - bgrState
	b1c                   uint8 // VIC register - bgrState
	b2c                   uint8 // VIC register - bgrState
	b3c                   uint8 // VIC register - bgrState
	badLineEnabler        bool  // Bad Lines enabled for this frame
	badLine               bool  // Current line is bad line
	sequencerFirstDmaLine uint16
	sequencerLastDmaLine  uint16
	foregroundSequencer   []func(int)
	backgroundSequencer   []func(int)
	memoryRead            [memoryUnitSize]func(uint16)
}

// NewGraphics initializes and returns a new GraphicsUnit instance with the provided VIC core, collision handler, and display buffer.
func NewGraphics(parent references.IComponent, factory references.IComponentFactory, label string, instance int, memory *MemoryUnit, collisions *CollisionsUnit, beam *RasterBeam, rasterYMax uint16, sequencerFirstDmaLine uint16, sequencerLastDmaLine uint16) *GraphicsUnit {
	gr := &GraphicsUnit{
		BaseComponent:         component.NewBaseComponent(),
		memory:                memory,
		collisions:            collisions,
		beam:                  beam,
		sequencerFirstDmaLine: sequencerFirstDmaLine,
		sequencerLastDmaLine:  sequencerLastDmaLine,
		gfxData:               0,
		colorData:             0,
		charData:              0,
		charDataLatch:         0,
		ecmBackgroundColor:    0,
		ecmForegroundColor:    0,
		baseOffset:            0,
		lineIndex:             0,
		xScroll:               0,
		yScroll:               0,
		displayMode:           0,
		b0c:                   0,
		b1c:                   0,
		b2c:                   0,
		b3c:                   0,
		videoMatrix:           make([]uint8, columnsMax),
		colorLine:             make([]uint8, columnsMax),
		textBuffer:            make([]uint8, (rasterYMax/8)*columnsMax),
		rowCounter:            rowsMax,
		videoCounter:          0,
		videoCounterLatch:     0,
		displayAccess:         0,
		bmm:                   false,
		ecm:                   false,
		badLine:               false,
		badLineEnabler:        false,
	}
	for x := range gr.memoryRead {
		gr.memoryRead[x] = gr.readGraphics
	}
	gr.memoryRead[0] = gr.readGraphicsFake

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
		gr.drawBackgroundInvalid,
		gr.drawBackgroundInvalid,
		gr.drawBackgroundInvalid,
	}
	//gr.BaseComponent.Register(factory, parent, "graphicsUnit", instance, gr, references.IdInternalComponent(label, instance, "GraphicsUnit"))
	gr.reflect = NewGraphicsUnitReflect(gr, factory, parent, "graphicsUnit", instance, references.IdInternalComponent(label, instance, "GraphicsUnit"))
	return gr
}

// Connect establishes the necessary connections or dependencies for the GraphicsUnit component to function properly.
// Returns an error if the initialization fails.
func (gr *GraphicsUnit) Connect() error {
	return nil
}

// EmulationRequired determines if the current bgrState configuration requires emulation for functionality.
func (gr *GraphicsUnit) EmulationRequired() bool {
	return false
}

// Emulate executes the main bgrState rendering loop, processing video memory, updating counters, and rendering components.
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

// BadLine checks if the bgrState unit is in a bad line state and returns true if it is, otherwise false.
func (gr *GraphicsUnit) BadLine() bool {
	return gr.badLine
}

// BadLineVerify updates the bad line condition based on the current raster position, DEN bit, and YSCROLL value.
// The bad line condition occurs when specific raster and scroll conditions are met, enabling certain VIC behavior.
// Bad Line Condition is given at any arbitrary clock cycle, if at the
// negative edge of ø0 at the beginning of the cycle RASTER >= $30 and RASTER <= $f7
// and the lower three bits of RASTER are equal to YSCROLL,
// and if the DEN bit has been set for at least one cycle somewhere in raster line $30
// So clearing the DEN bit will normally prevent Bad Lines
func (gr *GraphicsUnit) BadLineVerify(rasterY uint16, denBit bool) {
	if (rasterY >= gr.sequencerFirstDmaLine) && (rasterY <= gr.sequencerLastDmaLine) {
		if rasterY == gr.sequencerFirstDmaLine && denBit {
			//If YSCROLL=0, a Bad Line Condition occurs in raster line $30 as soon as the DEN bit
			gr.badLineEnabler = true
			if gr.yScroll == 0 {
				gr.badLine = true
				return
			}
		}
		if gr.badLineEnabler {
			gr.badLine = gr.yScroll == (rasterY & 7)
		}
	} else {
		gr.badLineEnabler = false
		gr.badLine = false
	}
}

// ReadB0c returns the value of b0c with the high nibble set to 1 (binary 1111), effectively OR-ing the value with 0xf0.
func (gr *GraphicsUnit) ReadB0c() uint8 {
	return gr.b0c | 0xf0
}

// ReadB1c retrieves the `b1c` value from the GraphicsUnit struct, applies a bitwise OR with 0xf0, and returns the result.
func (gr *GraphicsUnit) ReadB1c() uint8 {
	return gr.b1c | 0xf0
}

// ReadB2c retrieves the value of the b2c property with a bitwise OR operation applied, returning uint8 result.
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

// SetXScroll sets the horizontal scroll offset for the bgrState rendering system to the specified value.
func (gr *GraphicsUnit) SetXScroll(xScroll uint16) {
	gr.xScroll = xScroll
}

// GetXScroll retrieves the horizontal scroll offset of the bgrState rendering system.
// Returns the current value of `xScroll`.
func (gr *GraphicsUnit) GetXScroll() uint16 {
	return gr.xScroll
}

// SetYScroll sets the vertical scroll offset for the bgrState rendering system to the specified value.
func (gr *GraphicsUnit) SetYScroll(yScroll uint16) {
	gr.yScroll = yScroll
}

// GetYScroll retrieves the vertical scroll offset of the bgrState rendering system.
// Returns the current value of `yScroll`.
func (gr *GraphicsUnit) GetYScroll() uint16 {
	return gr.yScroll
}

// SetDisplayMode sets the current graphical display mode for the GraphicsUnit instance to the specified integer value.
func (gr *GraphicsUnit) SetDisplayMode(displayMode uint8) {
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

// ResetOffset resets the base offset of the GraphicsUnit to its default value of 0.
func (gr *GraphicsUnit) ResetOffset() {
	gr.baseOffset = 0
}

// CommitCharData processes the character data latch and updates the ECM background color based on bits 7 and 6 of the latch.
func (gr *GraphicsUnit) CommitCharData() {
	gr.charDataLatch = gr.charData
	ecmBackgroundMode := gr.charDataLatch >> 6 //bit 7 and 6
	switch ecmBackgroundMode {
	case 0b11:
		gr.ecmBackgroundColor = gr.b3c // Background color 3.
	case 0b10:
		gr.ecmBackgroundColor = gr.b2c // Background color 2.
	case 0b01:
		gr.ecmBackgroundColor = gr.b1c // Background color 1.
	case 0b00:
		gr.ecmBackgroundColor = gr.b0c // Background color 0.
	}
}

// AcquireDisplayAccessIfBadLine sets the displayAccess flag to true if the badLine flag in the core is active.
// This gives the CPU access to video memory during "bad lines".
func (gr *GraphicsUnit) AcquireDisplayAccessIfBadLine() {
	if gr.badLine {
		gr.displayAccess = 1
	}
}

// ResetRowCounterIfBadLine resets the row counter (RC) to 0 if the badLine in the core is true.
func (gr *GraphicsUnit) ResetRowCounterIfBadLine() {
	if gr.badLine {
		gr.rowCounter = 0
	}
}

// TryAcquireDisplayAccessOnScanlineEnd updates the display access state and row counter based on the current conditions and row limit.
// This function manages the timing of when the display access is granted to the CPU and when the row counter is incremented.
// It's called at the *end* of each scanline.
func (gr *GraphicsUnit) TryAcquireDisplayAccessOnScanlineEnd() {
	if gr.rowCounter >= rowsMax {
		// If we've reached the end of a character row (8 pixel rows), latch the video counter.
		gr.videoCounterLatch = gr.videoCounter
		gr.displayAccess = 0
	}
	if gr.badLine || gr.displayAccess != 0 {
		gr.rowCounter = (gr.rowCounter + 1) & rowsMax
		gr.displayAccess = 1
	}
}

// FetchMemory retrieves memory for a given raster line (rasterY) using the display access method.
func (gr *GraphicsUnit) FetchMemory(rasterY uint16) {
	gr.memoryRead[gr.displayAccess](rasterY)
}

// FetchData reads a character code from the video matrix at a calculated address based on videoCounter and base offset.
func (gr *GraphicsUnit) FetchData(base uint16) {
	gr.videoMatrix[gr.lineIndex] = gr.memory.ReadByte((gr.videoCounter & 0x3ff) | base)
}

// FetchDataFake fills the video matrix with fake data when AEC is high, simulating CPU access to the address bus.
func (gr *GraphicsUnit) FetchDataFake() {
	gr.videoMatrix[gr.lineIndex] = 0xff
}

// FetchColor reads color data from color RAM for the current video address and stores it in the active color line buffer.
func (gr *GraphicsUnit) FetchColor(base uint16) {
	gr.colorLine[gr.lineIndex] = gr.memory.readColorRam((gr.videoCounter & 0x3ff) | base)
}

// FetchColorFake fills the color line buffer with a fake data value when AEC is high and CPU has address bus access.
func (gr *GraphicsUnit) FetchColorFake() {
	gr.colorLine[gr.lineIndex] = 0xff
}

// DrawBackground renders the background using the current display mode
// and updates the bgrState offset and collision state.
func (gr *GraphicsUnit) DrawBackground() {
	gr.backgroundSequencer[gr.displayMode](gr.baseOffset)
	gr.baseOffset += characterPixelWidth
	gr.collisions.IncrementBackgroundOffset()
}

// DrawForeground renders the foreground bgrState based on the current display mode and x-scroll offset.
// It also increments the offset and updates the collisions.
func (gr *GraphicsUnit) DrawForeground() {
	gr.foregroundSequencer[gr.displayMode](gr.baseOffset + int(gr.xScroll))
	gr.baseOffset += characterPixelWidth
	gr.collisions.IncrementBackgroundOffset()
}

// setCharData sets the character data and determines the foreground color in Extended Color Mode based on character code bits 7 and 6.
func (gr *GraphicsUnit) setCharData(data uint8) {
	gr.charData = data
	ecmForegroundMode := gr.charData >> 6 //bit 7 and 6
	switch ecmForegroundMode {
	case 0b11:
		gr.ecmForegroundColor = gr.b3c // Foreground color 3.
	case 0b10:
		gr.ecmForegroundColor = gr.b2c // Foreground color 2.
	case 0b01:
		gr.ecmForegroundColor = gr.b1c // Foreground color 1.
	case 0b00:
		gr.ecmForegroundColor = gr.b0c // Foreground color 0.
	}
}

// setColorData sets the color data for the GraphicsUnit object using the provided 8-bit unsigned integer value.
func (gr *GraphicsUnit) setColorData(data uint8) {
	gr.colorData = data
}

// readGraphics reads bgrState data from memory based on the current rendering mode and updates internal bgrState states.
func (gr *GraphicsUnit) readGraphics(rasterY uint16) {
	var addr uint16
	if gr.bmm {
		// Bitmap Mode: Calculate the address based on the video counter, bitmap base address, and row counter.
		addr = ((gr.videoCounter & 0x03ff) << 3) | gr.memory.GetBitmapBase() | gr.rowCounter
	} else {
		// Text Mode: Calculate the address based on the character code from the video matrix, character base address, and row counter.
		addr = (uint16(gr.videoMatrix[gr.lineIndex]) << 3) | gr.memory.GetCharBase() | gr.rowCounter
	}
	if gr.ecm {
		// Extended Color Mode (ECM): Mask the address to use only the lower 13 bits of the character ROM address.
		addr &= 0xf9ff
	}
	gr.gfxData = gr.memory.ReadByte(addr)    // Operand the bgrState data (pixel data or character pattern).
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
}

// readGraphicsFake performs a fake read operation on memory, emulating data fetch based on ECM mode and resetting bgrState data.
func (gr *GraphicsUnit) readGraphicsFake( /*rasterY*/ _ uint16) {
	if gr.ecm {
		gr.gfxData = gr.memory.ReadByte(0x39ff) // Fake read (ECM).
	} else {
		gr.gfxData = gr.memory.ReadByte(0x3fff) // Fake read (non-ECM).
	}
	gr.setColorData(0) // Fake values.
	gr.setCharData(0)  // Fake values.
}

// drawBackgroundTextStandard renders the background text using the standard text mode based on the current GraphicsUnit settings.
// It uses the offset and core attributes of the GraphicsUnit instance to determine the drawing configuration.
func (gr *GraphicsUnit) drawBackgroundTextStandard(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundTextMulticolor renders a multicolor text background for the given GraphicsUnit object.
// Internally, it uses the _drawDefault function to set multicolor data based on the specified parameters.
// The GraphicsUnit parameter contains all the necessary data like offset, core state, and color information.
func (gr *GraphicsUnit) drawBackgroundTextMulticolor(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundBitmapMulticolor draws a multicolor bitmap background using the provided GraphicsUnit object.
// Updates the display buffer with colors based on the provided GraphicsUnit state and configuration.
func (gr *GraphicsUnit) drawBackgroundBitmapMulticolor(offset int) {
	gr.drawDefault(offset, gr.b0c)
}

// drawBackgroundBitmapStandard renders a standard bitmap background using the provided GraphicsUnit instance.
// It uses the _drawDefault function to handle the drawing process based on the current GraphicsUnit state.
// In standard bitmap mode, the background color for each 8x8 block is taken from the *previous* character's data.
func (gr *GraphicsUnit) drawBackgroundBitmapStandard(offset int) {
	gr.drawDefault(offset, gr.charDataLatch)
}

// drawBackgroundTextECM renders the background text in ECM (Extended Color Mode) based on character data and bitmask checks.
// It selects the appropriate color source from the GraphicsUnit core and applies it using the _drawDefault helper function.
func (gr *GraphicsUnit) drawBackgroundTextECM(offset int) {
	gr.drawDefault(offset, gr.ecmBackgroundColor)
}

// drawBackgroundDefault draws the default background by delegating to the _drawDefault function with the current offset.
// Draw 8 pixels with color 0 (usually black).
func (gr *GraphicsUnit) drawBackgroundInvalid(offset int) {
	gr.drawDefault(offset, 0)
}

// drawForegroundTextStandard renders foreground text using the standard bgrState mode at the specified offset.
func (gr *GraphicsUnit) drawForegroundTextStandard(offset int) {
	gr.drawStandard(offset, gr.b0c, gr.colorData)
}

// drawForegroundTextMulticolor renders multicolor text for the foreground depending on the color mode and provided offset.
// If the color mode indicates multicolor, it invokes `drawMulticolor`; otherwise, `drawStandard` is used.
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

// drawForegroundBitmapMulticolor renders a foreground bitmap in multicolor mode using the specified bgrState and offset.
func (gr *GraphicsUnit) drawForegroundBitmapMulticolor(offset int) {
	gr.drawMulticolor(offset, gr.b0c, gr.charData>>4, gr.charData, gr.colorData)
}

// drawForegroundTextECM handles rendering foreground text in Extended Color Mode (ECM) by selecting the correct color mapping.
func (gr *GraphicsUnit) drawForegroundTextECM(offset int) {
	gr.drawStandard(offset, gr.ecmForegroundColor, gr.colorData)
}

// drawForegroundBitmapStandardInvalid renders an invalid-standard-mode bitmap to the specified offset in the bgrState buffer.
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
	gr.beam.Draw8(offset, a)
}

// drawInvalidStandard updates bgrState buffer based on x-scroll and sets color values in the display buffer.
func (gr *GraphicsUnit) drawInvalidStandard(offset int, a uint8) {
	p1 := gr.gfxData >> gr.xScroll
	p2 := gr.gfxData << (7 - gr.xScroll)
	gr.collisions.UpdateBackground(p1, p2)
	gr.beam.Draw8(offset, a)
}

// drawInvalidMulticolor processes invalid multicolor bgrState and updates collision and display buffers accordingly.
func (gr *GraphicsUnit) drawInvalidMulticolor(offset int, a uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.xScroll
	p2 := p << (8 - gr.xScroll)
	gr.collisions.UpdateBackground(p1, p2)
	gr.beam.Draw8(offset, a)
}

// drawStandard renders 8 pixels in standard mode (1 bit per pixel). Uses colors 'a' (for 0 bits) and 'b' (for 1 bit).
func (gr *GraphicsUnit) drawStandard(offset int, a uint8, b uint8) {
	p1 := gr.gfxData >> gr.xScroll
	p2 := gr.gfxData << (7 - gr.xScroll)
	gr.collisions.UpdateBackground(p1, p2)
	gr.beam.DrawStandard(offset, a, b, gr.gfxData)
}

// drawMulticolor renders 8 pixels in multicolor mode (2 bits per pixel). Uses colors 'a', 'b', 'c', and 'd'.
func (gr *GraphicsUnit) drawMulticolor(offset int, a uint8, b uint8, c uint8, d uint8) {
	p := (gr.gfxData & 0xaa) | ((gr.gfxData & 0xaa) >> 1)
	p1 := p >> gr.xScroll
	p2 := p << (8 - gr.xScroll)
	gr.collisions.UpdateBackground(p1, p2)
	gr.beam.DrawMultiColor(offset, a, b, c, d, gr.gfxData)
}
