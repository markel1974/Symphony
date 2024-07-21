package vic

var _sprEmptyCollBuf = make([]uint8, DisplayXFillMax)

type Sprites struct {
	core       *Core
	foreMask   *ForeMask
	intr       IInterrupts
	db         IDisplayBuffer
	sprCollBuf []uint8   // Buffer for sprite-sprite collisions and priorities
	sprPtr     []uint16  // Sprite data pointers
	sprData    [][]uint8 // Sprite data read
	//sprDrawData [][]uint8 // Sprite data for drawing
	sprFlags    uint8    // 8 flags: Draw sprite in this line
	sprDrawData []uint32 // Sprite data for drawing
}

func NewSprites(core *Core, foreMask *ForeMask, db IDisplayBuffer) *Sprites {
	s := &Sprites{
		core:       core,
		foreMask:   foreMask,
		db:         db,
		intr:       nil,
		sprCollBuf: make([]uint8, DisplayXFillMax),
		sprPtr:     make([]uint16, SpriteNumber),
		sprData:    make([][]uint8, SpriteNumber),
		//sprDrawData: make([][]uint8, SpriteNumber),
		sprDrawData: make([]uint32, SpriteNumber),
	}
	for i := range s.sprData {
		s.sprData[i] = make([]uint8, 4)
	}
	//for i := range s.sprDrawData {
	//	s.sprDrawData[i] = make([]uint8, 4)
	//}
	return s
}

func (sp *Sprites) Setup(intr IInterrupts) {
	sp.intr = intr
}

func (sp *Sprites) GetSpritePtr(num int) uint16 {
	return sp.sprPtr[num]
}

func (sp *Sprites) SetSpritePtr(num int, data uint16) {
	sp.sprPtr[num] = data
}

func (sp *Sprites) SetSpriteData(num int, byteNum int, data uint8) {
	sp.sprData[num][byteNum] = data
}

func (sp *Sprites) SetSpriteFlags(spriteFlags uint8) {
	sp.sprFlags = spriteFlags
	if sp.sprFlags != 0 {
		for sNum := 0; sNum < len(sp.sprData); sNum++ {
			sp.sprDrawData[sNum] = (uint32(sp.sprData[sNum][0]) << 24) | (uint32(sp.sprData[sNum][1]) << 16) | (uint32(sp.sprData[sNum][2]) << 8)
		}
	}
}

/*
func (sp * Sprites) Rebuild() {
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		vic.sprDataCounter[idx] = vic.sprDataCounterBase[idx]
		if (vic.sprDMAFlags&mask) != 0 && (rasterY == uint16(vic.core.mXy[idx])) {
			vic.sprDisplayFlags |= mask
		}
	}
}
*/

func (sp *Sprites) Draw(lineStart int) {
	if sp.sprFlags == 0 {
		return
	}
	sprColl := uint8(0)
	gfxColl := uint8(0)
	copy(sp.sprCollBuf, _sprEmptyCollBuf)
	for sNum, sBit := uint8(0), uint8(1); sNum < 8; sNum, sBit = sNum+1, sBit<<1 {
		if sp.sprFlags&sBit != 0 {
			expanded := sp.core.mxe&sBit != 0
			multiColor := sp.core.mmc&sBit != 0
			if expanded {
				if multiColor {
					sp.drawSpriteExpandedMulticolor(lineStart, sNum, sBit, &gfxColl, &sprColl)
				} else {
					sp.drawSpriteExpandedStandard(lineStart, sNum, sBit, &gfxColl, &sprColl)
				}
			} else {
				if multiColor {
					sp.drawSpriteUnexpandedMulticolor(lineStart, sNum, sBit, &gfxColl, &sprColl)
				} else {
					sp.drawSpriteUnexpandedStandard(lineStart, sNum, sBit, &gfxColl, &sprColl)
				}
			}
		}
	}
	// sprite-sprite collisions
	if sp.core.sprClx != 0 {
		sp.core.sprClx |= sprColl
	} else {
		sp.core.sprClx |= sprColl
		sp.core.irqFlag |= 0x04
		if sp.core.irqMask&0x04 != 0 {
			sp.core.irqFlag |= 0x80
			sp.intr.TriggerVICIRQ()
		}
	}
	// sprite-background collisions
	if sp.core.sprClxBgr != 0 {
		sp.core.sprClxBgr |= gfxColl
	} else {
		sp.core.sprClxBgr |= gfxColl
		sp.core.irqFlag |= 0x02
		if sp.core.irqMask&0x02 != 0 {
			sp.core.irqFlag |= 0x80
			sp.intr.TriggerVICIRQ()
		}
	}
}

