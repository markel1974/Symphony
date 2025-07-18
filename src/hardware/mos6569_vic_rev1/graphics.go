package mos6569

import (
	"fmt"
	"github.com/markel1974/c64emu/src/references"
)

// columnsMax defines the maximum number of columns used for graphical or text-based data buffers in the graphics rendering.
// This constant represents the horizontal resolution in character units (40 columns).
const columnsMax = 40

// rowsMax is the maximum number of rows used for video display handling and row counter operations in the graphics logic.
// This constant represents the vertical resolution in character units (usually 8 pixel rows per character).
const rowsMax = 7

// Graphics represents the core structure handling graphical rendering and related operations in the system.
// It includes components for managing video memory, collisions, display buffer, and other graphical parameters.
// This struct encapsulates the state and behavior necessary for emulating the VIC-II's graphics rendering process.
type Graphics struct {
	core              *VIC
	collisions        *Collisions
	set8              func(int, *[8]uint8)
	setMulti8         func(int, uint8)
	gfxData           uint8
	colorData         uint8
	charData          uint8
	charDataLast      uint8
	offset            int     // Offset from bitmap spritesPresence
	lineIndex         int     // Index in video matrix / color line
	videoMatrix       []uint8 // Video matrix spritesPresence
	colorLine         []uint8 // Color line spritesPresence
	rowCounter        uint16  // Row counter
	videoCounter      uint16  // Video counter
	videoCounterLatch uint16  // Video counter base
	displayAccess     bool    // Display state
	textBuffer        []byte
}

