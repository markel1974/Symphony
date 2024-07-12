package vic

type Core struct {
	mXx            []uint16 // VIC registers [m0x - m1x - m2x - m3x - m4x - m5x - m6x - m7x]
	mXy            []uint8  // VIC registers [m0y - m1y - m2y - m3y - m4y - m5y - m6y - m7y]
	mx8            uint8    // VIC register
	cr1            uint8    // VIC register
	cr2            uint8    // VIC register
	lpx            uint8    // VIC register
	lpy            uint8    // VIC register
	me             uint8    // VIC register
	mxe            uint8    // VIC register
	mye            uint8    // VIC register
	mdp            uint8    // VIC register
	mmc            uint8    // VIC register
	ec             uint8    // VIC register
	b0c            uint8    // VIC register
	b1c            uint8    // VIC register
	b2c            uint8    // VIC register
	b3c            uint8    // VIC register
	mm0            uint8    // VIC register
	mm1            uint8    // VIC register
	mXc            []uint8  // VIC registers [m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c]
	ecColor        uint8    // Index ec Color Mapping
	b0cColor       uint8    // Index b0c Color Mapping
	b1cColor       uint8    // Index b1c Color Mapping
	b2cColor       uint8    // Index b2c Color Mapping
	b3cColor       uint8    // Index b3c Color Mapping
	mm0Color       uint8    // Index mm0 Color Mapping
	mm1Color       uint8    // Index mm1 Color Mapping
	mXcColor       []uint8  // Indices for m0c - m1c - m2c - m3c - m4c - m5c - m6c - m7c Color Mapping
	vaBase         uint8    // vaBase
	ciaVaBase      uint16   // CIA VA14/15 video base
	matrixBase     uint16   // Video matrix base
	charBase       uint16   // Character generator base
	bitmapBase     uint16   // Bitmap base
	displayIdx     int      // Index of current display mode
	xScroll        uint16   // X scroll value
	yScroll        uint16   // Y scroll valuej
	irqFlag        uint8    //
	irqMask        uint8    //
	irqRaster      uint16   // Interrupt raster line
	sprExpY        uint8    // 8 sprite y expansion FlipFlops
	sprClxBgr      uint8    // Sprite to background collision
	sprClx         uint8    // Sprite to sprite collision
	foreMaskOffset int      // Offset from foreMaskBuf
	foreMaskBuf    []uint8  // Foreground mask for sprite-graphics collisions and priorities
	colors         []uint8  // Indices of the 16 colors (16 times mirrored to avoid "& 0x0f")
	displayBuffer  []uint8  //
}

func NewCore() *Core {
	colors := make([]uint8, 256)
	for i := range colors {
		colors[i] = (uint8)(i & 0x0f)
	}
	c := &Core{
		mXx:            make([]uint16, SpriteNumber),
		mXy:            make([]uint8, SpriteNumber),
		mx8:            0,
		cr1:            0,
		cr2:            0,
		lpx:            0,
		lpy:            0,
		me:             0,
		mxe:            0,
		mye:            0,
		mdp:            0,
		mmc:            0,
		ec:             0,
		b0c:            0,
		b1c:            0,
		b2c:            0,
		b3c:            0,
		mm0:            0,
		mm1:            0,
		mXc:            make([]uint8, SpriteNumber),
		mXcColor:       make([]uint8, SpriteNumber),
		matrixBase:     0,
		charBase:       0,
		bitmapBase:     0,
		vaBase:         0,
		ciaVaBase:      0,
		displayIdx:     0,
		xScroll:        0,
		yScroll:        0,
		irqRaster:      0,
		irqFlag:        0,
		irqMask:        0,
		sprExpY:        0,
		sprClx:         0,
		sprClxBgr:      0,
		foreMaskOffset: 0,
		ecColor:        colors[0], // Preset colors to black
		b0cColor:       colors[0], // Preset colors to black
		b1cColor:       colors[0], // Preset colors to black
		b2cColor:       colors[0], // Preset colors to black
		b3cColor:       colors[0], // Preset colors to black
		mm0Color:       colors[0], // Preset colors to black
		mm1Color:       colors[0], // Preset colors to black
		colors:         colors,
		displayBuffer:  make([]uint8, DisplaySize),
		foreMaskBuf:    make([]uint8, DisplayXFill+1),
	}
	// Preset colors to black
	for i := range c.mXcColor {
		c.mXcColor[i] = c.colors[0]
	}
	return c
}
