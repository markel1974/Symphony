package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

const (
	borderDisplayXFill    = 0x1ff
	borderDisplayXFillMax = borderDisplayXFill + 64 //DisplayXFill + 1
)

// CollisionsUnit encapsulate collision detection functionality between sprites and graphics within a VIC system.
// It includes buffers for handling sprite-sprite and sprite-graphics collision data as well as priorities.
// The struct also manages foreground masks and offsets for sprite-graphics collision computations.
type CollisionsUnit struct {
	*component.BaseComponent
	graphics             uint8   // graphics represents the current collision state with graphics as an 8-bit unsigned integer.
	spritesCollision     uint8   // spritesCollision represent the state and collision mask of all active sprites in the current frame.
	spritesPresence      []uint8 // spritesPresence stores the state of sprite presence at each pixel, used for collision detection within the frame.
	spritesPresenceEmpty []uint8 // spritesPresenceEmpty is the initialized buffer used to reset sprite collision states for each frame.
	graphicsBuffer       []uint8 // graphicsBuffer holds the collision state for background graphics as a buffer of unsigned 8-bit integers.
	graphicsBufferEmpty  []uint8 // graphicsBufferEmpty holds an initialized empty buffer for resetting or clearing graphics collision data.
	graphicsBufferOffset int     // graphicsBufferOffset tracks the current position in the graphics buffer for collision updates.
	sprBgrClx            uint8   // Sprite to background collision
	sprSprClx            uint8   // Sprite to sprite collision
	irqEmit              func(irq uint8)
}

