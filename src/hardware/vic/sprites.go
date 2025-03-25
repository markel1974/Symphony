package mos6569

import (
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// Sprites represent the structure responsible for handling sprites, including their data, states, and configurations.
// It contains properties for managing sprite visual data, collision detection, and display control buffers.
// The type relies on various counters, flags, and pointers to handle sprite DMA and display activities effectively.
// It interacts with the VIC core and an implemented display buffer interface for rendering and collision processing.
type Sprites struct {
	core            *VIC                      // Pointer to the main VIC-II core.
	displayBuffer   references.IDisplayBuffer // Interface to the display buffer.
	collisions      *Collisions               // Pointer to the collision detection system.
	dataPtr         []uint16                  // Sprite data pointers (one per sprite).
	data            [][]uint8                 // Sprite data (up to 64 bytes per sprite).
	dmaFlags        uint8                     // Active DMA Sprite (bitmask: bit i = 1 means sprite is active).
	displayFlags    uint8                     // Active Display Sprite (bitmask).
	spriteFlags     uint8                     // Sprite in this line (bitmask).
	dataCounter     []uint16                  // Sprite counter data (one per sprite).
	dataCounterBase []uint16                  // Sprite base counter data (one per sprite).
	offset          int                       // Offset from bitmap spritesBuffer
}

// NewSprites initializes and returns a new instance of the Sprites struct with default settings and allocations.
// It sets up sprite data, counters, and dependencies using the provided VIC core, collisions, and display buffer.
func NewSprites(core *VIC, collisions *Collisions, displayBuffer references.IDisplayBuffer) *Sprites {
	s := &Sprites{
		core:            core,
		displayBuffer:   displayBuffer,
		collisions:      collisions,
		dataPtr:         make([]uint16, SpriteNumber),
		data:            make([][]uint8, SpriteNumber),
		displayFlags:    0,
		dmaFlags:        0,
		dataCounter:     make([]uint16, SpriteNumber),
		dataCounterBase: make([]uint16, SpriteNumber),
		offset:          0,
	}
	for i := range s.data {
		// Allocate space for sprite data (3 byte + extra).
		s.data[i] = make([]uint8, 4)
	}
	for i := range s.dataCounter {
		// Initialize sprite data counter to the last byte.
		s.dataCounter[i] = 63
	}
	return s
}

// Setup initializes the Sprites instance, preparing internal state and configurations needed for sprite operations.
func (sp *Sprites) Setup() {
	// Nothing to do here at the moment, as initialization is handled in NewSprites.
	// This function is kept for consistency and potential future use.
}

// SetOffset updates the `offset` value of the Sprites instance with the given value.
// This offset is used to calculate the starting position for sprite rendering on the current scanline.
func (sp *Sprites) SetOffset(offset int) {
	sp.offset = offset
}

// FetchPtr fetches the sprite pointer for the given sprite number if BA and AEC conditions are met,
// and updates its data pointer. Logs a warning if conditions are not met.
// This function is called during specific VIC-II cycles when sprite data pointers need to be fetched
// from memory.  The actual memory address is calculated based on the sprite pointer value
// and the VIC-II's memory mapping.
func (sp *Sprites) FetchPtr(num int) {
	if sp.core.baLow && sp.core.aecLow {
		addr := sp.core.matrixBase | 0x03f8 | uint16(num) // Calculate the address of the sprite pointer.
		data := sp.core.ReadByte(addr)                    // Read the sprite pointer from memory.
		ptr := uint16(data) << 6                          // Calculate the base address of the sprite data.(each sprite ptr points to a 64-byte block).
		sp.dataPtr[num] = ptr                             // Store the calculated pointer.
	} else {
		log.Printf("sprites: can't fetch sprite ptr %d", num) // Should not normally happen, as the VIC-II controls BA/AEC.
	}
}

// Fetch loads sprite data for the given sprite number and byte index from memory if BA and AEC lines are low.
// It updates the sprite's data array and increments its data counter.Logs an error if BA or AEC lines are high.
// This function is called during specific VIC-II cycles (typically cycles 49-55, three times per sprite) to fetch
// the actual sprite data (3 bytes per sprite per cycle).
func (sp *Sprites) Fetch(num int, bNum int) {
	if sp.core.baLow && sp.core.aecLow {
		ptr := sp.dataPtr[num]                     // Get the base address of the sprite data.
		addr := (sp.dataCounter[num] & 0x3f) | ptr // Calculate the address of the current byte within the sprite data.
		data := sp.core.ReadByte(addr)             // Read the byte from memory.
		sp.data[num][bNum] = data                  // Store the byte in the sprite's data buffer.
		sp.dataCounter[num]++                      // Increment the data counter for the next byte.
	} else {
		log.Printf("sprites: can't fetch sprite %d - %d", num, bNum) // Should not normally happen.
	}
}

// UpdateDisplayFlags updates the display flags for sprites by checking and clearing flags based on DMA activity status.
// It determines which sprites are currently active for display based on DMA and counter-status.
// Called at the *end* of each scanline (cycle 58).
func (sp *Sprites) UpdateDisplayFlags() {
	sp.spriteFlags = sp.displayFlags
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if ((sp.displayFlags & mask) != 0) && ((sp.dmaFlags & mask) == 0) {
			sp.displayFlags &= ^mask
		}
	}
}

