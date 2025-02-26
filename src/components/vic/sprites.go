package mos6569

import (
	"log"
)

// Sprites represents the structure responsible for handling sprites, including their data, states, and configurations.
// It contains properties for managing sprite visual data, collision detection, and display control buffers.
// The type relies on various counters, flags, and pointers to handle sprite DMA and display activities effectively.
// It interacts with the VIC core and an implemented display buffer interface for rendering and collision processing.
type Sprites struct {
	core            *VIC
	displayBuffer   IDisplayBuffer
	collisions      *Collisions
	dataPtr         []uint16  // Sprite data pointers
	data            [][]uint8 // Sprite data
	dmaFlags        uint8     // active DMA Sprite
	displayFlags    uint8     // active Display Sprite
	spriteFlags     uint8     // Sprite in this line
	dataCounter     []uint16  // Sprite counter data
	dataCounterBase []uint16  // Sprite base counter data
	offset          int       //
}

// NewSprites initializes and returns a new instance of the Sprites struct with default settings and allocations.
// It sets up sprite data, counters, and dependencies using the provided VIC core, collisions, and display buffer.
func NewSprites(core *VIC, collisions *Collisions, db IDisplayBuffer) *Sprites {
	s := &Sprites{
		core:            core,
		displayBuffer:   db,
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
		s.data[i] = make([]uint8, 4)
	}
	for i := range s.dataCounter {
		s.dataCounter[i] = 63
	}
	return s
}

// Setup initializes the Sprites instance, preparing internal state and configurations needed for sprite operations.
func (sp *Sprites) Setup() {
}

// SetOffset updates the `offset` value of the Sprites instance with the given value.
func (sp *Sprites) SetOffset(offset int) {
	sp.offset = offset
}

// FetchPtr fetches the sprite pointer for the given sprite number if BA and AEC conditions are met, and updates its data pointer.
// Logs a warning if conditions are not met.
func (sp *Sprites) FetchPtr(num int) {
	if sp.core.baLow && sp.core.aecLow {
		addr := sp.core.matrixBase | 0x03f8 | uint16(num)
		data := sp.core.ReadByte(addr)
		ptr := uint16(data) << 6
		sp.dataPtr[num] = ptr
	} else {
		log.Printf("sprites: can't fetch sprite ptr %d", num)
	}
}

// Fetch loads sprite data for the given sprite number and byte index from memory if BA and AEC lines are low.
// It updates the sprite's data array and increments its data counter.
// Logs an error if BA or AEC lines are high.
func (sp *Sprites) Fetch(num int, bNum int) {
	if sp.core.baLow && sp.core.aecLow {
		ptr := sp.dataPtr[num]
		addr := (sp.dataCounter[num] & 0x3f) | ptr
		data := sp.core.ReadByte(addr)
		sp.data[num][bNum] = data
		sp.dataCounter[num]++
	} else {
		log.Printf("sprites: can't fetch sprite %d - %d", num, bNum)
	}
}

// UpdateDisplayFlags updates the display flags for sprites by checking and clearing flags based on DMA activity status.
func (sp *Sprites) UpdateDisplayFlags() {
	sp.spriteFlags = sp.displayFlags
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if ((sp.displayFlags & mask) != 0) && ((sp.dmaFlags & mask) == 0) {
			sp.displayFlags &= ^mask
		}
	}
}

// UpdateCounterBase increments the base counters of sprites with vertical expansion enabled by 2 for each enabled sprite.
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
func (sp *Sprites) Draw() {
	if sp.spriteFlags == 0 {
		return
	}

	sp.collisions.Prepare()

	for sNum, sBit := uint8(0), uint8(1); sNum < SpriteNumber; sNum, sBit = sNum+1, sBit<<1 {
		if (sp.spriteFlags & sBit) != 0 {
			sColor := _colors[sp.core.mXc[sNum]]
			sData := (uint32(sp.data[sNum][0]) << 24) | (uint32(sp.data[sNum][1]) << 16) | (uint32(sp.data[sNum][2]) << 8)
			sOffset := int(sp.core.mXx[sNum]) + SpriteNumber
			lineOffset := sp.offset + sOffset // lineStart + sOffset
			m := sOffset / SpriteNumber
			s := sOffset & 7
			multiColor := (sp.core.mmc & sBit) != 0
			if expanded := (sp.core.mxe & sBit) != 0; expanded {
				if multiColor {
					sp.drawExpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit)
				} else {
					sp.drawExpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit)
				}
			} else {
				if multiColor {
					sp.drawUnexpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit)
				} else {
					sp.drawUnexpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit)
				}
			}
		}
	}

	sp.collisions.Detect()
}

