package mos6569

import "github.com/markel1974/c64emu/src/references"

const (
	dataCounterLastByte      = 0x3f
	spriteUnexpandedPixels   = 24
	spriteExpandedHalfPixels = 32
	spriteExpandedPixels     = 48
	planesMSB                = 1 << 31 //0x80000000
)

// Sprite represents a hardware sprite object used for rendering graphical elements on a screen.
// sNum is the sprite number, identifying the specific sprite.
// sBit is the sprite bit, used for masking and collision flag operations.
// core is a reference to the VIC, responsible for accessing hardware functionalities.
// data is a buffer storing the sprite's graphical data.
// dataPtr points to the memory address of the sprite data.
// counter is a counter used to track the fetching of sprite data.
// dataCounterBase serves as a base value for resetting the sprite data counter.
// set is a function for setting pixel data on the sprite's display buffer.
type Sprite struct {
	num         uint8
	max         int
	mask        uint8
	core        *VIC
	data        []uint8
	ptr         uint16
	counter     uint16
	counterBase uint16
	set         func(int, uint8)
}

// NewSprite initializes and returns a new Sprite instance with the provided VIC core, display buffer, and sprite number.
// It allocates memory for sprite data, sets initial values for counters, and configures the display function.
func NewSprite(core *VIC, displayBuffer references.IDisplayBuffer, sNum uint8, sMax int) *Sprite {
	sp := &Sprite{
		core:        core,
		set:         displayBuffer.Set,
		num:         sNum,
		mask:        uint8(1) << sNum,
		data:        make([]uint8, 4), // Allocate space for sprite data (3 byte + extra).
		counterBase: 0,
		counter:     dataCounterLastByte, // Initialize sprite data counter to the last byte.
		ptr:         0,
		max:         sMax,
	}
	return sp
}

// Number returns the unique identifier of the sprite as an 8-bit unsigned integer.
func (sp *Sprite) Number() uint8 {
	return sp.num
}

// Mask returns the bitmask representing the sprite's unique identifier for operations like collision detection.
func (sp *Sprite) Mask() uint8 {
	return sp.mask
}

// FetchPtr calculates and fetches the memory address of the sprite pointer from the VIC core and updates the sprite's data pointer.
func (sp *Sprite) FetchPtr() {
	addr := sp.core.matrixBase | 0x03f8 | uint16(sp.num) // Calculate the address of the sprite pointer.
	data := sp.core.ReadByte(addr)                       // Read the sprite pointer from memory.
	sp.ptr = uint16(data) << 6                           // Set the sprite's data pointer.
}

// FetchData retrieves the sprite data for the specified byte index and stores it in the sprite's data array.
// It calculates the memory address based on the sprite's data counter and pointer, reads the byte, and increments the counter.
func (sp *Sprite) FetchData(bNum uint8) {
	addr := (sp.counter & dataCounterLastByte) | sp.ptr // Calculate the address of the current byte within the sprite data.
	data := sp.core.ReadByte(addr)                      // Read the byte from memory.
	sp.data[bNum] = data                                // Store the byte in the sprite's data buffer.
	sp.counter++                                        // Increment the data counter for the next byte.
}

// CounterBase retrieves the base value of the sprite's data counter.
func (sp *Sprite) CounterBase() uint16 {
	return sp.counterBase
}

// CounterBaseReset resets the sprite's counterBase to zero. It is typically used to reinitialize the base counter.
func (sp *Sprite) CounterBaseReset() {
	sp.counterBase = 0
}

// CounterBaseIncrement increases the `counterBase` of the `Sprite` by the specified increment value.
func (sp *Sprite) CounterBaseIncrement(increment uint16) {
	sp.counterBase += increment
}

// CounterBaseApply sets the sprite's data counter to its base value stored in dataCounterBase.
func (sp *Sprite) CounterBaseApply() {
	sp.counter = sp.counterBase
}

