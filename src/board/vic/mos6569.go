package vic

import (
	"github.com/markel1974/c64emu/src/board/quartz"
	"github.com/markel1974/c64emu/src/preferences"
	"github.com/markel1974/c64emu/src/signals"
)

//https://dustlayer.com/c64-architecture

var _emptyForeMaskBuffer = make([]uint8, DisplayXDiv8)

type MOS6569 struct {
	core               *Core
	prefs              *preferences.Prefs
	sprites            *Sprites
	graphics           *Graphics
	foreMask           *ForeMask
	db                 IDisplayBuffer
	cycle              int      // Cycle
	lineStart          int      // Offset from current line in bitmap buffer
	lineOffset         int      // Offset from chunky bitmap buffer
	matrixLine         []uint8  // Buffer for video line, read in Bad Lines
	colorLine          []uint8  // Buffer for color line, read in Bad Lines
	rowCounter         uint16   // Row counter
	videoCounter       uint16   // Video counter
	videoCounterBase   uint16   // Video counter base
	borderULOn         bool     // Upper/lower border on
	displayOn          bool     // Display state
	drawThisLine       bool     // This line is drawn
	sprDMAFlags        uint8    // 8 flags: Sprite DMA active
	sprDisplayFlags    uint8    // 8 flags: Sprite display active
	sprDataCounter     []uint16 // Sprite data counters
	sprDataCounterBase []uint16 // Sprite data counter bases
	refreshCounter     uint8    // Refresh counter
	mlIndex            int      // Index in matrix/colorLine[]
	vBlanking          bool     // VBlank in next cycle
}

func NewMOS6569(db IDisplayBuffer) *MOS6569 {
	core := NewCore()
	foreMask := NewForeMask()
	//db := NewDisplayBuffer()
	vic := &MOS6569{
		core:               core,
		foreMask:           foreMask,
		db:                 db,
		graphics:           NewGraphics(core, foreMask, db),
		sprites:            NewSprites(core, foreMask, db),
		lineOffset:         0,
		rowCounter:         7,
		videoCounter:       0,
		videoCounterBase:   0,
		mlIndex:            0,
		cycle:              1,
		displayOn:          false,
		vBlanking:          false,
		drawThisLine:       false,
		sprDisplayFlags:    0,
		sprDMAFlags:        0,
		matrixLine:         make([]uint8, 40),
		colorLine:          make([]uint8, 40),
		sprDataCounter:     make([]uint16, SpriteNumber),
		sprDataCounterBase: make([]uint16, SpriteNumber),
	}
	for i := range vic.sprDataCounter {
		vic.sprDataCounter[i] = 63
	}
	return vic
}

func (vic *MOS6569) Setup(quartz *quartz.Quartz, intr IInterrupts, banks IBanks, prefs *preferences.Prefs) {
	//vic.board = board
	vic.prefs = prefs
	vic.core.Setup(quartz, intr, banks)
	vic.graphics.Setup()
	vic.sprites.Setup(intr)
}

//func (vic *MOS6569) GetDisplayBuffer() []uint8 {
//	return vic.db.Get()
//}

func (vic *MOS6569) Reset() {
	vic.core.ready = false
}

func (vic *MOS6569) GetLastByte() uint8 {
	return vic.core.lastByte
}

func (vic *MOS6569) GetBALow() bool {
	return vic.core.baLow
}

func (vic *MOS6569) GetReadySignal() *signals.Signal {
	return vic.core.readySignal
}

//func (vic *MOS6569) GetBALowSignal() *signals.Signal1[bool] {
//	return vic.core.baLowSignal
//}

func (vic *MOS6569) NewPrefs(_ *preferences.Prefs) {
	//vic.skipFrames = prefs.SkipFrames()
}

func (vic *MOS6569) ReadRegister(addr uint16) uint8 {
	return vic.core.ReadRegister(addr)
}

func (vic *MOS6569) WriteRegister(addr uint16, data uint8) {
	vic.core.WriteRegister(addr, data)
}

func (vic *MOS6569) ChangedVA(va uint8) {
	vic.core.ChangedVA(va)
}

