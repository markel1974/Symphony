package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	dataCounterLastByte      = 0x3f
	spriteUnexpandedPixels   = 24
	spriteExpandedHalfPixels = 32
	spriteExpandedPixels     = 48
	planesMSB                = 1 << 31 // planesMSB is the MSB Used to check the leftmost pixel containing sprite data.
)

const (
	dataAlignment = 3
)

const (
	LatchCount  = 2
	LatchOffset = 8
	LatchIndex0 = 0
	LatchIndex1 = 1 * LatchOffset
)

const (
	latchData0 = iota
	latchData1
	latchData2
	latchMDP
	latchMM0
	latchMM1
	latchMXC
	latchOffset
)

// Sprite represents a hardware sprite object used for rendering graphical elements on a screen.
// sNum is the sprite number, identifying the specific sprite.
// sBit is the sprite bit, used for masking and collision flag operations
// core is a reference to the VIC, responsible for accessing hardware functionalities
// data is a buffer storing the sprite's graphical data
// dataPtr points to the memory address of the sprite data
// counter is a counter used to track the fetching of sprite data
// dataCounterBase serves as a base value for resetting the sprite data counter
// set is a function for setting pixel data on the sprite's display buffer
type Sprite struct {
	*component.BaseComponent
	collisions       *CollisionsUnit
	beam             *RasterBeam
	num              uint8
	max              int
	mask             uint8
	latch            []int
	ptr              uint16
	counter          uint16
	counterBase      uint16
	counterIncrement uint16
	latchIndex       int
	draw             func(mdp uint8, mm0 uint8, mm1 uint8, mxc uint8, sOffset int, majorX int, minorX int, data0 uint8, data1 uint8, data2 uint8)
	//plane2ColorHelper [4]func(uint8) int
}

// NewSprite initializes and returns a new Sprite instance with the provided VIC core, display buffer, and sprite number.
// It allocates memory for sprite data, sets initial values for counters, and configures the display function.
func NewSprite(parent references.IComponent, factory references.IComponentFactory, label string, instance int, beam *RasterBeam, collisions *CollisionsUnit, sNum uint8, sMax int) *Sprite {
	sp := &Sprite{
		BaseComponent:    component.NewBaseComponent(),
		collisions:       collisions,
		beam:             beam,
		num:              sNum,
		mask:             uint8(1) << sNum,
		latch:            make([]int, (LatchCount*LatchOffset)+1),
		counterBase:      0,
		counterIncrement: 0,
		counter:          dataCounterLastByte, // Initialize sprite data counter to the last byte.
		ptr:              0,
		max:              sMax,
		latchIndex:       LatchIndex0,
	}
	//sp.plane2ColorHelper = sp.createPlane2ColorHelper()
	sp.draw = sp.drawUnexpandedStandard
	sp.BaseComponent.Register(factory, parent, "sprite", sp, references.IdInternalComponent(label, instance, "Sprite"))
	return sp
}

// Setup initializes the Sprite instance with necessary configurations prior to rendering or operation.
func (sp *Sprite) Setup() error {
	return nil
}

// Connect establishes a connection between the sprite and its core, enabling necessary hardware dependencies.
// Returns an error if the connection fails.
func (sp *Sprite) Connect() error {
	return nil
}

// EmulationRequired determines if emulation logic is required for the sprite, always returning false.
func (sp *Sprite) EmulationRequired() bool {
	return false
}

// Emulate executes the main processing logic for the sprite during each emulation cycle, updating its state and behavior.
func (sp *Sprite) Emulate() {
}

// Internal checks if the sprite operates in internal mode, typically for system-bound sprites or debug states.
func (sp *Sprite) Internal() bool {
	return true
}

// Reset reinitializes the sprite's internal state, counters, and data to their default configuration.
func (sp *Sprite) Reset() {
}

func (sp *Sprite) SetLatchIndex(index int) {
	sp.latchIndex = index
}

// Number returns the unique identifier of the sprite as an 8-bit unsigned integer.
func (sp *Sprite) Number() uint8 {
	return sp.num
}

func (sp *Sprite) Counter() uint16 {
	return sp.counter
}