// Draw renders the sprite on the screen at a specified starting scanline and manages collision detection.
func (sp *Sprite) Draw(lineStart int, collisions *Collisions) {
	sColor := _colors[sp.core.mXc[sp.num]]
	// Combine the three bytes of sprite data into a single 32-bit word for easier processing.
	sData := (uint32(sp.data[0]) << 24) | (uint32(sp.data[1]) << 16) | (uint32(sp.data[2]) << 8)
	// Calculate the sprite's X offset on the screen. Add 24 to account for the border.
	sOffset := int(sp.core.mXx[sp.num]) + sp.max //SpriteNumber
	// Calculate the final offset on the scanline, including the global offset.
	lineOffset := lineStart + sOffset // lineStart + sOffset
	// Calculate the "major" X coordinate (used for collision detection).This is essentially the character column.
	majorX := sOffset / sp.max
	// Calculate the "minor" X coordinate (used for collision detection).This is the pixel offset within the character column.
	minorX := sOffset & 7
	// Check if the sprite is in multicolor mode.
	multiColor := (sp.core.mmc & sp.mask) != 0
	if expanded := (sp.core.mxe & sp.mask) != 0; expanded {
		// If the sprite is expanded horizontally...
		if multiColor {
			sp.drawExpandedMulticolor(lineOffset, collisions, sColor, sData, sOffset, majorX, minorX)
		} else {
			sp.drawExpandedStandard(lineOffset, collisions, sColor, sData, sOffset, majorX, minorX)
		}
	} else {
		// If the sprite is *not* expanded horizontally...
		if multiColor {
			sp.drawUnexpandedMulticolor(lineOffset, collisions, sColor, sData, sOffset, majorX, minorX)
		} else {
			sp.drawUnexpandedStandard(lineOffset, collisions, sColor, sData, sOffset, majorX, minorX)
		}
	}
}

// drawExpandedMulticolor is responsible for rendering an expanded multicolor sprite on the display buffer.
// It handles sprite-to-graphics and sprite-to-sprite collision detection while processing each pixel's color.
// The method also respects sprite-to-background priority by masking sprite pixels when a graphics collision occurs.
func (sp *Sprite) drawExpandedMulticolor(lineOffset int, collisions *Collisions, sColor uint8, sData uint32, sOffset int, m int, s int) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := collisions.GetGraphicsL(m, s)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := collisions.GetGraphicsR(m, s)
	// Expand sprite data horizontally using lookup tables.The _multiExpTable expands each byte (8 bits)
	// into two bytes (16 bits), doubling the width.This is done for multicolor mode *before* bit-plane conversion.
	sDataL := (uint32(_multiExpTable[(sData>>24)&0xff]) << 16) | (uint32(_multiExpTable[(sData>>16)&0xff]))
	sDataR := uint32(_multiExpTable[(sData>>8)&0xff]) << 16
	// Convert sprite data to bit-planes for easier multicolor processing.In multicolor mode, each pixel is represented
	// by *two* bits (hence two bit-planes).
	plane0L := (sDataL & 0x55555555) | ((sDataL & 0x55555555) << 1) // Bit-plane 0 (left half).
	plane1L := (sDataL & 0xaaaaaaaa) | ((sDataL & 0xaaaaaaaa) >> 1) // Bit-plane 1 (left half).
	plane0R := (sDataR & 0x55555555) | ((sDataR & 0x55555555) << 1) // Bit-plane 0 (right half).
	plane1R := (sDataR & 0xaaaaaaaa) | ((sDataR & 0xaaaaaaaa) >> 1) // Bit-plane 1 (right half).
	// Collision with graphics? Check if any bits in the sprite's bit-planes overlap with the foreground mask.
	if ((foreMaskL & (plane0L | plane1L)) != 0) || ((foreMaskR & (plane0R | plane1R)) != 0) {
		// Set the sprite-to-graphics collision flag for this sprite.
		collisions.SetGraphicsCollision(sp.mask)
		if (sp.core.mdp & sp.mask) != 0 {
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
		if selectedColor := sp.planes2Color(plane0L, plane1L, sColor); selectedColor >= 0 {
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, uint8(selectedColor))
			}
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for idx := spriteExpandedHalfPixels; idx < spriteExpandedPixels; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		if selectedColor := sp.planes2Color(plane0R, plane1R, sColor); selectedColor >= 0 {
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, uint8(selectedColor))
			}
		}
	}
}

