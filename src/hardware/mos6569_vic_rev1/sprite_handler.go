package mos6569

import (
	"github.com/markel1974/c64emu/src/references"
	"log"
)

// SpriteHandler represent the structure responsible for handling sprites, including their data, states, and configurations.
// It contains properties for managing sprite visual data, collision detection, and display control buffers.
// The type relies on various counters, flags, and pointers to handle sprite DMA and display activities effectively.
// It interacts with the VIC core and an implemented display buffer interface for rendering and collision processing.
type SpriteHandler struct {
	core         *VIC        // Pointer to the main VIC-II core.
	collisions   *Collisions // Pointer to the collision detection system.
	sprites      []*Sprite   // Pointers to all Sprite objects managed by the SpriteHandler system.
	dmaFlags     uint8       // Active DMA Sprite (bitmask: bit i = 1 means sprite is active).
	displayFlags uint8       // Active Display Sprite (bitmask).
	spriteFlags  uint8       // Sprite in this line (bitmask).
	offset       int         // Offset from bitmap spritesBuffer
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

// UpdateCounterBase increments the base counters of sprites with vertical expansion enabled by 2 for each enabled sprite.
// Handles y-expansion
func (sp *SpriteHandler) UpdateCounterBase() {
	for _, sprite := range sp.sprites {
		if (sp.core.sprExpY & (1 << sprite.Number())) != 0 {
			sprite.DataCounterBaseIncrement(2)
		}
	}
}

// GetDMAFlag checks and returns the active DMA flag for the specified sprite(s) by performing a bitwise AND operation.
func (sp *SpriteHandler) GetDMAFlag(b uint8) uint8 {
	return sp.dmaFlags & b
}

// UpdateDMA updates the DMA status of sprites based on their raster line and enabled flags.
// Called in cycle 12.
func (sp *SpriteHandler) UpdateDMA() {
	rasterY := sp.core.rasterY & 0xff
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		if ((sp.core.me & mask) != 0) && (rasterY == uint16(sp.core.mXy[num])) {
			sp.dmaFlags |= mask
			sprite.ResetDataCounterBase()
			if (sp.core.mye & mask) != 0 {
				sp.core.sprExpY &= ^mask
			}
		}
	}
}

// UpdateCounterBaseDMA updates the base counters of sprites and manages DMA flags based on specific conditions.
// Called in cycle 13.
func (sp *SpriteHandler) UpdateCounterBaseDMA() {
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		if (sp.core.sprExpY & mask) != 0 {
			sprite.DataCounterBaseIncrement(1)
		}
		if (sprite.DataCounterBase() & 0x3f) == 0x3f {
			sp.dmaFlags &= ^mask
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
		sprite.ApplyDataCounterBase()
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
	for _, sNum := range activeSprites {
		sp.sprites[sNum].Draw(sp.offset, sp.collisions)
	}
	// Perform the final collision detection checks.
	sp.collisions.Detect()
}
