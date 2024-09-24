package mos6569

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

func (c *Collisions) Prepare() {
	c.sprites = uint8(0)
	c.graphics = uint8(0)
	copy(c.spritesBuffer, c.emptySpritesBuffer)
}

func (c *Collisions) SetGraphicsCollision(sBit uint8) {
	c.graphics |= sBit
}

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

func (c *Collisions) Detect() {
	c.core.CollisionApply(c.sprites, c.graphics)
}

func (c *Collisions) IncrementGraphicsOffset() {
	c.graphicsBufferOffset++
}

func (c *Collisions) ClearGraphics() {
	copy(c.graphicsBuffer, c.graphicsBufferEmpty)
	c.graphicsBufferOffset = 0
}

func (c *Collisions) UpdateGraphics(a uint8, b uint8) {
	c.graphicsBuffer[c.graphicsBufferOffset] |= a
	c.graphicsBuffer[c.graphicsBufferOffset+1] |= b
}

func (c *Collisions) GetGraphicsL(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m]) << 24) | (uint32(c.graphicsBuffer[m+1]) << 16) | (uint32(c.graphicsBuffer[m+2]) << 8) | (uint32(c.graphicsBuffer[m+3]))) << s) | (uint32(c.graphicsBuffer[m+4]) >> (8 - s))
	return f
}

func (c *Collisions) GetGraphicsR(m int, s int) uint32 {
	f := (((uint32(c.graphicsBuffer[m+4]) << 24) | (uint32(c.graphicsBuffer[m+5]) << 16) | (uint32(c.graphicsBuffer[m+6]) << 8) | (uint32(c.graphicsBuffer[m+7]))) << s) | (uint32(c.graphicsBuffer[m+8]) >> (8 - s))
	return f
}
