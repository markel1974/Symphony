package vic

var _sprEmptyCollBuf = make([]uint8, DisplayXFillMax)

type Sprites struct {
	core            *Core
	foreMask        *ForeMask
	displayBuffer   IDisplayBuffer
	collisionBuffer []uint8   // Buffer for sprite-sprite collisions and priorities
	dataPtr         []uint16  // Sprite data pointers
	data            [][]uint8 // Sprite data
	dmaFlags        uint8     // 8 flags: active DMA Sprite
	displayFlags    uint8     // 8 flags: active Display Sprite
	spriteFlags     uint8     // 8 flags: Sprite in this line
	dataCounter     []uint16  // Sprite counter data
	dataCounterBase []uint16  // Sprite base counter data
}

func NewSprites(core *Core, foreMask *ForeMask, db IDisplayBuffer) *Sprites {
	s := &Sprites{
		core:            core,
		foreMask:        foreMask,
		displayBuffer:   db,
		collisionBuffer: make([]uint8, DisplayXFillMax),
		dataPtr:         make([]uint16, SpriteNumber),
		data:            make([][]uint8, SpriteNumber),
		displayFlags:    0,
		dmaFlags:        0,
		dataCounter:     make([]uint16, SpriteNumber),
		dataCounterBase: make([]uint16, SpriteNumber),
	}
	for i := range s.data {
		s.data[i] = make([]uint8, 4)
	}
	for i := range s.dataCounter {
		s.dataCounter[i] = 63
	}
	return s
}

func (sp *Sprites) Setup() {
}

func (sp *Sprites) GetDMAFlag(bit uint8) uint8 {
	return sp.dmaFlags & bit
}

func (sp *Sprites) UpdateDisplayFlags() {
	sp.spriteFlags = sp.displayFlags

	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if (sp.displayFlags&mask) != 0 && (sp.dmaFlags&mask) == 0 {
			sp.displayFlags &= ^mask
		}
	}
}

func (sp *Sprites) FetchPtr(num int) {
	addr := sp.core.matrixBase | 0x03f8 | uint16(num)
	data := sp.core.ReadByte(addr)
	ptr := uint16(data) << 6
	sp.dataPtr[num] = ptr
}

func (sp *Sprites) Fetch(num int, bNum int) {
	if (sp.dmaFlags & (1 << num)) != 0 {
		ptr := sp.dataPtr[num]
		addr := (sp.dataCounter[num] & 0x3f) | ptr
		data := sp.core.ReadByte(addr)
		sp.data[num][bNum] = data
		sp.dataCounter[num]++
	} else if bNum == 1 {
		sp.core.ReadByte(0x3fff)
	}
}

func (sp *Sprites) UpdateDMA() {
	rasterY := sp.core.rasterY & 0xff
	for i, mask := 0, uint8(1); i < SpriteNumber; i, mask = i+1, mask<<1 {
		if (sp.core.me&mask) != 0 && rasterY == uint16(sp.core.mXy[i]) {
			sp.dmaFlags |= mask
			sp.dataCounterBase[i] = 0
			if (sp.core.mye & mask) != 0 {
				sp.core.sprExpY &= ^mask
			}
		}
	}
}

func (sp *Sprites) UpdateRasterYDisplayFlags() {
	rasterY := sp.core.rasterY & 0xff
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		sp.dataCounter[idx] = sp.dataCounterBase[idx]
		if (sp.dmaFlags&mask) != 0 && (rasterY == uint16(sp.core.mXy[idx])) {
			sp.displayFlags |= mask
		}
	}
}

func (sp *Sprites) UpdateDMACounterBase() {
	for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
		if (sp.core.sprExpY & mask) != 0 {
			sp.dataCounterBase[idx]++
		}
		if (sp.dataCounterBase[idx] & 0x3f) == 0x3f {
			sp.dmaFlags &= ^mask
		}
	}
}

func (sp *Sprites) UpdateCounterBase() {
	for idx := 0; idx < SpriteNumber; idx++ {
		if (sp.core.sprExpY & (1 << idx)) != 0 {
			sp.dataCounterBase[idx] += 2
		}
	}
}