// NewGraphics initializes and returns a new Graphics instance with the provided VIC core, collision handler, and display buffer.
func NewGraphics(core *VIC, collisions *Collisions, displayBuffer references.IDisplayBuffer) *Graphics {
	gr := &Graphics{
		core:              core,
		collisions:        collisions,
		set8:              displayBuffer.Set8,
		setMulti8:         displayBuffer.SetMulti8,
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

// Setup initializes the Graphics instance and prepares it for rendering operations.
// (Currently empty, but kept for consistency).
func (gr *Graphics) Setup() {
}

// GetText retrieves the text buffer content from the Graphics instance as a slice of bytes.
func (gr *Graphics) GetText() []byte {
	return gr.textBuffer
}

// PrintText outputs the contents of the textBuffer in a formatted manner, wrapping lines at every 40 characters.
// This is a helper function for debugging.
func (gr *Graphics) PrintText() {
	for x, v := range gr.textBuffer {
		if (x % 40) == 0 {
			fmt.Println()
		}
		fmt.Printf("%c", v)
	}
	fmt.Println()
}

// ResetVideoCounterLatch resets the video counter latch to zero.  This happens at the start of each frame.
func (gr *Graphics) ResetVideoCounterLatch() {
	gr.videoCounterLatch = 0
}

// UpdateVideoCounter updates the video counter to match the current video counter latch value (cycle 14).
func (gr *Graphics) UpdateVideoCounter() {
	gr.videoCounter = gr.videoCounterLatch
}

// ResetLineIndex resets the line index to zero. This happens at the beginning of each scanline (cycle 15).
func (gr *Graphics) ResetLineIndex() {
	gr.lineIndex = 0
}

// SetOffset sets the offset value for the Graphics instance to the given parameter.
// The offset represents the horizontal starting position (in pixels) for rendering.
func (gr *Graphics) SetOffset(offset int) {
	gr.offset = offset
}

// UpdateCharDataLast updates the `charDataLast` field with the current value of `charData`.
// This is needed for Extended Color Mode (ECM) where the background color is determined by the *previous* character.
func (gr *Graphics) UpdateCharDataLast() {
	gr.charDataLast = gr.charData
}

// TryResetRowCounter resets the row counter (RC) to 0 if the badLineCondition in the core is true (cycle 14).
func (gr *Graphics) TryResetRowCounter() {
	if gr.core.badLineCondition {
		gr.rowCounter = 0
	}
}

// TryAcquireDisplayAccess sets the displayAccess flag to true if the badLineCondition flag in the core is active.
// This gives the CPU access to video memory during "bad lines" (cycles 1-10 and 55-63).
func (gr *Graphics) TryAcquireDisplayAccess() {
	if gr.core.badLineCondition {
		gr.displayAccess = true
	}
}

// UpdateDisplayAccess updates the display access state and row counter based on the current conditions and row limit.
// This function manages the timing of when the display access is granted to the CPU and when the row counter is incremented.
// It's called at the *end* of each scanline (cycle 58).
func (gr *Graphics) UpdateDisplayAccess() {
	if gr.rowCounter == rowsMax {
		// If we've reached the end of a character row (8 pixel rows), latch the video counter.
		gr.videoCounterLatch = gr.videoCounter
		gr.displayAccess = false
	}
	if gr.core.badLineCondition || gr.displayAccess {
		// The & operator has precedence
		gr.rowCounter = (gr.rowCounter + 1) & rowsMax
		gr.displayAccess = true
	}
}

// TryGraphicsAccess fetches and processes graphics data from memory based on the current raster and display state.
// This function is the core of the VIC-II's graphics data fetching logic.  It's called in cycles 15-54, *if* displayAccess is true.
func (gr *Graphics) TryGraphicsAccess() {
	if gr.displayAccess {
		var addr uint16
		if gr.core.bmm {
			// Bitmap Mode: Calculate the address based on the video counter, bitmap base address, and row counter.
			addr = ((gr.videoCounter & 0x03ff) << 3) | gr.core.bitmapBase | gr.rowCounter // Bitmap
		} else {
			// Text Mode: Calculate the address based on the character code from the video matrix, character base address, and row counter.
			addr = (uint16(gr.videoMatrix[gr.lineIndex]) << 3) | gr.core.charBase | gr.rowCounter // Text
		}
		if gr.core.ecm {
			// Extended Color Mode (ECM): Mask the address to use only the lower 13 bits of the character ROM address.
			addr &= 0xf9ff
		}
		gr.gfxData = gr.core.ReadByte(addr)        // Read the graphics data (pixel data or character pattern).
		gr.charData = gr.videoMatrix[gr.lineIndex] // Get the character code from the video matrix.
		gr.colorData = gr.colorLine[gr.lineIndex]  // Get the color data from the color RAM.
		if gr.rowCounter == 0 {
			// At the beginning of a new character row (rowCounter == 0),
			// store the character code in the text buffer for debugging/display purposes.
			// https://sta.c64.org/cbm64scr.html
			index := columnsMax * (gr.core.rasterY / 8)
			gr.textBuffer[index+uint16(gr.lineIndex)] = _scCodesAscii[gr.charData&0x7f] // Convert screen code to ASCII
		}
		gr.lineIndex++    // Increment the line index to point to the next character/color data.
		gr.videoCounter++ // Increment the video counter.
	} else {
		// If display access is *not* granted, read from a "dummy" address.  The values read are not used.
		if gr.core.ecm {
			gr.gfxData = gr.core.ReadByte(0x39ff) // Dummy read (ECM).
		} else {
			gr.gfxData = gr.core.ReadByte(0x3fff) // Dummy read (non-ECM).
		}
		gr.colorData = 0 // Dummy values.
		gr.charData = 0  // Dummy values.
	}
}

// TryPhi2Access handles Phi2 clock phase access, updating videoMatrix and colorLine based on core conditions and memory.
// This function handles the memory access during the PHI2 phase of the CPU clock cycle. (cycles 15-54).
func (gr *Graphics) TryPhi2Access() {
	// Check if the Bus Available (BA) signal is low.
	if gr.core.baLow {
		// Check if the Address Enable Control (AEC) signal is low.
		if gr.core.aecLow {
			// If both BA and AEC are low, the VIC-II has access to the address bus.
			addr := (gr.videoCounter & 0x3ff) | gr.core.matrixBase  // Calculate address in video matrix.
			gr.videoMatrix[gr.lineIndex] = gr.core.ReadByte(addr)   // Read character code from video matrix.
			gr.colorLine[gr.lineIndex] = gr.core.readColorRam(addr) // Read color data from color RAM.
		} else {
			// If AEC is high, the CPU has access to the address bus, so we fill with dummy data.
			gr.colorLine[gr.lineIndex] = 0xff
			gr.videoMatrix[gr.lineIndex] = 0xff
		}
	}
}

// DrawBackground renders the background using the current display mode
// and updates the graphics offset and collision state (cycles 13-18 and 55-57).
func (gr *Graphics) DrawBackground() {
	// Call the appropriate background drawing function based on the current display mode.
	_backgroundSequencer[gr.core.displayMode](gr)
	// Increment the pixel offset by 8 (one character width).
	gr.offset += 8
	// Update the collision detection system's offset.
	gr.collisions.IncrementGraphicsOffset()
}

// DrawForeground renders the foreground graphics based on the current display mode and x-scroll offset.
// It also increments the offset and updates the collisions.
func (gr *Graphics) DrawForeground() {
	// Calculate the final offset, including x-scrolling.
	offset := gr.offset + int(gr.core.xScroll)
	// Call the appropriate foreground drawing function.
	_foregroundSequencer[gr.core.displayMode](gr, offset)
	// Increment the pixel offset by 8.
	gr.offset += 8
	// Update the collision detection system's offset.
	gr.collisions.IncrementGraphicsOffset()
}