func (sp *Sprite) SetPtr(ptr uint16) {
	sp.ptr = ptr
}

func (sp *Sprite) Ptr() uint16 {
	return sp.ptr
}

// Mask returns the bitmask representing the sprite's unique identifier for operations like collision detection.
func (sp *Sprite) Mask() uint8 {
	return sp.mask
}

func (sp *Sprite) SetData(bNum int, data uint8) {
	//latch
	sp.latch[sp.latchIndex+bNum] = int(data)
	sp.counter++
}

// LatchAttributes updates the sprite's color data, multicolor flags, and offset using the provided parameters.
func (sp *Sprite) LatchAttributes(mdp uint8, mm0 uint8, mm1 uint8, mXc uint8, mXx uint16) {
	sp.latch[sp.latchIndex+latchMDP] = int(mdp)
	sp.latch[sp.latchIndex+latchMM0] = int(mm0)
	sp.latch[sp.latchIndex+latchMM1] = int(mm1)
	sp.latch[sp.latchIndex+latchMXC] = int(mXc)
	sp.latch[sp.latchIndex+latchOffset] = int(mXx) + spriteNumber
}

// IncrementCounterBase updates the counterIncrement field of the sprite to align with the current data alignment value.
func (sp *Sprite) IncrementCounterBase() {
	sp.counterIncrement = dataAlignment
}

// CommitIncrementCounterBase updates the counter-base by adding the counter-increment, resets the increment, and checks a condition.
func (sp *Sprite) CommitIncrementCounterBase() bool {
	sp.counterBase += sp.counterIncrement
	sp.counterIncrement = 0
	return (sp.counterBase & 0x3f) == 0x3f
}

// ResetCounterBase resets the sprite's counterBase to zero. It is typically used to reinitialize the base counter.
func (sp *Sprite) ResetCounterBase() {
	sp.counterBase = 0
}

// CommitCounterBase sets the sprite's data counter to its base value stored in dataCounterBase.
func (sp *Sprite) CommitCounterBase() {
	sp.counter = sp.counterBase
}

// ModeUpdate updates the sprite's rendering mode pipeline based on multicolor and horizontal expansion flags.
func (sp *Sprite) ModeUpdate(mmc uint8, mxe uint8) {
	multiColor := (mmc & sp.mask) != 0
	expandedH := (mxe & sp.mask) != 0
	if expandedH {
		if multiColor {
			sp.draw = sp.drawExpandedMulticolor
		} else {
			sp.draw = sp.drawExpandedStandard
		}
	} else {
		if multiColor {
			sp.draw = sp.drawUnexpandedMulticolor
		} else {
			sp.draw = sp.drawUnexpandedStandard
		}
	}
}

func (sp *Sprite) Draw(latchIndex int) {
	sOffset := sp.latch[latchIndex+latchOffset]
	data0 := uint8(sp.latch[latchIndex+latchData0])
	data1 := uint8(sp.latch[latchIndex+latchData1])
	data2 := uint8(sp.latch[latchIndex+latchData2])
	mdp := uint8(sp.latch[latchIndex+latchMDP])
	mm0 := uint8(sp.latch[latchIndex+latchMM0])
	mm1 := uint8(sp.latch[latchIndex+latchMM1])
	mxc := uint8(sp.latch[latchIndex+latchMXC])
	charColumn := sOffset / spriteNumber // This is the character column (column char, 0-39)
	pixelOffset := sOffset & spriteIndex // This is the pixel offset within the character column (offset pixel inside column, 0-7).
	sp.draw(mdp, mm0, mm1, mxc, sOffset, charColumn, pixelOffset, data0, data1, data2)
}