// NewCollisions creates and returns a new CollisionsUnit instance, associated with the given VIC core.
func NewCollisions(parent references.IComponent, factory references.IComponentFactory, label string, instance int, irqEmit func(irq uint8), displayX int) *CollisionsUnit {
	c := &CollisionsUnit{
		BaseComponent:        component.NewBaseComponent(),
		irqEmit:              irqEmit,
		graphics:             0,
		spritesCollision:     0,
		spritesPresence:      make([]uint8, borderDisplayXFillMax), // Allocate the sprite buffer. Size is DisplayXFillMax (maximum X coordinate).
		spritesPresenceEmpty: make([]uint8, borderDisplayXFillMax), // Allocate and initialize the empty sprite buffer (all zeros).
		graphicsBuffer:       make([]uint8, borderDisplayXFill+1),  // Allocate the graphics buffer. Size is DisplayXFill+1. DisplayXFill seems to be 40
		graphicsBufferEmpty:  make([]uint8, displayX/8),
		graphicsBufferOffset: 0,
		sprSprClx:            0,
		sprBgrClx:            0,
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
	ret := c.sprSprClx
	c.sprSprClx = 0 // Read and clear
	return ret
}

// RetrieveSprite2Background reads and resets the sprite-to-background collision state, returning its value as an 8-bit unsigned integer.
func (c *CollisionsUnit) RetrieveSprite2Background() uint8 {
	// Sprite-background collision
	ret := c.sprBgrClx
	c.sprBgrClx = 0 // Read and clear
	return ret
}

// SetSprite updates the sprite-to-sprite collision state with the specified value.
func (c *CollisionsUnit) SetSprite(data uint8) {
	c.sprSprClx = data
}

// SetBackground sets the sprite-to-background collision state with the given value.
func (c *CollisionsUnit) SetBackground(data uint8) {
	c.sprBgrClx = data
}

// Prepare resets collision detection states and initializes sprite collision buffers for the next frame update.
// Called at the beginning of sprite drawing on each scanline.
func (c *CollisionsUnit) Prepare() {
	// Reset the sprite-to-sprite collision result.
	c.spritesCollision = 0
	// Reset the sprite-to-background collision result.
	c.graphics = 0
	// Reset the sprite buffer by copying the empty buffer.  This is *much* faster than iterating.
	copy(c.spritesPresence, c.spritesPresenceEmpty)
}

// SetGraphicsPresence sets a collision bit for graphics by performing a bitwise OR operation with the given bit.
// 'sBit' is a bitmask representing the sprite that collided with the background (1, 2, 4, 8, 16, 32, 64, 128).
func (c *CollisionsUnit) SetGraphicsPresence(sBit uint8) {
	// Set the corresponding bit in the 'graphics' collision result.
	c.graphics |= sBit
}

// SetSpritePresence checks and sets sprite collision at a specific index with a sprite bit and returns collision status.
// If a collision occurs, it updates the sprite collision state;
// otherwise, it updates the sprite buffer with the new bit.
func (c *CollisionsUnit) SetSpritePresence(index int, spiteBit uint8) bool {
	// Boundary check.
	if index >= borderDisplayXFillMax {
		return false
	}
	sBitPresence := c.spritesPresence[index]
	if sBitPresence == 0 {
		// mark this sprite as present at this pixel.
		c.spritesPresence[index] = spiteBit
		return false
	}
	// If any sprite is already present at this pixel...
	// Update the 'sprites' collision result with *both* the existing sprites *and* the new sprite.
	c.spritesCollision |= sBitPresence | spiteBit
	return true
}

// Commit triggers the collision application process using the stored sprite and graphics collision data.
// This method *actually writes* the collision results to the VIC-II's registers.
func (c *CollisionsUnit) Commit() {
	if c.sprSprClx != 0 {
		c.sprSprClx |= c.spritesCollision
	} else {
		c.sprSprClx |= c.spritesCollision
		c.irqEmit(irqSpriteToSpriteBit)
	}
	if c.sprBgrClx != 0 {
		c.sprBgrClx |= c.graphics
	} else {
		c.sprBgrClx |= c.graphics
		c.irqEmit(irqSpriteToGraphicBit)
	}
}

// IncrementGraphicsOffset increments the graphicsBufferOffset field
// to track the graphics buffer position during updates.
func (c *CollisionsUnit) IncrementGraphicsOffset() {
	// Increment the offset (in bytes).
	c.graphicsBufferOffset++
}

// ClearGraphics resets the graphics collision buffer and its offset to their initial empty states.
// Called at the *beginning* of each frame.
func (c *CollisionsUnit) ClearGraphics() {
	copy(c.graphicsBuffer, c.graphicsBufferEmpty)
	c.graphicsBufferOffset = 0
}

// UpdateGraphics updates the graphics collision buffer with provided foreground mask values a and b at the current offset.
// Called during *background* rendering.  'a' and 'b' represent the pixel data for two *consecutive* pixels.
func (c *CollisionsUnit) UpdateGraphics(a uint8, b uint8) {
	c.graphicsBuffer[c.graphicsBufferOffset] |= a   // Reset the graphics buffer by copying the empty buffer.
	c.graphicsBuffer[c.graphicsBufferOffset+1] |= b // Reset the offset.
}

// GetGraphicsL computes a 32-bit graphics mask from the graphics buffer starting at 'charColumn' with a 'pixelOffset' adjustment.
// The result is shifted and combined to align with the intended sub-pixel graphics position, returning a 32-bit uint representation.
func (c *CollisionsUnit) GetGraphicsL(charColumn int, pixelOffset int) uint32 {
	f := (((uint32(c.graphicsBuffer[charColumn]) << 24) | (uint32(c.graphicsBuffer[charColumn+1]) << 16) | (uint32(c.graphicsBuffer[charColumn+2]) << 8) | (uint32(c.graphicsBuffer[charColumn+3]))) << pixelOffset) | (uint32(c.graphicsBuffer[charColumn+4]) >> (8 - pixelOffset))
	return f
}

// GetGraphicsR calculates a 32-bit graphics mask from the graphics buffer starting at 'charColumn+4'.
// It shifts the mask by 'pixelOffset' bits to align with the intended sub-pixel graphics position.
// Returns a composite 32-bit unsigned integer representing the derived graphics mask.
func (c *CollisionsUnit) GetGraphicsR(charColumn int, pixelOffset int) uint32 {
	f := (((uint32(c.graphicsBuffer[charColumn+4]) << 24) | (uint32(c.graphicsBuffer[charColumn+5]) << 16) | (uint32(c.graphicsBuffer[charColumn+6]) << 8) | (uint32(c.graphicsBuffer[charColumn+7]))) << pixelOffset) | (uint32(c.graphicsBuffer[charColumn+8]) >> (8 - pixelOffset))
	return f
}
