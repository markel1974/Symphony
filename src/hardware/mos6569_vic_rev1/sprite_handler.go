package mos6569

import (
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// bitSprite0 to bitSprite7 represent bitmask values for individual sprite identifiers.
const (
	bitSprite0 = 0x1
	bitSprite1 = 0x2
	bitSprite2 = 0x4
	bitSprite3 = 0x8
	bitSprite4 = 0x10
	bitSprite5 = 0x20
	bitSprite6 = 0x40
	bitSprite7 = 0x80
)

// SpriteHandler manages sprite functionalities within the VIC-II system, including rendering, collisions, and DMA operations.
type SpriteHandler struct {
	core         *VIC        // core represents the pointer to the VIC instance used by the SpriteHandler for rendering and system control.
	collisions   *Collisions // collisions manages sprite collision detection within the SpriteHandler.
	sprites      []*Sprite   // sprites stores pointers to all Sprite instances managed by the SpriteHandler.
	dmaFlags     uint8       // dmaFlags represents the active Direct Memory Access (DMA) flags for sprite operations in the current cycle.
	displayFlags uint8       // displayFlags represents the active display flags for sprites, updated based on their DMA state and counters.
	spriteFlags  uint8       // spriteFlags represents the combined state of sprite display activity for the current line in the VIC-II system.
	offset       int         // offset is the horizontal offset used during sprite rendering to determine the starting position on the scanline.
}

// NewSprites initializes and returns a new instance of the SpriteHandler struct with default settings and allocations.
// It sets up sprite data, counters, and dependencies using the provided VIC core, collisions, and display buffer.
func NewSprites(core *VIC, collisions *Collisions, displayBuffer references.IDisplayBuffer) *SpriteHandler {
	s := &SpriteHandler{
		core:         core,
		collisions:   collisions,
		displayFlags: 0,
		dmaFlags:     0,
		offset:       0,
		sprites:      make([]*Sprite, SpriteNumber),
	}
	for i := range s.sprites {
		s.sprites[i] = NewSprite(core, displayBuffer, uint8(i), len(s.sprites))
	}
	return s
}

// Setup initializes the SpriteHandler instance, preparing internal state and configurations needed for sprite operations.
func (sp *SpriteHandler) Setup() {
	// Nothing to do here at the moment, as initialization is handled in NewSprites.
	// This function is kept for consistency and potential future use.
}

// SetOffset updates the `offset` value of the SpriteHandler instance with the given value.
// This offset is used to calculate the starting position for sprite rendering on the current scanline.
func (sp *SpriteHandler) SetOffset(offset int) {
	sp.offset = offset
}

// GetDMAFlag checks and returns the active DMA flag for the specified sprite(s) by performing a bitwise AND operation.
func (sp *SpriteHandler) GetDMAFlag(b uint8) uint8 {
	return sp.dmaFlags & b
}

// FetchPtr fetches the sprite pointer for the given sprite number if BA and AEC conditions are met,
// and updates its data pointer. Logs a warning if conditions are not met.
// This function is called during specific VIC-II cycles when sprite data pointers need to be fetched
// from memory.  The actual memory address is calculated based on the sprite pointer value
// and the VIC-II's memory mapping.
func (sp *SpriteHandler) FetchPtr(num uint8) {
	if sp.core.baLow && sp.core.aecLow {
		sp.sprites[num].FetchPtr()
	} else {
		log.Printf("sprites: can't fetch sprite ptr %d", num) // Should not normally happen, as the VIC-II controls BA/AEC.
	}
}

// FetchData loads sprite data for the given sprite number and byte index from memory if BA and AEC lines are low.
// It updates the sprite's data array and increments its data counter.Logs an error if BA or AEC lines are high.
// This function is called during specific VIC-II cycles (typically cycles 49-55, three times per sprite) to fetch
// the actual sprite data (3 bytes per sprite per cycle).
func (sp *SpriteHandler) FetchData(num uint8, bNum uint8) {
	if sp.core.baLow && sp.core.aecLow {
		sp.sprites[num].FetchData(bNum)
	} else {
		log.Printf("sprites: can't fetch sprite %d - %d", num, bNum) // Should not normally happen.
	}
}

// UpdateDisplayFlags updates the display flags for sprites by checking and clearing flags based on DMA activity status.
// It determines which sprites are currently active for display based on DMA and counter-status.
// Called at the *end* of each scanline (cycle 58).
func (sp *SpriteHandler) UpdateDisplayFlags() {
	sp.spriteFlags = sp.displayFlags
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		if ((sp.displayFlags & mask) != 0) && ((sp.dmaFlags & mask) == 0) {
			sp.displayFlags &= ^mask
		}
	}
}