// drawExpandedMulticolor is responsible for rendering an expanded multicolor sprite on the display buffer.
// It handles sprite-to-graphics and sprite-to-sprite collision detection while processing each pixel's color.
// The method also respects sprite-to-background priority by masking sprite pixels when a graphics collision occurs.
func (sp *Sprite) drawExpandedMulticolor(mdp uint8, mm0 uint8, mm1 uint8, mxc uint8, sOffset int, charColumn int, pixelOffset int, data0 uint8, data1 uint8, data2 uint8) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := sp.collisions.GetGraphicsL(charColumn, pixelOffset)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := sp.collisions.GetGraphicsR(charColumn, pixelOffset)
	// Expand sprite data horizontally using lookup tables.The _multiExpTable expands each byte (8 bits)
	// into two bytes (16 bits), doubling the width.This is done for multicolor mode *before* bit-plane conversion.
	sDataL := uint32(_multiExpTable[data0])<<16 | uint32(_multiExpTable[data1])
	sDataR := uint32(_multiExpTable[data2]) << 16
	// Convert sprite data to bit-planes for easier multicolor processing.In multicolor mode, each pixel is represented
	// by *two* bits (hence two bit-planes).
	plane0L := (sDataL & 0x55555555) | ((sDataL & 0x55555555) << 1) // Bit-plane 0 (left half).
	plane1L := (sDataL & 0xaaaaaaaa) | ((sDataL & 0xaaaaaaaa) >> 1) // Bit-plane 1 (left half).
	plane0R := (sDataR & 0x55555555) | ((sDataR & 0x55555555) << 1) // Bit-plane 0 (right half).
	plane1R := (sDataR & 0xaaaaaaaa) | ((sDataR & 0xaaaaaaaa) >> 1) // Bit-plane 1 (right half).
	// Collision with graphics? Check if any bits in the sprite's bit-planes overlap with the foreground mask.
	// Combine both bit-planes with logical OR
	// test collision with a single bitwise AND operation
	// handles all 48 pixels of the expanded sprite simultaneously
	if ((foreMaskL & (plane0L | plane1L)) != 0) || ((foreMaskR & (plane0R | plane1R)) != 0) {
		// Set the sprite-to-graphics collision flag for this sprite.
		sp.collisions.SetGraphicsPresence(sp.mask)
		if (mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled (MDP register), mask out the sprite pixels where a collision occurred.
			// This makes the background "show through" the sprite.
			plane0L &= ^foreMaskL // Mask left half.
			plane1L &= ^foreMaskL // Mask left half.
			plane0R &= ^foreMaskR // Mask right half.
			plane1R &= ^foreMaskR // Mask right half.
		}
	}

	// Draw the left half of the sprite (first 32 pixels). The sprite is expanded, so we draw 48 pixels total.
	for idx := 0; idx < spriteExpandedHalfPixels; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		if selectedColor := sp.planes2Color(plane0L, plane1L, mxc, mm0, mm1); selectedColor >= 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, uint8(selectedColor))
			}
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for idx := spriteExpandedHalfPixels; idx < spriteExpandedPixels; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		if selectedColor := sp.planes2Color(plane0R, plane1R, mxc, mm0, mm1); selectedColor >= 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, uint8(selectedColor))
			}
		}
	}
}

// drawExpandedStandard renders a horizontally expanded sprite in standard mode with collision detection and masking.
func (sp *Sprite) drawExpandedStandard(mdp uint8 /* mm0 */, _ uint8 /* mm1 */, _ uint8, mxc uint8, sOffset int, charColumn int, pixelOffset int, data0 uint8, data1 uint8, data2 uint8) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := sp.collisions.GetGraphicsL(charColumn, pixelOffset)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := sp.collisions.GetGraphicsR(charColumn, pixelOffset)
	// Expand sprite data horizontally using lookup tables.The _expTable expands each byte (8 bits) into two bytes (16 bits).
	sDataL := uint32(_expTable[data0])<<16 | uint32(_expTable[data1])
	sDataR := uint32(_expTable[data2]) << 16
	// Check for collisions with the foreground.
	if ((foreMaskL & sDataL) != 0) || ((foreMaskR & sDataR) != 0) {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsPresence(sp.mask)
		if (mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sDataL &= ^foreMaskL // Mask left half.
			sDataR &= ^foreMaskR // Mask right half.
		}
	}

	// Draw the left half of the sprite (first 32 pixels).
	for idx := 0; idx < spriteExpandedHalfPixels; idx, sDataL = idx+1, sDataL<<1 {
		if (sDataL & planesMSB) != 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, mxc)
			}
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for idx := spriteExpandedHalfPixels; idx < spriteExpandedPixels; idx, sDataR = idx+1, sDataR<<1 {
		if (sDataR & planesMSB) != 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, mxc)
			}
		}
	}
}

