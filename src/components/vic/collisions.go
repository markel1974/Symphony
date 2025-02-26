package mos6569

// Collisions encapsulate collision detection functionality between sprites and graphics within a VIC system.
// It includes buffers for handling sprite-sprite and sprite-graphics collision data as well as priorities.
// The struct also manages foreground masks and offsets for sprite-graphics collision computations.
type Collisions struct {
	core                 *VIC
	graphics             uint8
	sprites              uint8
	spritesBuffer        []uint8 // Buffer for sprite-sprite collisions and priorities
	emptySpritesBuffer   []uint8
	graphicsBuffer       []uint8 // Foreground mask for sprite-graphics collisions and priorities
	graphicsBufferEmpty  []uint8
	graphicsBufferOffset int
}

// NewCollisions creates a new instance of the Collisions struct, initializing buffers and setting default values.
func NewCollisions(core *VIC) *Collisions {
	return &Collisions{
		core:                 core,
		graphics:             0,
		sprites:              0,
		spritesBuffer:        make([]uint8, DisplayXFillMax),
		emptySpritesBuffer:   make([]uint8, DisplayXFillMax),
		graphicsBufferOffset: 0,
		graphicsBuffer:       make([]uint8, DisplayXFill+1),
		graphicsBufferEmpty:  make([]uint8, DisplayXDiv8),
		//graphicsBufferEmpty: make([]uint8, DisplayXFill+1),
	}
}

// Prepare resets collision detection states and initializes sprite collision buffers for the next frame update.
func (c *Collisions) Prepare() {
	c.sprites = uint8(0)
	c.graphics = uint8(0)
	copy(c.spritesBuffer, c.emptySpritesBuffer)
}

// SetGraphicsCollision sets a collision bit for graphics by performing a bitwise OR operation with the given bit.
func (c *Collisions) SetGraphicsCollision(sBit uint8) {
	c.graphics |= sBit
}

// SetSpriteCollision checks and sets sprite collision at a specific index with a sprite bit and returns collision status.
// If a collision occurs, it updates the sprite collision state;
// otherwise, it updates the sprite buffer with the new bit.
func (c *Collisions) SetSpriteCollision(collIdx int, sBit uint8) bool {
	collision := false
	if collIdx < DisplayXFillMax {
		if c.spritesBuffer[collIdx] != 0 {
			c.sprites |= c.spritesBuffer[collIdx] | sBit
			collision = true
		} else {
			c.spritesBuffer[collIdx] = sBit
		}
	}
	return collision
}

// Detect triggers the collision application process using the stored sprite and graphics collision data.
func (c *Collisions) Detect() {
	c.core.CollisionApply(c.sprites, c.graphics)
}

// IncrementGraphicsOffset increments the graphicsBufferOffset field
// to track the graphics buffer position during updates.
func (c *Collisions) IncrementGraphicsOffset() {
	c.graphicsBufferOffset++
}

// ClearGraphics resets the graphics collision buffer and its offset to their initial empty states.
func (c *Collisions) ClearGraphics() {
	copy(c.graphicsBuffer, c.graphicsBufferEmpty)
	c.graphicsBufferOffset = 0
}

// UpdateGraphics updates the graphics collision buffer with provided foreground mask values a and b at the current offset.
func (c *Collisions) UpdateGraphics(a uint8, b uint8) {
	c.graphicsBuffer[c.graphicsBufferOffset] |= a
	c.graphicsBuffer[c.graphicsBufferOffset+1] |= b
}

// GetGraphicsL computes a 32-bit mask
// by combining and shifting values from the graphicsBuffer at specified offsets and shifts.
func (c *Collisions) GetGraphicsL(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m]) << 24) | (uint32(c.graphicsBuffer[m+1]) << 16) | (uint32(c.graphicsBuffer[m+2]) << 8) | (uint32(c.graphicsBuffer[m+3]))) << s) | (uint32(c.graphicsBuffer[m+4]) >> (8 - s))
	return f
}

// GetGraphicsR calculates a 32-bit graphics mask from the graphics buffer using the specified offset and left shift.
func (c *Collisions) GetGraphicsR(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m+4]) << 24) | (uint32(c.graphicsBuffer[m+5]) << 16) | (uint32(c.graphicsBuffer[m+6]) << 8) | (uint32(c.graphicsBuffer[m+7]))) << s) | (uint32(c.graphicsBuffer[m+8]) >> (8 - s))
	return f
}