// UpdateDMA updates the DMA status of sprites based on their raster line and enabled flags.
func (sp *SpriteHandler) UpdateDMA() {
	rasterY := sp.core.rasterY & 0xff
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		if ((sp.core.me & mask) != 0) && (rasterY == uint16(sp.core.mXy[num])) {
			sp.dmaFlags |= mask
			sprite.CounterBaseReset()
			if (sp.core.mye & mask) != 0 {
				sp.core.sprExpY &= ^mask
			}
		}
	}
}

/*
// IncrementCounterBase increments the base counter of each sprite by the provided value and updates DMA flags accordingly.
func (sp *SpriteHandler) IncrementCounterBase(increment uint16) {
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		if (sp.core.sprExpY & mask) != 0 {
			if sprite.CounterBaseIncrement(increment) {
				sp.dmaFlags &= ^mask
			}
		}
	}
}
*/

// IncrementCounterBase advances the vertical position counter of each sprite by the specified increment, handling expansion logic.
// If a sprite is in vertical expansion mode and on the second line of an expanded pair, the counter is not advanced.
// Otherwise, the counter is incremented, and if the operation completes the sprite's row, its DMA flag is cleared.
func (sp *SpriteHandler) IncrementCounterBase(increment uint16) {
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		// Check if sprite is enabled for vertical expansion
		isExpanded := (sp.core.mye & mask) != 0
		// Check the flip-flop state for sprite expansion
		isSecondLineOfPair := (sp.core.sprExpY & mask) == 0
		if isExpanded && isSecondLineOfPair {
			// do nothing
		} else {
			// Otherwise (standard sprite OR first line of an expanded pair),
			// advance to the next row of sprite data.
			if sprite.CounterBaseIncrement(increment) {
				sp.dmaFlags &= ^mask
			}
		}
	}
}

// UpdateDisplayYFlags updates the display flags for sprites based on raster line position and active DMA flags.
// Called at the beginning of each scanline (cycle 14).
func (sp *SpriteHandler) UpdateDisplayYFlags() {
	rasterY := sp.core.rasterY & 0xff
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		sprite.CounterBaseApply()
		if ((sp.dmaFlags & mask) != 0) && (rasterY == uint16(sp.core.mXy[num])) {
			sp.displayFlags |= mask
		}
	}
}

// Draw renders all active sprites for the current line based on their flags, properties, and configurations.
// It handles both expanded and unexpanded sprites in standard and multicolor modes.
// Collision detection for sprites is carried out during the rendering process.
// Called in cycles 57-62.
func (sp *SpriteHandler) Draw() {
	activeSprites := _spritesData[sp.spriteFlags]
	if activeSprites == nil {
		return
	}
	// Prepare the collision detection system for this scanline.
	sp.collisions.Prepare()
	// Draw active sprites
	for _, sNum := range activeSprites {
		sp.sprites[sNum].Draw(sp.offset, sp.collisions)
	}
	// Perform the final collision detection checks.
	sp.collisions.Detect()
}

// ModeUpdate performs a mode-specific update for all sprites managed by the SpriteHandler by invoking their ModeUpdate method.
func (sp *SpriteHandler) ModeUpdate() {
	for _, sprite := range sp.sprites {
		sprite.ModeUpdate()
	}
}
