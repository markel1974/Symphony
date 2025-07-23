package mos6569

import (
	"github.com/markel1974/c64emu/src/common/bits"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	spriteNumber = 8
	spriteIndex  = spriteNumber - 1
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
	//core               *VIC        // core represents the pointer to the VIC instance used by the SpriteHandler for rendering and system control.
	collisions         *Collisions // collisions manages sprite collision detection within the SpriteHandler.
	sprites            []*Sprite   // sprites stores pointers to all Sprite instances managed by the SpriteHandler.
	dmaFlags           uint8       // dmaFlags represents the active Direct Memory Access (DMA) flags for sprite operations in the current cycle.
	spriteFlags        uint8       // spriteFlags represents the combined state of sprite display activity for the current line in the VIC-II system.
	currentSpriteFlags uint8       // currentSpriteFlags represents the active display flags for sprites, updated based on their DMA state and counters.
	offset             int         // offset is the horizontal offset used during sprite rendering to determine the starting position on the scanline.
	yExpansion         uint8       // 8 sprite y expansion FlipFlops
	mXx                []uint16    // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy                []uint8     // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mXc                []uint8     // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
	mx8                uint8       // VIC register
	mmc                uint8       // VIC register
	mxe                uint8       // VIC register
	me                 uint8       // Sprite Enabled
	mye                uint8       // VIC register
	mdp                uint8       // Sprite data priority
	mm0                uint8
	mm1                uint8
}