func (sp *Sprites) drawSpriteExpandedMulticolor(lineStart int, sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(sp.core.mXx[sNum]) + 8
	displayPtr := lineStart + q
	color := sp.core.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMaskL := sp.foreMask.GetL(m, s)
	foreMaskR := sp.foreMask.GetR(m, s)
	sData := sp.sprDrawData[sNum] //(uint32(sp.sprDrawData[sNum][0]) << 24) | (uint32(sp.sprDrawData[sNum][1]) << 16) | (uint32(sp.sprDrawData[sNum][2]) << 8)

	// Expand sprite data
	sDataL := uint32(_multiExpTable[sData>>24&0xff])<<16 | uint32(_multiExpTable[sData>>16&0xff])
	sDataR := uint32(_multiExpTable[sData>>8&0xff]) << 16
	// Convert sprite chunky pixels to bitPlanes
	plane0L := (sDataL & 0x55555555) | (sDataL&0x55555555)<<1
	plane1L := (sDataL & 0xaaaaaaaa) | (sDataL&0xaaaaaaaa)>>1
	plane0R := (sDataR & 0x55555555) | (sDataR&0x55555555)<<1
	plane1R := (sDataR & 0xaaaaaaaa) | (sDataR&0xaaaaaaaa)>>1
	// Collision with graphics?
	if (foreMaskL&(plane0L|plane1L)) != 0 || (foreMaskR&(plane0R|plane1R)) != 0 {
		*gfxColl |= sBit
		if sp.core.mdp&sBit != 0 {
			// Mask sprite if in background
			plane0L &= ^foreMaskL
			plane1L &= ^foreMaskL
			plane0R &= ^foreMaskR
			plane1R &= ^foreMaskR
		}
	}
	// Paint sprite
	idx := 0
	for ; idx < 32; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		selectedColor := uint8(0)
		if plane1L&0x80000000 != 0 {
			if plane0L&0x80000000 != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if plane0L&0x80000000 != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawSpritePixel(displayPtr, q, idx, sBit, selectedColor, sprColl)
	}
	for ; idx < 48; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		selectedColor := uint8(0)
		if plane1R&0x80000000 != 0 {
			if plane0R&0x80000000 != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if plane0R&0x80000000 != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawSpritePixel(displayPtr, q, idx, sBit, selectedColor, sprColl)
	}
}

func (sp *Sprites) drawSpriteExpandedStandard(lineStart int, sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(sp.core.mXx[sNum]) + 8
	displayPtr := lineStart + q
	color := sp.core.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := sp.foreMask.GetA(m, s)
	sData := sp.sprDrawData[sNum] //(uint32(sp.sprDrawData[sNum][0]) << 24) | (uint32(sp.sprDrawData[sNum][1]) << 16) | (uint32(sp.sprDrawData[sNum][2]) << 8)
	// Check graphics collision
	if (foreMask & sData) != 0 {
		*gfxColl |= sBit
		if sp.core.mdp&sBit != 0 {
			// Mask sprite if in background
			sData &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			sp.drawSpritePixel(displayPtr, q, idx, sBit, color, sprColl)
		}
	}
}

func (sp *Sprites) drawSpriteUnexpandedMulticolor(lineStart int, sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(sp.core.mXx[sNum]) + 8
	displayPtr := lineStart + q
	color := sp.core.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := sp.foreMask.GetL(m, s) //sp.foreMask.GetB(m, s)
	sData := sp.sprDrawData[sNum]      //(uint32(sp.sprDrawData[sNum][0]) << 24) | (uint32(sp.sprDrawData[sNum][1]) << 16) | (uint32(sp.sprDrawData[sNum][2]) << 8)
	// Convert sprite pixels to bitPlanes
	plane0 := (sData & 0x55555555) | (sData&0x55555555)<<1
	plane1 := (sData & 0xaaaaaaaa) | (sData&0xaaaaaaaa)>>1
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		*gfxColl |= sBit
		if sp.core.mdp&sBit != 0 {
			// Mask sprite if in background
			plane0 &= ^foreMask
			plane1 &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		var selectedColor uint8
		if (plane1 & 0x80000000) != 0 {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = color
			}
		} else {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawSpritePixel(displayPtr, q, idx, sBit, selectedColor, sprColl)
	}
}

func (sp *Sprites) drawSpriteUnexpandedStandard(lineStart int, sNum uint8, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	q := int(sp.core.mXx[sNum]) + 8
	displayPtr := lineStart + q
	color := sp.core.mXcColor[sNum]
	m := q / 8
	s := q & 7
	foreMask := sp.foreMask.GetL(m, s) //sp.foreMask.GetC(m, s)
	sData := sp.sprDrawData[sNum]      //(uint32(sp.sprDrawData[sNum][0]) << 24) | (uint32(sp.sprDrawData[sNum][1]) << 16) | (uint32(sp.sprDrawData[sNum][2]) << 8)
	// Check graphics collision
	if (foreMask & sData) != 0 {
		*gfxColl |= sBit
		if sp.core.mdp&sBit != 0 {
			// Mask sprite if in background
			sData &= ^foreMask
		}
	}
	// Paint sprite
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			sp.drawSpritePixel(displayPtr, q, idx, sBit, color, sprColl)
		}
	}
}

func (sp *Sprites) drawSpritePixel(displayPtr int, q int, idx int, sBit uint8, selColor uint8, sprColl *uint8) {
	if collIdx := q + idx; collIdx < DisplayXFillMax {
		// Check graphics collision
		if (sp.sprCollBuf[collIdx]) != 0 {
			// Collision with sprite?
			*sprColl |= sp.sprCollBuf[collIdx] | sBit
		} else {
			// Draw pixel if no collision
			sp.db.Set(displayPtr+idx, selColor)
			sp.sprCollBuf[collIdx] = sBit
		}
	}
}
