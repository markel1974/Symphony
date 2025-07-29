package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	collisionsSize = 1 << 9
	collisionsMask = collisionsSize - 1
)

// CollisionsUnit encapsulate collision detection functionality between sprites and bgrState within a VIC system.
// It includes buffers for handling sprite-sprite and sprite-bgrState collision data as well as priorities.
// The struct also manages foreground masks and offsets for sprite-bgrState collision computations.
type CollisionsUnit struct {
	*component.BaseComponent
	irqEmit          func(irq uint8)
	bgrState         uint8   // bgrState represents the collision state of all graphics in the current line.
	sprState         uint8   // sprState represent the state and collision mask of all active sprites in the current line.
	sprPresence      []uint8 // sprPresence stores the state of sprite presence at each pixel, used for collision detection within the line.
	bgrBuffer        []uint8 // bgrBuffer holds the collision state for background bgrState.
	bgrBufferOffset  int     // bgrBufferOffset tracks the current position in the bgrState buffer for collision updates.
	spr2BgrClx       uint8   // spr2BgrClx sprite to background collision
	spr2SprClx       uint8   // spr2SprClx sprite to sprite collision
	sprPresenceEmpty []uint8 // sprPresenceEmpty is the initialized buffer used to reset sprite collision states for each line.
	bgrBufferEmpty   []uint8 // bgrBufferEmpty holds an initialized empty buffer for resetting or clearing bgrState collision data.
}

// NewCollisions creates and returns a new CollisionsUnit instance, associated with the given VIC core.
func NewCollisions(parent references.IComponent, factory references.IComponentFactory, label string, instance int, irqEmit func(irq uint8)) *CollisionsUnit {
	c := &CollisionsUnit{
		BaseComponent:    component.NewBaseComponent(),
		irqEmit:          irqEmit,
		sprState:         0,
		sprPresence:      make([]uint8, collisionsSize),
		sprPresenceEmpty: make([]uint8, collisionsSize),
		spr2SprClx:       0,
		spr2BgrClx:       0,
		bgrState:         0,
		bgrBuffer:        make([]uint8, collisionsSize),
		bgrBufferEmpty:   make([]uint8, collisionsSize),
		bgrBufferOffset:  0,
	}
	c.BaseComponent.Register(factory, parent, "collisionsUnit", instance, c, references.IdInternalComponent(label, instance, "CollisionsUnit"))
	return c
}

// Setup initializes the collision detection system, preparing necessary configurations for its operation.
func (c *CollisionsUnit) Setup() error {
	return nil
}

// Connect establishes necessary connections required by the CollisionsUnit instance and initializes its dependencies.
func (c *CollisionsUnit) Connect() error {
	return nil
}

// EmulationRequired determines if emulation is required for the current state of the collision system. Always returns false.
func (c *CollisionsUnit) EmulationRequired() bool {
	return false
}

// Emulate performs the collision emulation process for the current frame, updating internal state based on collisions detected.
func (c *CollisionsUnit) Emulate() {
}

// Internal checks if the collision detection system is operating in internal mode and returns a boolean value.
func (c *CollisionsUnit) Internal() bool {
	return true
}

// Reset clears all collision states and buffers, preparing the collision system for a new frame.
func (c *CollisionsUnit) Reset() {
}

// RetrieveSprite2Sprite reads and clears the sprite-to-sprite collision state, returning its value as an 8-bit unsigned integer.
func (c *CollisionsUnit) RetrieveSprite2Sprite() uint8 {
	ret := c.spr2SprClx
	c.spr2SprClx = 0 // Read and clear
	return ret
}

// RetrieveSprite2Background reads and resets the sprite-to-background collision state, returning its value as an 8-bit unsigned integer.
func (c *CollisionsUnit) RetrieveSprite2Background() uint8 {
	// Sprite-background collision
	ret := c.spr2BgrClx
	c.spr2BgrClx = 0 // Read and clear
	return ret
}

// SetSprite updates the sprite-to-sprite collision state with the specified value.
func (c *CollisionsUnit) SetSprite(data uint8) {
	c.spr2SprClx = data
}

// SetBackground sets the sprite-to-background collision state with the given value.
func (c *CollisionsUnit) SetBackground(data uint8) {
	c.spr2BgrClx = data
}

// Prepare resets collision detection states and initializes sprite collision buffers for the next frame update.
// Called at the beginning of sprite drawing on each scanline.
func (c *CollisionsUnit) Prepare() {
	// Reset the sprite-to-sprite collision result.
	c.sprState = 0
	// Reset the sprite-to-background collision result.
	c.bgrState = 0
	// Reset the sprite buffer by copying the empty buffer.  This is *much* faster than iterating.
	copy(c.sprPresence, c.sprPresenceEmpty)
}

