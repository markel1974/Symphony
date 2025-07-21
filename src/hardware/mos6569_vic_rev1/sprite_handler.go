package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
	"log"
)

const (
	spriteNumber = 8
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
	*component.BaseComponent
	core               *VIC        // core represents the pointer to the VIC instance used by the SpriteHandler for rendering and system control.
	collisions         *Collisions // collisions manages sprite collision detection within the SpriteHandler.
	sprites            []*Sprite   // sprites stores pointers to all Sprite instances managed by the SpriteHandler.
	dmaFlags           uint8       // dmaFlags represents the active Direct Memory Access (DMA) flags for sprite operations in the current cycle.
	spriteFlags        uint8       // spriteFlags represents the combined state of sprite display activity for the current line in the VIC-II system.
	currentSpriteFlags uint8       // currentSpriteFlags represents the active display flags for sprites, updated based on their DMA state and counters.
	offset             int         // offset is the horizontal offset used during sprite rendering to determine the starting position on the scanline.
}

// NewSprites initializes and returns a new instance of the SpriteHandler struct with default settings and allocations.
// It sets up sprite data, counters, and dependencies using the provided VIC core, collisions, and display buffer.
func NewSprites(parent references.IComponent, factory references.IComponentFactory, label string, instance int, core *VIC, collisions *Collisions, displayBuffer references.IDisplayBuffer) *SpriteHandler {
	s := &SpriteHandler{
		BaseComponent:      component.NewBaseComponent(),
		core:               core,
		collisions:         collisions,
		currentSpriteFlags: 0,
		dmaFlags:           0,
		offset:             0,
		sprites:            make([]*Sprite, spriteNumber),
	}
	s.BaseComponent.Register(factory, parent, "spriteHandler", s, references.IdInternalComponent(label, instance, "SpriteHandler"))
	for i := range s.sprites {
		s.sprites[i] = NewSprite(s, factory, "SpriteHandler", i, core, displayBuffer, uint8(i), len(s.sprites))
	}
	return s
}

// Setup initializes the SpriteHandler instance, preparing internal state and configurations needed for sprite operations.
func (sp *SpriteHandler) Setup() error {
	// Nothing to do here at the moment, as initialization is handled in NewSprites.
	// This function is kept for consistency and potential future use.
	return nil
}

// Connect establishes and initializes the necessary connections for the SpriteHandler. Returns an error if connection fails.
func (sp *SpriteHandler) Connect() error {
	return nil
}

// EmulationRequired determines if sprite emulation logic is required for the current system configuration.
func (sp *SpriteHandler) EmulationRequired() bool {
	return false
}

// Emulate executes a single emulation step for sprite processing, handling rendering, collisions, and DMA updates.
func (sp *SpriteHandler) Emulate() {
}

// Internal determines whether the SpriteHandler operates in an internal, self-contained mode, returning a boolean value.
func (sp *SpriteHandler) Internal() bool {
	return true
}

// Reset reinitializes the SpriteHandler's internal state and configurations to its default settings.
func (sp *SpriteHandler) Reset() {
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

// UpdateDMA updates the Direct Memory Access (DMA) flags for sprites at the current raster line position.
// It resets the counter-base for active sprites and handles vertical expansion mode and line-specific configurations.
func (sp *SpriteHandler) UpdateDMA() {
	rasterY := uint8(sp.core.rasterY & 0xff)
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		if ((sp.core.me & mask) != 0) && (rasterY == sp.core.mXy[num]) {
			sp.dmaFlags |= mask
			sprite.ResetCounterBase()
			if (sp.core.mye & mask) != 0 {
				sp.core.sprExpY &= ^mask
			}
		}
	}
}

// TryIncrementCounterBase checks sprite vertical expansion state and conditionally increments the sprite's internal row counter.
func (sp *SpriteHandler) TryIncrementCounterBase() {
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
			sprite.IncrementCounterBase()
		}
	}
}

// CommitIncrementCounterBase ensures sprite row counters are updated, and disables DMA flags for sprites reaching their row end.
func (sp *SpriteHandler) CommitIncrementCounterBase() {
	for _, sprite := range sp.sprites {
		if sprite.CommitIncrementCounterBase() {
			mask := sprite.Mask()
			sp.dmaFlags &= ^mask
		}
	}
}

// PrepareSpriteFlags updates the current sprite display flags based on DMA activity and the current raster line position.
func (sp *SpriteHandler) PrepareSpriteFlags() {
	rasterY := uint8(sp.core.rasterY & 0xff)
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		sprite.CommitCounterBase()
		if ((sp.dmaFlags & mask) != 0) && (rasterY == sp.core.mXy[num]) {
			sp.currentSpriteFlags |= mask
		}
	}
}

// CommitSpriteFlags updates the display flags for sprites by checking and clearing flags based on DMA activity status.
// It determines which sprites are currently active for display based on DMA and counter-status.
func (sp *SpriteHandler) CommitSpriteFlags() {
	sp.spriteFlags = sp.currentSpriteFlags
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		if ((sp.currentSpriteFlags & mask) != 0) && ((sp.dmaFlags & mask) == 0) {
			sp.currentSpriteFlags &= ^mask
		}
	}
}

// Draw renders all active sprites for the current line based on their flags, properties, and configurations.
// It handles both expanded and unexpanded sprites in standard and multicolor modes.
// Collision detection for sprites is carried out during the rendering process.
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
	sp.collisions.Commit()
}

// ModeUpdate performs a mode-specific update for all sprites managed by the SpriteHandler by invoking their ModeUpdate method.
func (sp *SpriteHandler) ModeUpdate() {
	for _, sprite := range sp.sprites {
		sprite.ModeUpdate()
	}
}

/*
func (sp *SpriteHandler) TryIncrementCounterBaseOld() {
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		if (sp.core.sprExpY & mask) != 0 {
			sprite.IncrementCounterBase()
		}
	}
}
*/