// drawExpandedStandard renders a horizontally expanded sprite in standard mode with collision detection and masking.
func (sp *Sprite) drawExpandedStandard(lineOffset int, collisions *Collisions, sColor uint8, sData uint32, sOffset int, m int, s int) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := collisions.GetGraphicsL(m, s)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := collisions.GetGraphicsR(m, s)
	// Expand sprite data horizontally using lookup tables.The _expTable expands each byte (8 bits) into two bytes (16 bits).
	sDataL := uint32(_expTable[(sData>>24)&0xff])<<16 | uint32(_expTable[(sData>>16)&0xff])
	sDataR := uint32(_expTable[(sData>>8)&0xff]) << 16
	// Check for collisions with the foreground.
	if ((foreMaskL & sDataL) != 0) || ((foreMaskR & sDataR) != 0) {
		// Set the sprite-to-graphics collision flag.
		collisions.SetGraphicsCollision(sp.mask)
		if (sp.core.mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sDataL &= ^foreMaskL // Mask left half.
			sDataR &= ^foreMaskR // Mask right half.
		}
	}

	// Draw the left half of the sprite (first 32 pixels).
	for idx := 0; idx < spriteExpandedHalfPixels; idx, sDataL = idx+1, sDataL<<1 {
		if (sDataL & planesMSB) != 0 {
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, sColor)
			}
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for idx := spriteExpandedHalfPixels; idx < spriteExpandedPixels; idx, sDataR = idx+1, sDataR<<1 {
		if (sDataR & planesMSB) != 0 {
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, sColor)
			}
		}
	}
}

// drawUnexpandedMulticolor renders a non-expanded multicolor sprite, handling graphics collisions and color selection.
func (sp *Sprite) drawUnexpandedMulticolor(lineOffset int, collisions *Collisions, sColor uint8, sData uint32, sOffset int, m int, s int) {
	// Get the foreground mask for the sprite's character column.
	foreMask := collisions.GetGraphicsL(m, s)
	// Convert sprite data to bit-planes.  No expansion is needed here since the sprite is *not* expanded.
	p0 := sData & 0x55555555 // sprite to bit-planes 0.
	p1 := sData & 0xaaaaaaaa // sprite to bit-Planes 1.
	plane0 := p0 | (p0 << 1) // Combine bits for plane 0.
	plane1 := p1 | (p1 >> 1) // Combine bits for plane 1.
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		// Set the sprite-to-graphics collision flag.
		collisions.SetGraphicsCollision(sp.mask)
		if (sp.core.mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			plane0 &= ^foreMask // Mask plane 0.
			plane1 &= ^foreMask // Mask plane 1.
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < spriteUnexpandedPixels; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		if selectedColor := sp.planes2Color(plane0, plane1, sColor); selectedColor >= 0 {
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, uint8(selectedColor))
			}
		}
	}
}

// drawUnexpandedStandard renders a 24-pixel wide unexpanded standard sprite, managing collisions and pixel masking.
func (sp *Sprite) drawUnexpandedStandard(lineOffset int, collisions *Collisions, sColor uint8, sData uint32, sOffset int, m int, s int) {
	// Get the foreground mask for the sprite's character column.
	foreMask := collisions.GetGraphicsL(m, s)
	// Check for collisions with the foreground.
	if (foreMask & sData) != 0 {
		// Set the sprite-to-graphics collision flag.
		collisions.SetGraphicsCollision(sp.mask)
		if (sp.core.mdp & sp.mask) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sData &= ^foreMask // Mask the sprite data.
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < spriteUnexpandedPixels; idx, sData = idx+1, sData<<1 {
		if (sData & planesMSB) != 0 {
			// Check for sprite-to-sprite collisions *before* drawing.
			if !collisions.SetSpriteCollision(sOffset+idx, sp.mask) {
				sp.set(lineOffset+idx, sColor)
			}
		}
	}
}

// planesColor evaluates the color of a sprite pixel based on the provided plane data and sprite color, returning a color value.
func (sp *Sprite) planes2Color(plane0 uint32, plane1 uint32, sColor uint8) int {
	p1 := (plane1 & planesMSB) >> 0x1e //bit 30
	p0 := (plane0 & planesMSB) >> 0x1f //bit 31
	switch p1 | p0 {
	case 0b00:
		return -1 // transparent
	case 0b01:
		return int(_colors[sp.core.mm0]) //mm0
	case 0b10:
		return int(sColor) // color
	case 0b11:
		return int(_colors[sp.core.mm1]) //mm1
	}
	return -1
}