// NewSprites initializes and returns a new instance of the SpriteHandler struct with default settings and allocations.
// It sets up sprite data, counters, and dependencies using the provided VIC core, collisions, and display buffer.
func NewSprites(parent references.IComponent, factory references.IComponentFactory, label string, instance int, memory *Memory, collisions *Collisions, displayBuffer references.IDisplayBuffer) *SpriteHandler {
	s := &SpriteHandler{
		BaseComponent:      component.NewBaseComponent(),
		collisions:         collisions,
		currentSpriteFlags: 0,
		dmaFlags:           0,
		offset:             0,
		sprites:            make([]*Sprite, spriteNumber),
		mXx:                make([]uint16, spriteNumber),
		mXy:                make([]uint8, spriteNumber),
		mXc:                make([]uint8, spriteNumber),
		mx8:                0,
		mmc:                0,
		mxe:                0,
		me:                 0,
		mye:                0,
		mdp:                0,
		mm0:                0,
		mm1:                0,
	}
	s.BaseComponent.Register(factory, parent, "spriteHandler", s, references.IdInternalComponent(label, instance, "SpriteHandler"))
	for i := range s.sprites {
		s.sprites[i] = NewSprite(s, factory, "SpriteHandler", i, displayBuffer, memory, collisions, uint8(i), len(s.sprites))
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

// ReadMM0 returns the mm0 value of the sprite with the upper nibble set to 0xf.
func (sp *SpriteHandler) ReadMM0() uint8 {
	return sp.mm0 | 0xf0
}

// ReadMM1 returns the result of the bitwise OR operation between the mm1 field and the constant 0xf0.
func (sp *SpriteHandler) ReadMM1() uint8 {
	return sp.mm1 | 0xf0
}

// ReadMe returns the current value of the sprite enable flag stored in the 'me' field.
func (sp *SpriteHandler) ReadMe() uint8 {
	return sp.me // Sprite enable
}

// ReadMDp retrieves the sprite data priority (MDp) from the SpriteHandler and returns it as an unsigned 8-bit integer.
func (sp *SpriteHandler) ReadMDp() uint8 {
	return sp.mdp // Sprite data priority
}

// ReadMYe retrieves the current Sprite Y expansion value.
func (sp *SpriteHandler) ReadMYe() uint8 {
	return sp.mye // Sprite Y expansion
}

// ReadMMc returns the sprite multicolor configuration stored in the `mmc` field of the SpriteHandler.
func (sp *SpriteHandler) ReadMMc() uint8 {
	return sp.mmc // Sprite multicolor
}

// ReadMXe retrieves the current state of the sprite's X expansion flag.
func (sp *SpriteHandler) ReadMXe() uint8 {
	return sp.mxe // Sprite X expansion
}

// ReadMXc0 reads and returns the first element of the mXc array combined with a fixed mask of 0xf0.
func (sp *SpriteHandler) ReadMXc0() uint8 {
	return sp.mXc[0] | 0xf0
}

// ReadMXc1 returns the second element of the mXc array bitwise OR'd with 0xf0.
func (sp *SpriteHandler) ReadMXc1() uint8 {
	return sp.mXc[1] | 0xf0
}

// ReadMXc2 reads the third value from the mXc array and combines it with 0xf0 using bitwise OR, returning the result.
func (sp *SpriteHandler) ReadMXc2() uint8 {
	return sp.mXc[2] | 0xf0
}

// ReadMXc3 reads the fourth element of the mXc array, applies a bitwise OR with 0xf0, and returns the resulting value.
func (sp *SpriteHandler) ReadMXc3() uint8 {
	return sp.mXc[3] | 0xf0
}

// ReadMXc4 returns the value of the fifth element of the mXc array OR'ed with 0xf0.
func (sp *SpriteHandler) ReadMXc4() uint8 {
	return sp.mXc[4] | 0xf0
}

// ReadMXc5 retrieves and returns the modified value of the sixth multi-color register for the sprite system.
func (sp *SpriteHandler) ReadMXc5() uint8 {
	return sp.mXc[5] | 0xf0
}

// ReadMXc6 retrieves the seventh sprite's X-coordinate high bit, combined with a fixed high nibble value of 0xf.
func (sp *SpriteHandler) ReadMXc6() uint8 {
	return sp.mXc[6] | 0xf0
}

// ReadMXc7 retrieves the 8th sprite's collision configuration, combining it with a constant value of 0xf0.
func (sp *SpriteHandler) ReadMXc7() uint8 {
	return sp.mXc[7] | 0xf0
}

// ReadMXx0 returns the least-significant byte of the X-coordinate for the first sprite stored in the mXx array.
func (sp *SpriteHandler) ReadMXx0() uint8 {
	return uint8(sp.mXx[0])
}

// ReadMXy0 retrieves the first element of the mXy array, representing the Y-coordinate-related sprite data.
func (sp *SpriteHandler) ReadMXy0() uint8 {
	return sp.mXy[0]
}

// ReadMXx1 retrieves the most significant bit of the x-coordinate for sprite 1 from the internal mXx array.
func (sp *SpriteHandler) ReadMXx1() uint8 {
	return uint8(sp.mXx[1])
}

// ReadMXy1 retrieves the Y position of the second sprite from the mXy array.
func (sp *SpriteHandler) ReadMXy1() uint8 {
	return sp.mXy[1]
}

// ReadMXx2 retrieves the X-coordinate of the third sprite from the internal mXx array as an unsigned 8-bit integer.
func (sp *SpriteHandler) ReadMXx2() uint8 {
	return uint8(sp.mXx[2])
}

// ReadMXy2 retrieves the Y-coordinate value of the third sprite (index 2) from the internal mXy array.
func (sp *SpriteHandler) ReadMXy2() uint8 {
	return sp.mXy[2]
}

// ReadMXx3 retrieves the value of the fourth element in the mXx slice of SpriteHandler as an unsigned 8-bit integer.
func (sp *SpriteHandler) ReadMXx3() uint8 {
	return uint8(sp.mXx[3])
}

// ReadMXy3 retrieves the Y coordinate value of the third sprite from the internal mXy array.
// It returns a uint8 representing the Y coordinate.
func (sp *SpriteHandler) ReadMXy3() uint8 {
	return sp.mXy[3]
}

// ReadMXx4 retrieves the fourth horizontal position value for a sprite from the mXx array in the SpriteHandler.
func (sp *SpriteHandler) ReadMXx4() uint8 {
	return uint8(sp.mXx[4])
}

// ReadMXy4 retrieves the Y-coordinate value for the sprite stored at index 4 in the mXy array.
func (sp *SpriteHandler) ReadMXy4() uint8 {
	return sp.mXy[4]
}

// ReadMXx5 returns the value of the 6th element in the mXx array as an unsigned 8-bit integer.
func (sp *SpriteHandler) ReadMXx5() uint8 {
	return uint8(sp.mXx[5])
}

// ReadMXy5 retrieves the Y-coordinate of the sixth sprite from the mXy array in the SpriteHandler struct.
func (sp *SpriteHandler) ReadMXy5() uint8 {
	return sp.mXy[5]
}

// ReadMXx6 retrieves the X-coordinate value of the 6th sprite from the internal mXx array of the SpriteHandler.
func (sp *SpriteHandler) ReadMXx6() uint8 {
	return uint8(sp.mXx[6])
}

// ReadMXy6 retrieves the Y-coordinate value of the 6th sprite from the mXy array.
func (sp *SpriteHandler) ReadMXy6() uint8 {
	return sp.mXy[6]
}

// ReadMXx7 returns the 8th element of the mXx array as a uint8 representing the X-coordinate of sprite 7.
func (sp *SpriteHandler) ReadMXx7() uint8 {
	return uint8(sp.mXx[7])
}

// ReadMXy7 retrieves the value stored at index 7 of the mXy array in the SpriteHandler.
func (sp *SpriteHandler) ReadMXy7() uint8 {
	return sp.mXy[7]
}

// ReadMX8 retrieves the most significant bit (MSB) of sprite X positions.
// Returns the value of the mx8 field as an 8-bit unsigned integer.
func (sp *SpriteHandler) ReadMX8() uint8 {
	// Sprite X position MSB
	return sp.mx8
}

// WriteMM0 assigns the provided 8-bit unsigned integer value to the mm0 field of the SpriteHandler instance.
func (sp *SpriteHandler) WriteMM0(data uint8) {
	sp.mm0 = data
}

// WriteMM1 sets the value of the mm1 property in the SpriteHandler instance using the provided 8-bit unsigned integer.
func (sp *SpriteHandler) WriteMM1(data uint8) {
	sp.mm1 = data
}

// WriteMe sets the sprite enable register to the provided data value.
func (sp *SpriteHandler) WriteMe(data uint8) {
	sp.me = data // Sprite enable
}

// WriteMYe sets the sprite Y expansion register with the provided data.
// Updates the internal `yExpansion` state based on the input value.
func (sp *SpriteHandler) WriteMYe(data uint8) {
	sp.mye = data // Sprite Y expansion
	sp.yExpansion |= ^data
}

// WriteMDp sets the sprite data priority (mdp) to the provided uint8 value.
func (sp *SpriteHandler) WriteMDp(data uint8) {
	sp.mdp = data // Sprite data priority
}

// WriteMMc updates the mmc field and calls ModeUpdate on all sprites to refresh their rendering mode with updated flags.
func (sp *SpriteHandler) WriteMMc(data uint8) {
	sp.mmc = data
	for _, sprite := range sp.sprites {
		sprite.ModeUpdate(sp.mmc, sp.mxe)
	}
}

// WriteMXe updates the horizontal expansion flag (MXE) for sprites and triggers their mode update accordingly.
func (sp *SpriteHandler) WriteMXe(data uint8) {
	sp.mxe = data
	for _, sprite := range sp.sprites {
		sprite.ModeUpdate(sp.mmc, sp.mxe)
	}
}

// WriteMXx0 updates the lower byte of the first mXx element with the provided data while preserving the upper byte.
func (sp *SpriteHandler) WriteMXx0(data uint8) {
	sp.mXx[0] = (sp.mXx[0] & 0xff00) | uint16(data)
}

// WriteMXy0 updates the first element of the SpriteHandler's mXy slice with the provided data value.
func (sp *SpriteHandler) WriteMXy0(data uint8) {
	sp.mXy[0] = data
}

// WriteMXx1 updates the lower 8 bits of the mXx[1] field with the provided data value, preserving the upper 8 bits.
func (sp *SpriteHandler) WriteMXx1(data uint8) {
	sp.mXx[1] = (sp.mXx[1] & 0xff00) | uint16(data)
}

// WriteMXy1 writes the given uint8 data to the second index (1) of the mXy slice in the SpriteHandler struct.
func (sp *SpriteHandler) WriteMXy1(data uint8) {
	sp.mXy[1] = data
}

// WriteMXx2 updates the lower byte of the third element in the mXx array with the provided uint8 data value.
func (sp *SpriteHandler) WriteMXx2(data uint8) {
	sp.mXx[2] = (sp.mXx[2] & 0xff00) | uint16(data)
}

// WriteMXy2 sets the third index of the mXy slice to the specified data value.
func (sp *SpriteHandler) WriteMXy2(data uint8) {
	sp.mXy[2] = data
}

// WriteMXx3 updates the lowest 8 bits of the fourth element in the mXx array with the provided data.
func (sp *SpriteHandler) WriteMXx3(data uint8) {
	sp.mXx[3] = (sp.mXx[3] & 0xff00) | uint16(data)
}

// WriteMXy3 sets the value of the fourth element in the mXy slice to the specified data.
func (sp *SpriteHandler) WriteMXy3(data uint8) {
	sp.mXy[3] = data
}

// WriteMXx4 updates the lower 8 bits of the fifth element in the mXx array with the given uint8 data value.
func (sp *SpriteHandler) WriteMXx4(data uint8) {
	sp.mXx[4] = (sp.mXx[4] & 0xff00) | uint16(data)
}

// WriteMXy4 sets the fifth element of the mXy array in SpriteHandler to the provided uint8 data value.
func (sp *SpriteHandler) WriteMXy4(data uint8) {
	sp.mXy[4] = data
}

// WriteMXx5 writes an 8-bit value to the lower 8 bits of the 6th element in the mXx array.
func (sp *SpriteHandler) WriteMXx5(data uint8) {
	sp.mXx[5] = (sp.mXx[5] & 0xff00) | uint16(data)
}

// WriteMXy5 updates the sixth element (index 5) of the mXy slice with the provided data value.
func (sp *SpriteHandler) WriteMXy5(data uint8) {
	sp.mXy[5] = data
}

// WriteMXx6 writes the provided 8-bit data to the lower byte of the mXx[6] field, preserving the upper byte.
func (sp *SpriteHandler) WriteMXx6(data uint8) {
	sp.mXx[6] = (sp.mXx[6] & 0xff00) | uint16(data)
}

// WriteMXy6 sets the seventh element of the mXy array to the specified data value.
func (sp *SpriteHandler) WriteMXy6(data uint8) {
	sp.mXy[6] = data
}

// WriteMXx7 updates the lower 8 bits of the 7th element in the mXx array with the provided data while preserving the upper bits.
func (sp *SpriteHandler) WriteMXx7(data uint8) {
	sp.mXx[7] = (sp.mXx[7] & 0xff00) | uint16(data)
}

// WriteMXy7 sets the 8th index of the mXy slice in the SpriteHandler to the provided uint8 value.
func (sp *SpriteHandler) WriteMXy7(data uint8) {
	sp.mXy[7] = data
}

// WriteMX8 updates the MSBs of X coordinates for sprite handling based on the given data.
func (sp *SpriteHandler) WriteMX8(data uint8) { //MSBs of X coordinates
	sp.mx8 = data
	for i := range sp.mXx {
		if (data & bits.Uint8s[i]) != 0 {
			sp.mXx[i] |= 0x100
		} else {
			sp.mXx[i] &= 0xff
		}
	}
}

// WriteMXc0 assigns the provided uint8 data value to the first element of the mXc array in the SpriteHandler.
func (sp *SpriteHandler) WriteMXc0(data uint8) {
	sp.mXc[0] = data
}

// WriteMXc1 sets the second element of the mXc slice to the provided data value.
func (sp *SpriteHandler) WriteMXc1(data uint8) {
	sp.mXc[1] = data
}

// WriteMXc2 sets the third element (index 2) of the mXc array in the SpriteHandler to the specified data value.
func (sp *SpriteHandler) WriteMXc2(data uint8) {
	sp.mXc[2] = data
}

// WriteMXc3 assigns the given data to the fourth index of the mXc array in the SpriteHandler.
func (sp *SpriteHandler) WriteMXc3(data uint8) {
	sp.mXc[3] = data
}

// WriteMXc4 writes a uint8 value to the fifth element of the mXc array within the SpriteHandler struct.
func (sp *SpriteHandler) WriteMXc4(data uint8) {
	sp.mXc[4] = data
}

// WriteMXc5 updates the fifth element of the mXc slice with the provided 8-bit unsigned integer data.
func (sp *SpriteHandler) WriteMXc5(data uint8) {
	sp.mXc[5] = data
}

// WriteMXc6 writes the provided data to the 6th index of the mXc slice in the SpriteHandler.
func (sp *SpriteHandler) WriteMXc6(data uint8) {
	sp.mXc[6] = data
}

// WriteMXc7 sets the 7th index of the mXc array to the provided uint8 data value.
func (sp *SpriteHandler) WriteMXc7(data uint8) {
	sp.mXc[7] = data
}

// UpdateYExpansion adjusts the sprite's vertical expansion state based on the MYE register using an inversion technique.
func (sp *SpriteHandler) UpdateYExpansion() {
	// Invert y expansion FlipFlop (if MYE bit is set)
	for idx, mask := 0, uint8(1); idx < spriteNumber; idx, mask = idx+1, mask<<1 {
		if (sp.mye & mask) != 0 {
			sp.yExpansion ^= mask
		}
	}
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
	sp.sprites[num].FetchPtr()
}

// FetchData0 fetches the first byte of data for the sprite identified by sNum and stores it in the sprite's data buffer.
func (sp *SpriteHandler) FetchData0(sNum uint8) {
	sp.sprites[sNum].FetchData(0)
}

// FetchData1 fetches the second byte of data for the sprite identified by sNum and stores it in the sprite's data buffer.
func (sp *SpriteHandler) FetchData1(sNum uint8) {
	sp.sprites[sNum].FetchData(1)
}

// FetchData2 retrieves the 'byte 2' data for the specified sprite and computes additional data configuration for rendering.
func (sp *SpriteHandler) FetchData2(sNum uint8) {
	sp.sprites[sNum].FetchData(2)

	sp.sprites[sNum].SetData(sp.mdp, sp.mm0, sp.mm1, sp.mXc[sNum], sp.mXx[sNum])
}

// UpdateDMA updates the Direct Memory Access (DMA) flags for sprites at the current raster line position.
// It resets the counter-base for active sprites and handles vertical expansion mode and line-specific configurations.
func (sp *SpriteHandler) UpdateDMA(rasterY uint16) {
	rasterYLow := uint8(rasterY & 0xff)
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		if ((sp.me & mask) != 0) && (rasterYLow == sp.mXy[num]) {
			sp.dmaFlags |= mask
			sprite.ResetCounterBase()
			if (sp.mye & mask) != 0 {
				sp.yExpansion &= ^mask
			}
		}
	}
}

// TryIncrementCounterBase checks sprite vertical expansion state and conditionally increments the sprite's internal row counter.
func (sp *SpriteHandler) TryIncrementCounterBase() {
	for _, sprite := range sp.sprites {
		mask := sprite.Mask()
		// Check if sprite is enabled for vertical expansion
		isExpanded := (sp.mye & mask) != 0
		// Check the flip-flop state for sprite expansion
		isSecondLineOfPair := (sp.yExpansion & mask) == 0
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
func (sp *SpriteHandler) PrepareSpriteFlags(rasterY uint16) {
	rasterYLow := uint8(rasterY & 0xff)
	for _, sprite := range sp.sprites {
		num := sprite.Number()
		mask := sprite.Mask()
		sprite.CommitCounterBase()
		if ((sp.dmaFlags & mask) != 0) && (rasterYLow == sp.mXy[num]) {
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
	const spriteIndex = spriteNumber - 1
	activeSprites := _spritesData[sp.spriteFlags]
	if activeSprites == nil {
		return
	}
	// Prepare the collision detection system for this scanline.
	sp.collisions.Prepare()
	for _, sNum := range activeSprites {
		sp.sprites[sNum].Draw(sp.offset)
	}
	// Perform the final collision detection checks.
	sp.collisions.Commit()
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