// UpdateCounterBase increments the base counters of sprites with vertical expansion enabled by 2 for each enabled sprite.
// Handles y-expansion
func (sp *Sprites) UpdateCounterBase() {
	for idx := 0; idx < SpriteNumber; idx++ {
		if (sp.core.sprExpY & (1 << idx)) != 0 {
			sp.dataCounterBase[idx] += 2
		}
	}
}

// GetDMAFlag checks and returns the active DMA flag for the specified sprite(s) by performing a bitwise AND operation.
func (sp *Sprites) GetDMAFlag(b uint8) uint8 {
	return sp.dmaFlags & b
}

// UpdateDMA updates the DMA status of sprites based on their raster line and enabled flags.
// Called in cycle 12.
func (sp *Sprites) UpdateDMA() {
	rasterY := sp.core.rasterY & 0xff
	for i, mask := 0, uint8(1); i < SpriteNumber; i, mask = i+1, mask<<1 {
		if ((sp.core.me & mask) != 0) && (rasterY == uint16(sp.core.mXy[i])) {
			sp.dmaFlags |= mask
			sp.dataCounterBase[i] = 0
			if (sp.core.mye & mask) != 0 {
				sp.core.sprExpY &= ^mask
			}
		}
	}
}

// UpdateCounterBaseDMA updates the base counters of sprites and manages DMA flags based on specific conditions.
// Called in cycle 13.
func (sp *Sprites) UpdateCounterBaseDMA() {
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if (sp.core.sprExpY & mask) != 0 {
			sp.dataCounterBase[idx]++
		}
		if (sp.dataCounterBase[idx] & 0x3f) == 0x3f {
			sp.dmaFlags &= ^mask
		}
	}
}

// UpdateDisplayYFlags updates the display flags for sprites based on raster line position and active DMA flags.
// Called at the beginning of each scanline (cycle 14).
func (sp *Sprites) UpdateDisplayYFlags() {
	rasterY := sp.core.rasterY & 0xff
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		sp.dataCounter[idx] = sp.dataCounterBase[idx]
		if ((sp.dmaFlags & mask) != 0) && (rasterY == uint16(sp.core.mXy[idx])) {
			sp.displayFlags |= mask
		}
	}
}

// Draw renders all active sprites for the current line based on their flags, properties, and configurations.
// It handles both expanded and unexpanded sprites in standard and multicolor modes.
// Collision detection for sprites is carried out during the rendering process.
// Called in cycles 57-62.
func (sp *Sprites) Draw() {
	if sp.spriteFlags == 0 {
		// If no sprites are active on this scanline, return early.
		return
	}

	// Prepare the collision detection system for this scanline.
	sp.collisions.Prepare()

	for sNum, sBit := uint8(0), uint8(1); sNum < SpriteNumber; sNum, sBit = sNum+1, sBit<<1 {
		if (sp.spriteFlags & sBit) != 0 {
			// If this sprite is active on this scanline...
			// Get the sprite's color from the VIC-II registers.
			sColor := _colors[sp.core.mXc[sNum]]
			// Combine the three bytes of sprite data into a single 32-bit word for easier processing.
			sData := (uint32(sp.data[sNum][0]) << 24) | (uint32(sp.data[sNum][1]) << 16) | (uint32(sp.data[sNum][2]) << 8)
			// Calculate the sprite's X offset on the screen. Add 24 to account for the border.
			sOffset := int(sp.core.mXx[sNum]) + SpriteNumber
			// Calculate the final offset on the scanline, including the global offset.
			lineOffset := sp.offset + sOffset // lineStart + sOffset
			// Calculate the "major" X coordinate (used for collision detection).This is essentially the character column.
			m := sOffset / SpriteNumber
			// Calculate the "minor" X coordinate (used for collision detection).This is the pixel offset within the character column.
			s := sOffset & 7
			// Check if the sprite is in multicolor mode.
			multiColor := (sp.core.mmc & sBit) != 0
			if expanded := (sp.core.mxe & sBit) != 0; expanded {
				// If the sprite is expanded horizontally...
				if multiColor {
					// ...and in multicolor mode.
					sp.drawExpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit)
				} else {
					// ...and in standard color mode.
					sp.drawExpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit)
				}
			} else {
				// If the sprite is *not* expanded horizontally...
				if multiColor {
					// ...and in multicolor mode.
					sp.drawUnexpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit)
				} else {
					// ...and in standard color mode.
					sp.drawUnexpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit)
				}
			}
		}
	}
	// Perform the final collision detection checks.
	sp.collisions.Detect()
}

