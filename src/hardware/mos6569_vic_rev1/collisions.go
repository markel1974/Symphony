package mos6569

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/references"
)

// Collisions encapsulate collision detection functionality between sprites and graphics within a VIC system.
// It includes buffers for handling sprite-sprite and sprite-graphics collision data as well as priorities.
// The struct also manages foreground masks and offsets for sprite-graphics collision computations.
type Collisions struct {
	*component.BaseComponent
	core                 *VIC    // core references the VIC system used for handling collision detection and related graphics processing.
	graphics             uint8   // graphics represents the current collision state with graphics as an 8-bit unsigned integer.
	spritesCollision     uint8   // spritesCollision represent the state and collision mask of all active sprites in the current frame.
	spritesPresence      []uint8 // spritesPresence stores the state of sprite presence at each pixel, used for collision detection within the frame.
	spritesPresenceEmpty []uint8 // spritesPresenceEmpty is the initialized buffer used to reset sprite collision states for each frame.
	graphicsBuffer       []uint8 // graphicsBuffer holds the collision state for background graphics as a buffer of unsigned 8-bit integers.
	graphicsBufferEmpty  []uint8 // graphicsBufferEmpty holds an initialized empty buffer for resetting or clearing graphics collision data.
	graphicsBufferOffset int     // graphicsBufferOffset tracks the current position in the graphics buffer for collision updates.
}

// NewCollisions creates and returns a new Collisions instance, associated with the given VIC core.
func NewCollisions(parent references.IComponent, factory references.IComponentFactory, label string, instance int, core *VIC) *Collisions {
	c := &Collisions{
		BaseComponent:        component.NewBaseComponent(),
		core:                 core,
		graphics:             0,
		spritesCollision:     0,
		spritesPresence:      make([]uint8, DisplayXFillMax), // Allocate the sprite buffer. Size is DisplayXFillMax (maximum X coordinate).
		spritesPresenceEmpty: make([]uint8, DisplayXFillMax), // Allocate and initialize the empty sprite buffer (all zeros).
		graphicsBuffer:       make([]uint8, DisplayXFill+1),  // Allocate the graphics buffer. Size is DisplayXFill+1. DisplayXFill seems to be 40
		graphicsBufferEmpty:  make([]uint8, DisplayXDiv8),    // DisplayXDiv8 seems to be 52
		graphicsBufferOffset: 0,
	}
	c.BaseComponent.Register(factory, parent, "collisions", c, references.IdInternalComponent(label, instance, "Collisions"))
	return c
}

// Setup initializes the collision detection system, preparing necessary configurations for its operation.
func (c *Collisions) Setup() error {
	return nil
}

// Connect establishes necessary connections required by the Collisions instance and initializes its dependencies.
func (c *Collisions) Connect() error {
	return nil
}

// EmulationRequired determines if emulation is required for the current state of the collision system. Always returns false.
func (c *Collisions) EmulationRequired() bool {
	return false
}

// Emulate performs the collision emulation process for the current frame, updating internal state based on collisions detected.
func (c *Collisions) Emulate() {
}

// Internal checks if the collision detection system is operating in internal mode and returns a boolean value.
func (c *Collisions) Internal() bool {
	return true
}

// Reset clears all collision states and buffers, preparing the collision system for a new frame.
func (c *Collisions) Reset() {
}

// Prepare resets collision detection states and initializes sprite collision buffers for the next frame update.
// Called at the beginning of sprite drawing on each scanline.
func (c *Collisions) Prepare() {
	// Reset the sprite-to-sprite collision result.
	c.spritesCollision = uint8(0)
	// Reset the sprite-to-background collision result.
	c.graphics = uint8(0)
	// Reset the sprite buffer by copying the empty buffer.  This is *much* faster than iterating.
	copy(c.spritesPresence, c.spritesPresenceEmpty)
}

// SetGraphicsPresence sets a collision bit for graphics by performing a bitwise OR operation with the given bit.
// 'sBit' is a bitmask representing the sprite that collided with the background (1, 2, 4, 8, 16, 32, 64, 128).
func (c *Collisions) SetGraphicsPresence(sBit uint8) {
	// Set the corresponding bit in the 'graphics' collision result.
	c.graphics |= sBit
}

// SetSpritePresence checks and sets sprite collision at a specific index with a sprite bit and returns collision status.
// If a collision occurs, it updates the sprite collision state;
// otherwise, it updates the sprite buffer with the new bit.
func (c *Collisions) SetSpritePresence(index int, sBit uint8) bool {
	// Boundary check.
	if index >= DisplayXFillMax {
		return false
	}
	sBitPresence := c.spritesPresence[index]
	if sBitPresence == 0 {
		// mark this sprite as present at this pixel.
		c.spritesPresence[index] = sBit
		return false
	}

	// If any sprite is already present at this pixel...
	// Update the 'sprites' collision result with *both* the existing sprites *and* the new sprite.
	c.spritesCollision |= sBitPresence | sBit
	// Indicate that a collision occurred.
	return true
}

// Commit triggers the collision application process using the stored sprite and graphics collision data.
// This method *actually writes* the collision results to the VIC-II's registers.
func (c *Collisions) Commit() {
	// Call the CollisionApply method on the VIC core, passing the collision results.
	c.core.CollisionApply(c.spritesCollision, c.graphics)
}

// IncrementGraphicsOffset increments the graphicsBufferOffset field
// to track the graphics buffer position during updates.
func (c *Collisions) IncrementGraphicsOffset() {
	// Increment the offset (in bytes).
	c.graphicsBufferOffset++
}

// ClearGraphics resets the graphics collision buffer and its offset to their initial empty states.
// Called at the *beginning* of each frame.
func (c *Collisions) ClearGraphics() {
	copy(c.graphicsBuffer, c.graphicsBufferEmpty)
	c.graphicsBufferOffset = 0
}

// UpdateGraphics updates the graphics collision buffer with provided foreground mask values a and b at the current offset.
// Called during *background* rendering.  'a' and 'b' represent the pixel data for two *consecutive* pixels.
func (c *Collisions) UpdateGraphics(a uint8, b uint8) {
	// Reset the graphics buffer by copying the empty buffer.
	c.graphicsBuffer[c.graphicsBufferOffset] |= a
	// Reset the offset.
	c.graphicsBuffer[c.graphicsBufferOffset+1] |= b
}

// GetGraphicsL calculates a 32-bit graphics mask for collision detection.
// It assembles data from 5 consecutive bytes of the graphics buffer (from m to m+4)
// and shifts them by 's' bits to perfectly align with the sprite's sub-pixel position.
// This creates a "window" of graphics data that can be compared with the sprite.
func (c *Collisions) GetGraphicsL(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m]) << 24) | (uint32(c.graphicsBuffer[m+1]) << 16) | (uint32(c.graphicsBuffer[m+2]) << 8) | (uint32(c.graphicsBuffer[m+3]))) << s) | (uint32(c.graphicsBuffer[m+4]) >> (8 - s))
	return f
}

// GetGraphicsR calculates a 32-bit graphics mask from the graphics buffer using the specified offset and left shift.
// m is the major coordinate
// s is the shift
func (c *Collisions) GetGraphicsR(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m+4]) << 24) | (uint32(c.graphicsBuffer[m+5]) << 16) | (uint32(c.graphicsBuffer[m+6]) << 8) | (uint32(c.graphicsBuffer[m+7]))) << s) | (uint32(c.graphicsBuffer[m+8]) >> (8 - s))
	return f
}