// SetSprite2BackgroundPresence sets a collision bit for bgrState by performing a bitwise OR operation with the given bit.
// 'sBit' is a bitmask representing the sprite that collided with the background (1, 2, 4, 8, 16, 32, 64, 128).
func (c *CollisionsUnit) SetSprite2BackgroundPresence(sBit uint8) {
	// Set the corresponding bit in the 'bgrState' collision result.
	c.bgrState |= sBit
}

// SetSprite2SpritePresence checks and sets sprite collision at a specific index with a sprite bit and returns collision status.
// If a collision occurs, it updates the sprite collision state;
// otherwise, it updates the sprite buffer with the new bit.
func (c *CollisionsUnit) SetSprite2SpritePresence(index int, spriteBit uint8) bool {
	sBitPresence := c.sprPresence[index&collisionsMask]
	if sBitPresence == 0 {
		// mark this sprite as present at this pixel.
		c.sprPresence[index&collisionsMask] = spriteBit
		return false
	}
	// If any sprite is already present at this pixel...
	// Update the 'sprites' collision result with *both* the existing sprites *and* the new sprite.
	c.sprState |= sBitPresence | spriteBit
	return true
}

// CommitSprite updates the sprite-to-sprite collision state and triggers an interrupt if a collision occurred.
func (c *CollisionsUnit) CommitSprite() {
	if c.spr2SprClx != 0 {
		c.spr2SprClx |= c.sprState
	} else {
		c.spr2SprClx |= c.sprState
		c.irqEmit(irqSpriteToSpriteBit)
	}
}

// CommitBackground updates the sprite-to-background collision state and triggers an interrupt if a collision is detected.
func (c *CollisionsUnit) CommitBackground() {
	if c.spr2BgrClx != 0 {
		c.spr2BgrClx |= c.bgrState
	} else {
		c.spr2BgrClx |= c.bgrState
		c.irqEmit(irqSpriteToGraphicBit)
	}
}

// IncrementBackgroundOffset increments the bgrBufferOffset field
// to track the bgrState buffer position during updates.
func (c *CollisionsUnit) IncrementBackgroundOffset() {
	// Increment the offset (in bytes).
	c.bgrBufferOffset++
}

// ClearBackground resets the bgrState collision buffer and its offset to their initial empty states.
// Called at the *beginning* of each frame.
func (c *CollisionsUnit) ClearBackground() {
	c.bgrBufferOffset = 0
	copy(c.bgrBuffer, c.bgrBufferEmpty)
}

// UpdateBackground updates the bgrState collision buffer with provided foreground mask values a and b at the current offset.
// Called during *background* rendering.  'a' and 'b' represent the pixel data for two *consecutive* pixels.
func (c *CollisionsUnit) UpdateBackground(a uint8, b uint8) {
	c.bgrBuffer[c.bgrBufferOffset] |= a   // Reset the bgrState buffer by copying the empty buffer.
	c.bgrBuffer[c.bgrBufferOffset+1] |= b // Reset the offset.
}

// GetGraphicsL computes a 32-bit bgrState mask from the bgrState buffer starting at 'charColumn' with a 'pixelOffset' adjustment.
// The result is shifted and combined to align with the intended sub-pixel bgrState position, returning a 32-bit uint representation.
func (c *CollisionsUnit) GetGraphicsL(charColumn int, pixelOffset int) uint32 {
	f := (((uint32(c.bgrBuffer[charColumn]) << 24) | (uint32(c.bgrBuffer[charColumn+1]) << 16) | (uint32(c.bgrBuffer[charColumn+2]) << 8) | (uint32(c.bgrBuffer[charColumn+3]))) << pixelOffset) | (uint32(c.bgrBuffer[charColumn+4]) >> (8 - pixelOffset))
	return f
}

// GetGraphicsR calculates a 32-bit bgrState mask from the bgrState buffer starting at 'charColumn+4'.
// It shifts the mask by 'pixelOffset' bits to align with the intended sub-pixel bgrState position.
// Returns a composite 32-bit unsigned integer representing the derived bgrState mask.
func (c *CollisionsUnit) GetGraphicsR(charColumn int, pixelOffset int) uint32 {
	f := (((uint32(c.bgrBuffer[charColumn+4]) << 24) | (uint32(c.bgrBuffer[charColumn+5]) << 16) | (uint32(c.bgrBuffer[charColumn+6]) << 8) | (uint32(c.bgrBuffer[charColumn+7]))) << pixelOffset) | (uint32(c.bgrBuffer[charColumn+8]) >> (8 - pixelOffset))
	return f
}