func (sp *Sprites) Draw(lineStart int) {
	if sp.spriteFlags == 0 {
		return
	}
	sprColl := uint8(0)
	gfxColl := uint8(0)
	copy(sp.collisionBuffer, _sprEmptyCollBuf)
	for sNum, sBit := uint8(0), uint8(1); sNum < SpriteNumber; sNum, sBit = sNum+1, sBit<<1 {
		if (sp.spriteFlags & sBit) != 0 {
			sColor := sp.core.mXcColor[sNum]
			sData := (uint32(sp.data[sNum][0]) << 24) | (uint32(sp.data[sNum][1]) << 16) | (uint32(sp.data[sNum][2]) << 8)
			sOffset := int(sp.core.mXx[sNum]) + SpriteNumber
			lineOffset := lineStart + sOffset
			m := sOffset / SpriteNumber
			s := sOffset & 7

			expanded := sp.core.mxe&sBit != 0
			multiColor := sp.core.mmc&sBit != 0
			if expanded {
				if multiColor {
					sp.drawExpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit, &gfxColl, &sprColl)
				} else {
					sp.drawExpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit, &gfxColl, &sprColl)
				}
			} else {
				if multiColor {
					sp.drawUnexpandedMulticolor(lineOffset, sColor, sData, sOffset, m, s, sBit, &gfxColl, &sprColl)
				} else {
					sp.drawUnexpandedStandard(lineOffset, sColor, sData, sOffset, m, s, sBit, &gfxColl, &sprColl)
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
		if (sp.core.irqMask & 0x04) != 0 {
			sp.core.irqFlag |= 0x80
			sp.core.signalIRQTrigger.Emit(intrVicId)
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
			sp.core.signalIRQTrigger.Emit(intrVicId)
		}
	}
}

func (sp *Sprites) drawExpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	foreMaskL := sp.foreMask.GetL(m, s)
	foreMaskR := sp.foreMask.GetR(m, s)
	// Expand sprite data
	sDataL := uint32(_multiExpTable[sData>>24&0xff])<<16 | uint32(_multiExpTable[sData>>16&0xff])
	sDataR := uint32(_multiExpTable[sData>>8&0xff]) << 16
	// Convert sprite pixels to bitPlanes
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
	idx := 0
	for ; idx < 32; idx, plane0L, plane1L = idx+1, plane0L<<1, plane1L<<1 {
		selectedColor := uint8(0)
		if (plane1L & 0x80000000) != 0 {
			if (plane0L & 0x80000000) != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0L & 0x80000000) != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawPixel(lineOffset, sOffset, idx, sBit, selectedColor, sprColl)
	}
	for ; idx < 48; idx, plane0R, plane1R = idx+1, plane0R<<1, plane1R<<1 {
		selectedColor := uint8(0)
		if (plane1R & 0x80000000) != 0 {
			if (plane0R & 0x80000000) != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0R & 0x80000000) != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawPixel(lineOffset, sOffset, idx, sBit, selectedColor, sprColl)
	}
}

func (sp *Sprites) drawExpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	foreMaskL := sp.foreMask.GetL(m, s)
	foreMaskR := sp.foreMask.GetR(m, s)
	sDataL := uint32(_expTable[sData>>24&0xff])<<16 | uint32(_expTable[sData>>16&0xff])
	sDataR := uint32(_expTable[sData>>8&0xff]) << 16
	if (foreMaskL&sDataL) != 0 || (foreMaskR&sDataR) != 0 {
		*gfxColl |= sBit
		if (sp.core.mdp & sBit) != 0 {
			sDataL &= ^foreMaskL
			sDataR &= ^foreMaskR
		}
	}
	var idx = 0
	for ; idx < 32; idx, sDataL = idx+1, sDataL<<1 {
		if (sDataL & 0x80000000) != 0 {
			sp.drawPixel(lineOffset, sOffset, idx, sBit, sColor, sprColl)
		}
	}
	for ; idx < 48; idx, sDataR = idx+1, sDataR<<1 {
		if (sDataR & 0x80000000) != 0 {
			sp.drawPixel(lineOffset, sOffset, idx, sBit, sColor, sprColl)
		}
	}
}

func (sp *Sprites) drawUnexpandedMulticolor(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	foreMask := sp.foreMask.GetL(m, s)
	// Convert sprite pixels to bitPlanes
	p0 := sData & 0x55555555
	p1 := sData & 0xaaaaaaaa
	plane0 := p0 | (p0 << 1)
	plane1 := p1 | (p1 >> 1)
	// Check graphics collision
	if (foreMask & (plane0 | plane1)) != 0 {
		*gfxColl |= sBit
		if (sp.core.mdp & sBit) != 0 {
			// Mask sprite if in background
			plane0 &= ^foreMask
			plane1 &= ^foreMask
		}
	}
	for idx := 0; idx < 24; idx, plane0, plane1 = idx+1, plane0<<1, plane1<<1 {
		var selectedColor uint8
		if (plane1 & 0x80000000) != 0 {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = sp.core.mm1Color
			} else {
				selectedColor = sColor
			}
		} else {
			if (plane0 & 0x80000000) != 0 {
				selectedColor = sp.core.mm0Color
			} else {
				continue
			}
		}
		sp.drawPixel(lineOffset, sOffset, idx, sBit, selectedColor, sprColl)
	}
}

func (sp *Sprites) drawUnexpandedStandard(lineOffset int, sColor uint8, sData uint32, sOffset int, m int, s int, sBit uint8, gfxColl *uint8, sprColl *uint8) {
	foreMask := sp.foreMask.GetL(m, s)
	// Check graphics collision
	if (foreMask & sData) != 0 {
		*gfxColl |= sBit
		if sp.core.mdp&sBit != 0 {
			// Mask sprite if in background
			sData &= ^foreMask
		}
	}
	for idx := 0; idx < 24; idx, sData = idx+1, sData<<1 {
		if (sData & 0x80000000) != 0 {
			sp.drawPixel(lineOffset, sOffset, idx, sBit, sColor, sprColl)
		}
	}
}

func (sp *Sprites) drawPixel(lineOffset int, sOffset int, idx int, sBit uint8, selColor uint8, sprColl *uint8) {
	collision := false
	if collIdx := sOffset + idx; collIdx < DisplayXFillMax {
		if sp.collisionBuffer[collIdx] != 0 {
			*sprColl |= sp.collisionBuffer[collIdx] | sBit
			collision = true
		} else {
			sp.collisionBuffer[collIdx] = sBit
		}
	}
	if !collision {
		sp.displayBuffer.Set(lineOffset+idx, selColor)
	}
}