// drawExpandedMulticolor renders expanded multicolor sprites with applied graphics and collision handling logic.
func (sp *Sprites) drawExpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := sp.collisions.GetGraphicsL(m, s)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := sp.collisions.GetGraphicsR(m, s)
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
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			// If sprite-to-background priority is enabled (MDP register), mask out the sprite pixels where a collision occurred.
			// This makes the background "show through" the sprite.
			plane0L &= ^foreMaskL // Mask left half.
			plane1L &= ^foreMaskL // Mask left half.
			plane0R &= ^foreMaskR // Mask right half.
			plane1R &= ^foreMaskR // Mask right half.
		}
	}
	idx := 0
	// Draw the left half of the sprite (first 32 pixels).  The sprite is expanded, so we draw 48 pixels total.
	for ; idx < 32; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		selectedColor := uint8(0)
		// Determine the color of the current pixel based on the bit-planes.
		if (plane1L & 0x80000000) != 0 {
			// Check the most significant bit of plane 1.
			if (plane0L & 0x80000000) != 0 {
				// Check the most significant bit of plane 0.
				selectedColor = _colors[sp.core.mm1] // 11: Use color from MM1 register.
			} else {
				selectedColor = sColor // 10: Use the sprite's individual color.
			}
		} else {
			if (plane0L & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm0] // 01: Use color from MM0 register.
			} else {
				continue // 00: Transparent - don't draw anything.
			}
		}
		// Check for sprite-to-sprite collisions *before* drawing the pixel.
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			// Draw the pixel to the display buffer.
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for ; idx < 48; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		selectedColor := uint8(0)
		// Determine the color of the current pixel (same logic as above).
		if (plane1R & 0x80000000) != 0 {
			if (plane0R & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm1]
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0R & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm0]
			} else {
				continue // Transparent.
			}
		}
		// Check for sprite-to-sprite collisions.
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
}

// drawExpandedStandard draws an expanded standard sprite on the display buffer with collision detection and masking checks.
func (sp *Sprites) drawExpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	// Get the foreground mask for the left half of the sprite's character column.
	foreMaskL := sp.collisions.GetGraphicsL(m, s)
	// Get the foreground mask for the right half of the sprite's character column.
	foreMaskR := sp.collisions.GetGraphicsR(m, s)
	// Expand sprite data horizontally using lookup tables.The _expTable expands each byte (8 bits) into two bytes (16 bits).
	sDataL := uint32(_expTable[(sData>>24)&0xff])<<16 | uint32(_expTable[(sData>>16)&0xff])
	sDataR := uint32(_expTable[(sData>>8)&0xff]) << 16
	// Check for collisions with the foreground.
	if ((foreMaskL & sDataL) != 0) || ((foreMaskR & sDataR) != 0) {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sDataL &= ^foreMaskL // Mask left half.
			sDataR &= ^foreMaskR // Mask right half.
		}
	}
	var idx = 0
	// Draw the left half of the sprite (first 32 pixels).
	for ; idx < 32; idx, sDataL = idx+1, sDataL<<1 {
		if (sDataL & 0x80000000) != 0 {
			// Check the most significant bit.
			// Check for sprite-to-sprite collisions *before* drawing.
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
	// Draw the right half of the sprite (remaining 16 pixels).
	for ; idx < 48; idx, sDataR = idx+1, sDataR<<1 {
		if (sDataR & 0x80000000) != 0 {
			// Check the most significant bit.
			// Check for sprite-to-sprite collisions.
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
}

// drawUnexpandedMulticolor renders an unexpanded multicolor sprite onto the display buffer with collision detection.
func (sp *Sprites) drawUnexpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	// Get the foreground mask for the sprite's character column.
	foreMask := sp.collisions.GetGraphicsL(m, s)
	// Convert sprite data to bit-planes.  No expansion is needed here since the sprite is *not* expanded.
	p0 := sData & 0x55555555 // sprite to bit-planes 0.
	p1 := sData & 0xaaaaaaaa // sprite to bit-Planes 1.
	plane0 := p0 | (p0 << 1) // Combine bits for plane 0.
	plane1 := p1 | (p1 >> 1) // Combine bits for plane 1.
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			plane0 &= ^foreMask // Mask plane 0.
			plane1 &= ^foreMask // Mask plane 1.
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < 24; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		var selectedColor uint8
		// Determine the color of the current pixel based on the bit-planes.
		if (plane1 & 0x80000000) != 0 {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm1] // 11: Use color from MM1 register.
			} else {
				selectedColor = sColor // 10: Use the sprite's individual color.
			}
		} else {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm0] // 01: Use color from MM0 register.
			} else {
				continue // 00: Transparent - don't draw anything.
			}
		}
		// Check for sprite-to-sprite collisions *before* drawing.
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
}

// drawUnexpandedStandard renders a non-expanded standard sprite onto the display buffer and manages collision logic.
func (sp *Sprites) drawUnexpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	// Get the foreground mask for the sprite's character column.
	foreMask := sp.collisions.GetGraphicsL(m, s)
	// Check for collisions with the foreground.
	if (foreMask & sData) != 0 {
		// Set the sprite-to-graphics collision flag.
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			// If sprite-to-background priority is enabled, mask the sprite data.
			sData &= ^foreMask // Mask the sprite data.
		}
	}
	// Draw the sprite (24 pixels).
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			// Check the most significant bit.
			// Check for sprite-to-sprite collisions *before* drawing.
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
}