// drawExpandedMulticolor renders expanded multicolor sprites with applied graphics and collision handling logic.
func (sp *Sprites) drawExpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	foreMaskL := sp.collisions.GetGraphicsL(m, s)
	foreMaskR := sp.collisions.GetGraphicsR(m, s)
	// Expand sprite
	sDataL := (uint32(_multiExpTable[(sData>>24)&0xff]) << 16) | (uint32(_multiExpTable[(sData>>16)&0xff]))
	sDataR := uint32(_multiExpTable[(sData>>8)&0xff]) << 16
	plane0L := (sDataL & 0x55555555) | ((sDataL & 0x55555555) << 1) // convert sprite to bitPlanes
	plane1L := (sDataL & 0xaaaaaaaa) | ((sDataL & 0xaaaaaaaa) >> 1) // convert sprite to bitPlanes
	plane0R := (sDataR & 0x55555555) | ((sDataR & 0x55555555) << 1) // convert sprite to bitPlanes
	plane1R := (sDataR & 0xaaaaaaaa) | ((sDataR & 0xaaaaaaaa) >> 1) // convert sprite to bitPlanes
	// Collision with graphics?
	if ((foreMaskL & (plane0L | plane1L)) != 0) || ((foreMaskR & (plane0R | plane1R)) != 0) {
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			plane0L &= ^foreMaskL //background, mask sprite
			plane1L &= ^foreMaskL //background, mask sprite
			plane0R &= ^foreMaskR //background, mask sprite
			plane1R &= ^foreMaskR //background, mask sprite
		}
	}
	idx := 0
	for ; idx < 32; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		selectedColor := uint8(0)
		if (plane1L & 0x80000000) != 0 {
			if (plane0L & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm1]
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0L & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm0]
			} else {
				continue
			}
		}
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
	for ; idx < 48; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		selectedColor := uint8(0)
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
				continue
			}
		}
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
}

// drawExpandedStandard draws an expanded standard sprite on the display buffer with collision detection and masking checks.
func (sp *Sprites) drawExpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	foreMaskL := sp.collisions.GetGraphicsL(m, s)
	foreMaskR := sp.collisions.GetGraphicsR(m, s)
	sDataL := uint32(_expTable[(sData>>24)&0xff])<<16 | uint32(_expTable[(sData>>16)&0xff])
	sDataR := uint32(_expTable[(sData>>8)&0xff]) << 16
	if ((foreMaskL & sDataL) != 0) || ((foreMaskR & sDataR) != 0) {
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			sDataL &= ^foreMaskL
			sDataR &= ^foreMaskR
		}
	}
	var idx = 0
	for ; idx < 32; idx, sDataL = idx+1, sDataL<<1 {
		if (sDataL & 0x80000000) != 0 {
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
	for ; idx < 48; idx, sDataR = idx+1, sDataR<<1 {
		if (sDataR & 0x80000000) != 0 {
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
}

// drawUnexpandedMulticolor renders an unexpanded multicolor sprite onto the display buffer with collision detection.
func (sp *Sprites) drawUnexpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	foreMask := sp.collisions.GetGraphicsL(m, s)
	p0 := sData & 0x55555555 // sprite to bitPlanes
	p1 := sData & 0xaaaaaaaa // sprite to bitPlanes
	plane0 := p0 | (p0 << 1)
	plane1 := p1 | (p1 >> 1)
	// check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		sp.collisions.SetGraphicsCollision(sBit)
		if (sp.core.mdp & sBit) != 0 {
			plane0 &= ^foreMask //background, mask sprite
			plane1 &= ^foreMask //background, mask sprite
		}
	}
	for idx := 0; idx < 24; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		var selectedColor uint8
		if (plane1 & 0x80000000) != 0 {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm1]
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = _colors[sp.core.mm0]
			} else {
				continue
			}
		}
		if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
			sp.displayBuffer.Set(lineOffset+idx, selectedColor)
		}
	}
}

// drawUnexpandedStandard renders a non-expanded standard sprite onto the display buffer and manages collision logic.
func (sp *Sprites) drawUnexpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8) {
	foreMask := sp.collisions.GetGraphicsL(m, s)
	if (foreMask & sData) != 0 {
		sp.collisions.SetGraphicsCollision(sBit)
		if sp.core.mdp&sBit != 0 {
			sData &= ^foreMask //background, mask sprite
		}
	}
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			if !sp.collisions.SetSpriteCollision(sOffset+idx, sBit) {
				sp.displayBuffer.Set(lineOffset+idx, sColor)
			}
		}
	}
}