// drawUnexpandedMulticolor renders a non-expanded multicolor sprite, handling graphics collisions and color selection.
func (sp *Sprite) drawUnexpandedMulticolor(mdp uint8, mm0 uint8, mm1 uint8, mxc uint8, sOffset int, charColumn int, pixelOffset int, data0 uint8, data1 uint8, data2 uint8) {
	// Get the foreground mask for the sprite's character column.
	foreMask := sp.collisions.GetGraphicsL(charColumn, pixelOffset)
	// Combine the three bytes of sprite data into a single 32-bit word for easier processing.
	sData := (uint32(data0) << 24) | (uint32(data1) << 16) | (uint32(data2) << 8)
	// Convert sprite data to bit-planes.  No expansion is needed here since the sprite is *not* expanded.
	p0 := sData & 0x55555555 //bit-plane 0: Extracts odd bits (0x55555555 = 01010101...)
	p1 := sData & 0xaaaaaaaa //bit-plane 1: Extracts even bits (0xaaaaaaaa = 10101010...)
	//bit duplication (<< 1 and >> 1) is needed for multicolor mode where each logical pixel uses 2 physical bits
	plane0 := p0 | (p0 << 1) // Combine bits for plane 0.
	plane1 := p1 | (p1 >> 1) // Combine bits for plane 1.
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsPresence(sp.mask)
		if (mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			plane0 &= ^foreMask
			plane1 &= ^foreMask
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < spriteUnexpandedPixels; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		if selectedColor := sp.planes2Color(plane0, plane1, mxc, mm0, mm1); selectedColor >= 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, uint8(selectedColor))
			}
		}
	}
}

// drawUnexpandedStandard renders a 24-pixel wide unexpanded standard sprite, managing collisions and pixel masking.
func (sp *Sprite) drawUnexpandedStandard(mdp uint8 /* mm0 */, _ uint8 /* mm1 */, _ uint8, mxc uint8, sOffset int, charColumn int, pixelOffset int, data0 uint8, data1 uint8, data2 uint8) {
	// Get the foreground mask for the sprite's character column.
	foreMask := sp.collisions.GetGraphicsL(charColumn, pixelOffset)
	// Combine the three bytes of sprite data into a single 32-bit word for easier processing.
	sData := (uint32(data0) << 24) | (uint32(data1) << 16) | (uint32(data2) << 8)
	// Check for collisions with the foreground.
	if (foreMask & sData) != 0 {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsPresence(sp.mask)
		if (mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sData &= ^foreMask
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < spriteUnexpandedPixels; idx, sData = idx+1, sData<<1 {
		if (sData & planesMSB) != 0 {
			if !sp.collisions.SetSpritePresence(sOffset+idx, sp.mask) {
				sp.beam.Draw(sOffset+idx, mxc)
			}
		}
	}
}

// planesColor evaluates the color of a sprite pixel based on the provided plane data and sprite color, returning a color value.
func (sp *Sprite) planes2Color(plane0 uint32, plane1 uint32, sColor uint8, mm0 uint8, mm1 uint8) int {
	bit1 := (plane1 >> 31) & 1
	bit0 := (plane0 >> 31) & 1
	index := (bit1 << 1) | bit0
	//return sp.plane2ColorHelper[index&4](mxc)
	switch index {
	case 0b00:
		return -1 // transparent
	case 0b01:
		return int(mm0) //mm0
	case 0b10:
		return int(sColor) // color
	case 0b11:
		return int(mm1) //mm1
	}
	return -1
}

/*
func (sp *Sprite) createPlane2ColorHelper() [4]func(uint8) int {
	var p2c [4]func(uint8) int
	p2c[0] = func(uint8) int {
		return -1 // transparent
	}
	p2c[1] = func(uint8) int {
		return int(_colors[sp.core.mm0]) //mm0
	}
	p2c[2] = func(mxc uint8) int {
		return int(mxc) // color
	}
	p2c[3] = func(uint8) int {
		return int(_colors[sp.core.mm1]) //mm1
	}
	return p2c
}
*/