func (vic *MOS6569) LightPenTrigger() {
	vic.core.LightPenTrigger()
}

func (vic *MOS6569) Emulate() (bool, bool) {
	vBlank := false
	lastCycle := false

	switch vic.cycle {
	case 1:
		// Fetch sprite pointer 3, increment raster counter, trigger raster IRQ,
		// test for Bad Line, reset BA if sprites 3 and 4 off, read data of sprite 3
		if vic.core.rasterY == RasterYMax {
			// Trigger VBlank in cycle 2
			vic.vBlanking = true
		} else {
			// Increment raster counter
			vic.core.rasterY++
			// Trigger raster IRQ if IRQ line reached
			if vic.core.rasterY == vic.core.irqRaster {
				vic.core.rasterIrq()
			}
			// In line $30, the DEN bit controls if Bad Lines can occur
			if vic.core.rasterY == 0x30 {
				vic.core.badLinesEnabled = vic.core.cr1&0x10 != 0
			}
			// Bad Line condition?
			vic.core.badLineUpdate()
			// Don't draw all lines, hide some at the top and bottom
			vic.drawThisLine = (vic.core.rasterY >= FirstDisplayedLine) && (vic.core.rasterY <= LastDisplayedLine)
		}
		// First sample of border state
		vic.graphics.SetBorderOnSample(0)
		vic.fetchSpriteDataPtr(3)
		vic.fetchSpriteData(3, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x18) == 0 {
			vic.core.ClearBALow()
		}

	case 2:
		// Set BA for sprite 5, read data of sprite 3
		if vic.vBlanking {
			vBlank = true
			// Vertical blank, reset counters
			vic.videoCounterBase = 0
			vic.refreshCounter = 0xff
			vic.vBlanking = false
			vic.lineStart = 0
			vic.core.ResetCounters()
		}
		// Our output goes here
		vic.lineOffset = vic.lineStart
		vic.foreMask.Clear()

		vic.fetchSpriteData(3, 1)
		vic.fetchSpriteData(3, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x20) != 0 {
			vic.core.SetBALow()
		}
	case 3:
		// Fetch sprite pointer 4, reset BA is sprite 4 and 5 off
		vic.fetchSpriteDataPtr(4)
		vic.fetchSpriteData(4, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x30) == 0 {
			vic.core.ClearBALow()
		}
	case 4:
		// Set BA for sprite 6, read data of sprite 4
		vic.fetchSpriteData(4, 1)
		vic.fetchSpriteData(4, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x40) != 0 {
			vic.core.SetBALow()
		}
	case 5:
		// Fetch sprite pointer 5, reset BA if sprite 5 and 6 off
		vic.fetchSpriteDataPtr(5)
		vic.fetchSpriteData(5, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x60) == 0 {
			vic.core.ClearBALow()
		}
	case 6:
		// Set BA for sprite 7, read data of sprite 5
		vic.fetchSpriteData(5, 1)
		vic.fetchSpriteData(5, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x80) != 0 {
			vic.core.SetBALow()
		}
	case 7:
		// Fetch sprite pointer 6, reset BA if sprite 6 and 7 off
		vic.fetchSpriteDataPtr(6)
		vic.fetchSpriteData(6, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0xc0) == 0 {
			vic.core.ClearBALow()
		}
	case 8:
		// Read data of sprite 6
		vic.fetchSpriteData(6, 1)
		vic.fetchSpriteData(6, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
	case 9:
		// Fetch sprite pointer 7, reset BA if sprite 7 off
		vic.fetchSpriteDataPtr(7)
		vic.fetchSpriteData(7, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x80) == 0 {
			vic.core.ClearBALow()
		}
	case 10:
		// Read data of sprite 7
		vic.fetchSpriteData(7, 1)
		vic.fetchSpriteData(7, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
	case 11:
		// Refresh, reset BA
		vic.refreshAccess()
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		vic.core.ClearBALow()

	case 12:
		// Refresh, turn on matrix access if Bad Line
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}

	case 13:
		// Refresh, turn on matrix access if Bad Line, reset rasterX, graphics display starts here
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		vic.core.rasterX = 0xfffc

	case 14:
		// Refresh, videoCounter -> videoCounterBase, turn on matrix access and reset RC if Bad Line
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access and reset RC if Bad Line
		if vic.core.isBadLine {
			vic.rowCounter = 0
			vic.displayOn = true
			vic.core.SetBALow()
		}
		vic.videoCounter = vic.videoCounterBase

	case 15:
		// Refresh and matrix access, increment sprDataCounterBase by 2 if y expansion FlipFlop is set
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.refreshAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		for idx := 0; idx < 8; idx++ {
			if (vic.core.sprExpY & (1 << idx)) != 0 {
				vic.sprDataCounterBase[idx] += 2
			}
		}
		vic.mlIndex = 0
		vic.matrixAccess()

	case 16:
		// Graphics and matrix access, increment sprDataCounterBase by 1 if y expansion FlipFlop is set
		// and check if sprite DMA can be turned off
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.core.sprExpY & mask) != 0 {
				vic.sprDataCounterBase[idx]++
			}
			if (vic.sprDataCounterBase[idx] & 0x3f) == 0x3f {
				vic.sprDMAFlags &= ^mask
			}
		}
		vic.matrixAccess()

	case 17:
		// Graphics and matrix access, turn off border in 40 column mode, display window starts here
		if (vic.core.cr2 & 8) != 0 {
			if vic.core.rasterY == vic.core.dyStop {
				vic.borderULOn = true
			} else {
				if (vic.core.cr1 & 0x10) != 0 {
					if vic.core.rasterY == vic.core.dyStart {
						vic.borderULOn = false
						vic.graphics.SetBorderOn(false)
					} else if !vic.borderULOn {
						vic.graphics.SetBorderOn(false)
					}
				} else if !vic.borderULOn {
					vic.graphics.SetBorderOn(false)
				}
			}
		}
		// Second sample of border state
		vic.graphics.SetBorderOnSample(1)
		if vic.drawThisLine {
			if vic.borderULOn {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawBackground(vic.lineOffset)
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		vic.matrixAccess()

	case 18:
		// Turn off border in 38 column mode
		if (vic.core.cr2 & 8) == 0 {
			if vic.core.rasterY == vic.core.dyStop {
				vic.borderULOn = true
			} else {
				if (vic.core.cr1 & 0x10) != 0 {
					if vic.core.rasterY == vic.core.dyStart {
						vic.borderULOn = false
						vic.graphics.SetBorderOn(false)
					} else if !vic.borderULOn {
						vic.graphics.SetBorderOn(false)
					}
				} else {
					if !vic.borderULOn {
						vic.graphics.SetBorderOn(false)
					}
				}
			}
		}

		// Third sample of border state
		vic.graphics.SetBorderOnSample(2)
		if vic.drawThisLine {
			if vic.borderULOn {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}

		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		vic.matrixAccess()

		vic.graphics.SetCharDataLast()

	case 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54:
		if vic.drawThisLine {
			if vic.borderULOn {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		// Turn on display and matrix access if Bad Line
		if vic.core.isBadLine {
			vic.displayOn = true
			vic.core.SetBALow()
		}
		vic.matrixAccess()
		vic.graphics.SetCharDataLast()

	case 55:
		// Last graphics access, turn off matrix access, turn on sprite DMA if Y coordinate is
		// right and sprite is enabled, handle sprite y expansion, set BA for sprite 0
		if vic.drawThisLine {
			if vic.borderULOn {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		vic.graphicsAccess()
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		// Invert y expansion FlipFlop (if MYE bit is set)
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.core.mye & mask) != 0 {
				vic.core.sprExpY ^= mask
			}
		}
		vic.checkSpriteDMA()
		if (vic.sprDMAFlags & 0x01) != 0 {
			vic.core.SetBALow()
		} else {
			vic.core.ClearBALow()
		}

	case 56:
		// Turn on border in 38 column mode, turn on sprite DMA if Y coordinate is right and
		// sprite is enabled, set BA for sprite 0, display window ends here
		if (vic.core.cr2 & 8) == 0 {
			vic.graphics.SetBorderOn(true)
		}
		// Fourth sample of border state
		vic.graphics.SetBorderOnSample(3)
		if vic.drawThisLine {
			if vic.borderULOn {
				vic.graphics.DrawBackground(vic.lineOffset)
			} else {
				vic.graphics.DrawGraphics(vic.lineOffset)
			}
			vic.sampleBorder()
		}
		//idleAccess
		vic.core.readByte(0x3fff)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		vic.checkSpriteDMA()
		if (vic.sprDMAFlags & 0x01) != 0 {
			vic.core.SetBALow()
		}

	// Turn on border in 40 column mode, set BA for sprite 1, paint sprites
	case 57:
		if (vic.core.cr2 & 8) != 0 {
			vic.graphics.SetBorderOn(true)
		}
		// Fifth sample of border state
		vic.graphics.SetBorderOnSample(4)
		vic.sprites.SetSpriteFlags(vic.sprDisplayFlags)
		// Turn off sprite display if DMA is off
		for idx, mask := 0, uint8(1); idx < 8; idx, mask = idx+1, mask<<1 {
			if (vic.sprDisplayFlags&mask) != 0 && (vic.sprDMAFlags&mask) == 0 {
				vic.sprDisplayFlags &= ^mask
			}
		}
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		//idleAccess
		vic.core.readByte(0x3fff)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x02) != 0 {
			vic.core.SetBALow()
		}
	case 58:
		// Fetch sprite pointer 0, sprDataCounterBase->sprDataCounter, turn on sprite display if necessary,
		// turn off display if RC=7, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		rasterY := vic.core.rasterY & 0xff
		for idx, mask := 0, uint8(1); idx < SpriteNumber; idx, mask = idx+1, mask<<1 {
			vic.sprDataCounter[idx] = vic.sprDataCounterBase[idx]
			if (vic.sprDMAFlags&mask) != 0 && (rasterY == uint16(vic.core.mXy[idx])) {
				vic.sprDisplayFlags |= mask
			}
		}
		vic.fetchSpriteDataPtr(0)
		vic.fetchSpriteData(0, 0)
		if vic.rowCounter == 7 {
			vic.videoCounterBase = vic.videoCounter
			vic.displayOn = false
		}
		if vic.core.isBadLine || vic.displayOn {
			vic.rowCounter = (vic.rowCounter + 1) & 7
			vic.displayOn = true
		}
	case 59:
		// Set BA for sprite 2, read data of sprite 0
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
		}
		vic.fetchSpriteData(0, 1)
		vic.fetchSpriteData(0, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x04) != 0 {
			vic.core.SetBALow()
		}

	// Fetch sprite pointer 1, reset BA if sprite 1 and 2 off, graphics display ends here
	case 60:
		if vic.drawThisLine {
			vic.graphics.DrawBackground(vic.lineOffset)
			vic.sampleBorder()
			vic.sprites.Draw(vic.lineStart)
			vic.graphics.DrawBorder(vic.lineStart)
			// Increment pointer in chunky buffer
			vic.lineStart += DisplayX
		}
		vic.fetchSpriteDataPtr(1)
		vic.fetchSpriteData(1, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x06) == 0 {
			vic.core.ClearBALow()
		}

	// Set BA for sprite 3, read data of sprite 1
	case 61:
		vic.fetchSpriteData(1, 1)
		vic.fetchSpriteData(1, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x08) != 0 {
			vic.core.SetBALow()
		}

	// Read sprite pointer 2, reset BA if sprite 2 and 3 off, read data of sprite 2
	case 62:
		vic.fetchSpriteDataPtr(2)
		vic.fetchSpriteData(2, 0)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if (vic.sprDMAFlags & 0x0c) == 0 {
			vic.core.ClearBALow()
		}

	// Set BA for sprite 4, read data of sprite 2
	case 63:
		vic.fetchSpriteData(2, 1)
		vic.fetchSpriteData(2, 2)
		if vic.core.isBadLine {
			vic.displayOn = true
		}
		if vic.core.rasterY == vic.core.dyStop {
			vic.borderULOn = true
		} else if (vic.core.cr1&0x10) != 0 && vic.core.rasterY == vic.core.dyStart {
			vic.borderULOn = false
		}
		if (vic.sprDMAFlags & 0x10) != 0 {
			vic.core.SetBALow()
		}
		lastCycle = true
	}
	vic.core.rasterX += 8
	if lastCycle {
		vic.cycle = 1
	} else {
		vic.cycle++
	}
	return vBlank, lastCycle
}

func (vic *MOS6569) refreshAccess() {
	_ = vic.core.readByte(0x3f00 | uint16(vic.refreshCounter))
	vic.refreshCounter--
}

func (vic *MOS6569) graphicsAccess() {
	if vic.displayOn {
		var addr uint16
		if (vic.core.cr1 & 0x20) != 0 {
			addr = ((vic.videoCounter & 0x03ff) << 3) | vic.core.bitmapBase | vic.rowCounter // Bitmap
		} else {
			addr = (uint16(vic.matrixLine[vic.mlIndex]) << 3) | vic.core.charBase | vic.rowCounter // Text
		}
		if (vic.core.cr1 & 0x40) != 0 {
			addr &= 0xf9ff // ECM
		}
		vic.graphics.SetGfxData(vic.core.readByte(addr))
		vic.graphics.SetCharData(vic.matrixLine[vic.mlIndex])
		vic.graphics.SetColorData(vic.colorLine[vic.mlIndex])
		vic.mlIndex++
		vic.videoCounter++
	} else {
		if (vic.core.cr1 & 0x40) != 0 {
			vic.graphics.SetGfxData(vic.core.readByte(0x39ff))
		} else {
			vic.graphics.SetGfxData(vic.core.readByte(0x3fff))
		}
		vic.graphics.SetColorData(0)
		vic.graphics.SetCharData(0)
	}
}

func (vic *MOS6569) matrixAccess() {
	if vic.core.baLow {
		if vic.core.quartz.Cycle()-vic.core.baLowFirstCycle < 3 {
			vic.colorLine[vic.mlIndex] = 0xff
			vic.matrixLine[vic.mlIndex] = 0xff
		} else {
			addr := (vic.videoCounter & 0x03ff) | vic.core.matrixBase
			vic.matrixLine[vic.mlIndex] = vic.core.readByte(addr)
			vic.colorLine[vic.mlIndex] = vic.core.banks.ReadColor(addr & 0x03ff)
		}
	}
}

func (vic *MOS6569) fetchSpriteDataPtr(num int) {
	addr := vic.core.matrixBase | 0x03f8 | uint16(num)
	data := vic.core.readByte(addr)
	ptr := uint16(data) << 6
	vic.sprites.SetSpritePtr(num, ptr)
}

func (vic *MOS6569) fetchSpriteData(num int, byteNum int) {
	if (vic.sprDMAFlags & (1 << num)) != 0 {
		ptr := vic.sprites.GetSpritePtr(num)
		addr := (vic.sprDataCounter[num] & 0x3f) | ptr
		data := vic.core.readByte(addr)
		vic.sprites.SetSpriteData(num, byteNum, data)
		vic.sprDataCounter[num]++
	} else if byteNum == 1 {
		//idleAccess
		vic.core.readByte(0x3fff)
	}
}

func (vic *MOS6569) checkSpriteDMA() {
	rasterY := vic.core.rasterY & 0xff
	for i, mask := 0, uint8(1); i < 8; i, mask = i+1, mask<<1 {
		if (vic.core.me&mask) != 0 && rasterY == uint16(vic.core.mXy[i]) {
			vic.sprDMAFlags |= mask
			vic.sprDataCounterBase[i] = 0
			if (vic.core.mye & mask) != 0 {
				vic.core.sprExpY &= ^mask
			}
		}
	}
}

func (vic *MOS6569) sampleBorder() {
	vic.graphics.SetBorderColorSample(vic.cycle)
	vic.lineOffset += 8
	vic.foreMask.Increment()
}
